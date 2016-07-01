package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/redis.v4"

	"github.com/TRUMPTOWERS/polypi/ingest"
	"github.com/TRUMPTOWERS/polypi/label"
	"github.com/TRUMPTOWERS/polypi/pie"
)

const piesURL = "http://stash.truex.com/tech/bakeoff/pies.json"
const redisURL = "localhost:6379"

func main() {
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
	labels := make(map[string]*label.Label)
	pies := make([]pie.DataPie, len(in.Pies))

	for i, p := range in.Pies {
		// Create or fetch existing labels
		var pLabels []*label.Label
		for _, l := range p.Labels {
			var lp *label.Label
			var ok bool
			if lp, ok = labels[l]; !ok {
				lp = &label.Label{
					label.DataLabel{Name: l},
				}
			}
			pLabels = append(pLabels, lp)
			labels[l] = lp
		}

		// Create pie objects
		pies[i] = pie.DataPie{
			ID:       p.ID,
			Name:     p.Name,
			ImageURL: p.ImageURL,
			Price:    p.PricePerSlice,
			Slices:   p.Slices,
			Labels:   pLabels,
		}
	}

	// Save pies to redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Make or get all labels
	for _, l := range labels {
		existing := client.Get(fmt.Sprintf("labelName_%v", l.Name))
		id, err := existing.Int64()
		if err == redis.Nil {
			// Create new label
			idRes := client.Incr("labelCount")
			id, err = idRes.Result()
			if err != nil {
				log.Fatalf("Could not increment labelCount: %v", err)
			}
			l.ID = id

			// ID is set, make database match
			data, err := json.Marshal(l)
			if err != nil {
				log.Fatalf("Could not marshal the label: %v", err)
			}
			err = client.Set(fmt.Sprintf("label_%d", l.ID), data, 0).Err()
			if err != nil {
				log.Fatalf("Could not set the id key: %v", err)
			}
			err = client.Set(fmt.Sprintf("labelName_%v", l.Name), l.ID, 0).Err()
			if err != nil {
				log.Fatalf("Could not set the name key: %v", err)
			}
		} else if err != nil {
			log.Fatalf("Could not retrieve label id for %v: %v", l.Name, err)
		}
		l.ID = id
	}

	// Create or replace all pies
	for _, p := range pies {
		data, err := json.Marshal(p)
		if err != nil {
			log.Fatalf("Could not marshal ze pie: %v", err)
		}
		err = client.Set(fmt.Sprintf("pie_%d", p.ID), data, 0).Err()
		if err != nil {
			log.Fatalf("Could not set ze pie: %v", err)
		}
	}
}
