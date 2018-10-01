package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Antall ... Resultatstrukturen til grensesnittet
type Antall struct {
	Stillinger int `json:"stillinger"`
	Annonser   int `json:"annonser"`
	Nye        int `json:"nye"`
}

var antallStillinger Antall
var mutex = &sync.Mutex{}

func hentAntall() Antall {

	// Generert fra eksempel-resultat
	// Slettet en del overflødig fra resultatet.
	// https://mholt.github.io/json-to-go/
	type AutoGenerated struct {
		Aggregations struct {
			Published struct {
				DocCount int `json:"doc_count"`
				Range    struct {
					Buckets []struct {
						DocCount int `json:"doc_count"`
					} `json:"buckets"`
				} `json:"range"`
			} `json:"published"`
		} `json:"aggregations"`
	}

	query := `
	{
		"size": 0,
		"post_filter": {
		  "bool" :  {
			"filter":  [ { "term" : { "status" : "ACTIVE" } } ]
		  }
		},
		"aggs": {
		  "published": {
			"filter": {
			  "bool" :  {
				"filter":  [ { "term" : { "status" : "ACTIVE" } } ]
			  }
			},
			"aggs": {
				"range": {
					"date_range": {
					"field": "published",
					"ranges": [
						{
						"key": "now-1d",
						"from": "now-1d"
						}
					]
					}
				}
			}
		}
	  }
	}`
	tmp := Antall{}

	res, err := http.Post("https://stillingsok.nav.no/pam-stillingsok/search-api/stillingsok/ad/_search",
		"application/json", strings.NewReader(query))

	if err != nil {
		log.Println("ERROR: ", err)
		return tmp
	}

	var body AutoGenerated
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Println("ERROR: ", err)
		return tmp
	}

	// XXX: Vi aggregerer ikke antall stillinger.
	tmp.Annonser = body.Aggregations.Published.DocCount
	tmp.Stillinger = body.Aggregations.Published.DocCount
	if len(body.Aggregations.Published.Range.Buckets) > 0 {
		tmp.Nye = body.Aggregations.Published.Range.Buckets[0].DocCount
	}
	return tmp
}

func stillinger(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	tmp := antallStillinger
	mutex.Unlock()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s, _ := json.Marshal(tmp)
	w.Write(s)
}

func main() {
	var (
		portnummer string
		prefix     string
	)
	flag.StringVar(&portnummer, "p", ":8085", "Hvilket portnummer/adresse")
	flag.StringVar(&prefix, "prefix", "/api", "Løsningens adresse")
	flag.Parse()

	{
		ticker := time.NewTicker(5000 * time.Millisecond)
		go func() {
			for t := range ticker.C {
				tmp := hentAntall()
				mutex.Lock()
				antallStillinger = tmp
				mutex.Unlock()
				log.Println("Oppdaterer ", t)
			}
		}()
	}

	http.HandleFunc(prefix+"/stillinger", stillinger) // set router
	err := http.ListenAndServe(portnummer, nil)       // set listen port
	if err != nil {
		log.Fatal("ERROR:", err)
	}
	log.Println("API stillinger, lytter til adresse  ", portnummer)
}
