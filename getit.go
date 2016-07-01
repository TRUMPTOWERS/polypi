package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/redis.v4"

	"github.com/TRUMPTOWERS/polypi/ingest"
	"github.com/TRUMPTOWERS/polypi/label"
	"github.com/TRUMPTOWERS/polypi/pie"
)

const piesURL = "http://stash.truex.com/tech/bakeoff/pies.json"
const redisURL = "localhost:6379"

func main() {
	logger := log.New(os.Stderr, "ingestion: ", log.Llongfile)
	redis.SetLogger(logger)

	in := ingest.Ingest{}
	// HTTP to get pies
	res, err := http.Get(piesURL)
	if err != nil {
		log.Fatalf("Could not get pies: %v", err)
	}

	// Decode the JSON
	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&in)
	if err != nil {
		log.Fatalf("Could not decode shtuff: %v", err)
	}

	// Create pies objects
	labels := make(map[string]label.DataLabel)

	for _, p := range in.Pies {
		// Make sure all labels for this pie exist in the labels map
		for _, l := range p.Labels {
			var lp label.DataLabel
			var ok bool
			if lp, ok = labels[l]; !ok {
				lp = label.DataLabel{Name: l}
			}
			labels[l] = lp
		}
	}

	// We officially need a client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Make or get all labels
	for i := range labels {
		existing := client.Get(fmt.Sprintf(label.IDByName, labels[i].Name))
		id, err := existing.Int64()
		if err == redis.Nil {
			// Create new label
			id, err = client.Incr(label.CountKey).Result()
			if err != nil {
				log.Fatalf("Could not increment labelCount: %v", err)
			}

			// ID is set, make database match
			data, err := json.Marshal(labels[i])
			if err != nil {
				log.Fatalf("Could not marshal the label: %v", err)
			}
			err = client.Set(fmt.Sprintf(label.Key, id), data, 0).Err()
			if err != nil {
				log.Fatalf("Could not set the id key: %v", err)
			}
			err = client.Set(fmt.Sprintf(label.IDByName, labels[i].Name), id, 0).Err()
			if err != nil {
				log.Fatalf("Could not set the name key: %v", err)
			}
		} else if err != nil {
			log.Fatalf("Could not retrieve label id for %s: %v", labels[i].Name, err)
		}
		this := labels[i]
		this.ID = id
		labels[i] = this
	}

	// Create or replace all pies
	for _, p := range in.Pies {
		// Make the datapie struct
		thisPie := pie.DataPie{
			ID:       p.ID,
			Name:     p.Name,
			ImageURL: p.ImageURL,
			Price:    p.PricePerSlice,
			Slices:   p.Slices,
			Labels:   make([]int64, len(p.Labels)),
		}

		// Fill in the label ids
		for i, l := range p.Labels {
			thisPie.Labels[i] = labels[l].ID
		}

		data, err := json.Marshal(thisPie)
		if err != nil {
			log.Fatalf("Could not marshal ze pie: %v", err)
		}
		err = client.Set(fmt.Sprintf(pie.Key, p.ID), data, 0).Err()
		if err != nil {
			log.Fatalf("Could not set ze pie: %v", err)
		}
	}
}
