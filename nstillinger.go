package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func stillinger(w http.ResponseWriter, r *http.Request) {
	// TODO: Unmarshall
	mock := `{"stillinger":"000","annonser":"000","nye":"0"}`
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(mock))
}

func oppdaterStillinger() {
}

func main() {
	var (
		portnummer string
		prefix     string
	)
	flag.StringVar(&portnummer, "p", ":8085", "Hvilket portnummer/adresse")
	flag.StringVar(&prefix, "prefix", "/api", "Hvilket adresse ligger løsningen på")
	flag.Parse()

	{
		ticker := time.NewTicker(5000 * time.Millisecond)
		go func() {
			for t := range ticker.C {
				oppdaterStillinger()
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
