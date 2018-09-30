package main

import (
	"log"
	"testing"
)

func TestHentAntall(T *testing.T) {
	tmp := hentAntall()

	log.Println("Kjørte test ", tmp)
}
