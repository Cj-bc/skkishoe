package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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


func usage() {
	fmt.Fprintf(os.Stderr, `skkishoe usage:
  skkishoe DICTIONARY... # Start skkishoe server using DICTIONARY

Flags:
`)
	flag.PrintDefaults()
	os.Exit(0)
}

func init() {
	flag.Var(&flag_dicts, "dict", "Dictionaries to use. Must be a Valid file path joined by `:'\ne.g. foo/bar.dict:bar/baz.dict")
}

func main() {
	port := flag.Int("port", 8080, "server port number")

	flag.Parse()
	if flag.NArg() <= 0 {
		usage()
	}

	slog.Info("Setting up dictionaries", "dictionaries", flag.Args())
	dict := skkdic.New()
	var err error
	for _, d := range flag.Args() {
		f, err := os.Open(d)
		defer f.Close()
		if err != nil {
			slog.Warn("Failed to open dictionary", "dictionary", d, "error", err)
		}

		if err = dict.Load(f); err != nil {
			slog.Warn("Get error while reading dictionary file", "dictionary", d, "error", err)
		} else {
			slog.Info(fmt.Sprintf("dictionary loaded: %s", d))
		}
	}

	if err != nil {
		slog.Error("Unknown error occured", "error", err.Error())
		os.Exit(1)
	}

	service := CandidatesService{dict}

	slog.Info("Setting up Server")
	srv, err := oas.NewServer(service, oas.WithMiddleware(StoreRawRequestMiddleware))
	if err != nil {
		slog.Error("Unknown error occured", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("Server is ready. Start listening", "port", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), srv); err != nil {
		slog.Error("Unknown error occured", "error", err.Error())
	}
	os.Exit(1)
}
