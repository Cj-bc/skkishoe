package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/mattn/go-skkdic"
	"github.com/ogen-go/ogen/middleware"

	oas "github.com/Cj-bc/skkishoe/internal/oas"
)

// Insert raw *http.Request pointer to context
func StoreRawRequestMiddleware(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	req.SetContext(context.WithValue(req.Context, "rawRequest", req.Raw))
	return next(req)
}

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

	srv, err := oas.NewServer(service, oas.WithMiddleware(StoreRawRequestMiddleware))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe("localhost:8080", srv))
}
