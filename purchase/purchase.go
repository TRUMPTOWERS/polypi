package purchase

import "github.com/TRUMPTOWERS/polypi/customer"

type DataPurchase struct {
	Customer customer.Customer
	Amount   int
}

type Purchase struct {
	DataPurchase
}
