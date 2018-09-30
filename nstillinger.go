package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

// Antall ... holder resultatet
type Antall struct {
	Stillinger int `json:"stillinger"`
	Annonser   int `json:"annonser"`
	Nye        int `json:"nye"`
}

var antallStillinger Antall

func stillinger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s, _ := json.Marshal(antallStillinger)
	w.Write(s)
}

func oppdaterAntall() {
	// TODO:
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
		///antallStillinger := Antall{}
		ticker := time.NewTicker(5000 * time.Millisecond)
		go func() {
			for t := range ticker.C {
				oppdaterAntall()
				log.Println("Oppdaterer ", t)
			}
		}()
	}

	checkErr := func(err error) {
		if err != nil {
			log.Fatal("ERROR:", err)
		}
	}
	http.HandleFunc(prefix+"/stillinger", stillinger) // set router
	err := http.ListenAndServe(portnummer, nil)       // set listen port
	checkErr(err)
	log.Println("API stillinger, lytter til adresse  ", portnummer)
}
