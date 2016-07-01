package piehandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TRUMPTOWERS/polypi/customer"
	"github.com/TRUMPTOWERS/polypi/pie"
	"github.com/TRUMPTOWERS/polypi/purchase"
	"github.com/gorilla/mux"

	"gopkg.in/redis.v4"
)

// sendBack is how we respond to requests for pie data
type sendBack struct {
	Name            string         `json:"name"`
	ImageURL        string         `json:"image_url"`
	Price           float64        `json:"price_per_slice"`
	RemainingSlices int            `json:"remaining_slices"`
	Purchases       []userPurchase `json:"purchases"`
}

type userPurchase struct {
	Username string `json:"username"`
	Slices   int    `json:"slices"`
}

// Handler impliments the http.Handler interface
type Handler struct {
	// DS is our datastore
	DS *redis.Client
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	// id isn't an integer
	if err != nil {
		http.NotFound(rw, r)
		return
	}

	data, err := h.DS.Get(fmt.Sprintf(pie.Key, id)).Bytes()
	if err != nil {
		http.Error(rw, "Could not retrieve data", http.StatusInternalServerError)
		return
	}

	thisPie := &pie.Pie{}

	err = json.Unmarshal(data, &thisPie.DataPie)
	if err != nil {
		http.Error(rw, "Data appears malformed", http.StatusInternalServerError)
	}

	thisPie.Populate()

	retPie := &sendBack{
		Name:            thisPie.Name,
		ImageURL:        thisPie.ImageURL,
		Price:           thisPie.Price,
		RemainingSlices: thisPie.RemainingSlices,
	}

	// retrieve data about the purchases
	thesePurchases := map[string]*userPurchase{}

	for _, p := range thisPie.Purchases {
		thisP := purchase.DataPurchase{}
		data, err := h.DS.Get(fmt.Sprintf(purchase.Key, p)).Bytes()
		if err != nil {
			http.Error(rw, "Could not get purchase data", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(data, &thisP)
		if err != nil {
			http.Error(rw, "Purchase data appears malformed", http.StatusInternalServerError)
			return
		}

		thisCustomer := customer.DataCustomer{}
		data, err = h.DS.Get(fmt.Sprintf(customer.Key, thisP.Customer)).Bytes()
		if err != nil {
			http.Error(rw, "Could not get customer for relevant purchase", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(data, &thisCustomer)
		if err != nil {
			http.Error(rw, "Customer data for relevenat purchase appears malformed", http.StatusInternalServerError)
		}

		if cs, ok := thesePurchases[thisCustomer.Username]; ok {
			cs.Slices++
		} else {
			thesePurchases[thisCustomer.Username] = &userPurchase{
				Username: cs.Username,
				Slices:   1,
			}
		}
	}

	retPie.Purchases = make([]userPurchase, len(thesePurchases))
	i := 0
	for _, p := range thesePurchases {
		retPie.Purchases[i] = *p
		i++
	}

	encoder := json.NewEncoder(rw)
	encoder.Encode(retPie)
}
