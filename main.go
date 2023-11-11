package main

import (
	"log"
	"net/http"

	oas "github.com/Cj-bc/skkishoe/internal/oas"
)

func main() {
	tmp := map[string][]oas.Candidate {
		"おむすび": {
			{
				Candidate: oas.NewOptString("御掬"),
				Annotation: oas.NewOptString("VTuber 御掬この子"),
			},
		},
		"きょう": {
			{
				Candidate: oas.NewOptString("今日"),
				Annotation: oas.OptString{},
			},
			{
				Candidate: oas.NewOptString("京"),
				Annotation: oas.OptString{},
			},
		},
	}

	service := CandidatesService{tmp}

	srv, err := oas.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe("localhost:8080", srv))
}
