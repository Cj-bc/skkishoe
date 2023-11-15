package main

import (
	"log"
	"net/http"
	"github.com/mattn/go-skkdic"
	"os"

	oas "github.com/Cj-bc/skkishoe/internal/oas"
)

func main() {
	f, err := os.Open("/usr/share/skk/SKK-JISYO.L.utf-8")
	if err != nil {
		log.Fatal(err)
	}

	dict := skkdic.New()
	err = dict.Load(f)
	if err != nil {
		log.Fatal(err)
	}

	service := CandidatesService{dict}

	srv, err := oas.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe("localhost:8080", srv))
}
