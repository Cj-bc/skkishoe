package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mattn/go-skkdic"
	"github.com/ogen-go/ogen/middleware"

	oas "github.com/Cj-bc/skkishoe/internal/oas"
)

// Insert raw *http.Request pointer to context
func StoreRawRequestMiddleware(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	req.SetContext(context.WithValue(req.Context, "rawRequest", req.Raw))
	return next(req)
}

type dicts struct {
	dicts []string
}

func (d dicts) String() string {
	return strings.Join(d.dicts, ":")
}

func (d *dicts) Set(value string) error {
	for _, path := range strings.Split(value, ":") {
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}
		d.dicts = append(d.dicts, path)
	}
	return nil
}

var (
	flag_dicts dicts
)

func init() {
	flag.Var(&flag_dicts, "dict", "Dictionaries to use. Must be a Valid file path joined by `:'\ne.g. foo/bar.dict:bar/baz.dict")
}

func main() {
	port := flag.Int("port", 8080, "server port number")

	flag.Parse()

	dict := skkdic.New()
	var err error
	if len(flag_dicts.dicts) > 0 {
		for _, d := range flag_dicts.dicts {
			f, err := os.Open(d)
			defer f.Close()
			if err != nil {
				log.Fatalf("Coul'd not open specified dictionary %s: %w", d, err)
			}

			err = dict.Load(f)
		}
	} else {
		log.Fatal("You need to specify Dictionary to use. Aborting...")
	}
	if err != nil {
		log.Fatal(err)
	}

	service := CandidatesService{dict}

	srv, err := oas.NewServer(service, oas.WithMiddleware(StoreRawRequestMiddleware))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), srv))
}
