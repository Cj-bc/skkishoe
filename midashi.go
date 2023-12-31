package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	oas "github.com/Cj-bc/skkishoe/internal/oas"
	"github.com/mattn/go-skkdic"
)

type MidashiService struct {
	dict *skkdic.Dict
}


// Construct result for 'text/*' Mime types
func TextResult(cs []oas.Candidate) *oas.MidashisMidashiGetOKTextPlain {
	cands := []string{}
	for _, c := range cs {
		cands = append(cands, c.Candidate + ";" + c.Annotation.Or(""))
	}
	reader := strings.NewReader(strings.Join(cands, "/"))
	res := oas.MidashisMidashiGetOKTextPlain{Data: reader}
	return &res
}

func isAlphabet(r rune) bool {
		return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

// Convert skkdic.Entry to oas.Candidate
func entryToCandidates(e skkdic.Entry) []oas.Candidate {
	candidates := []oas.Candidate{}
	for _, w := range e.Words {
		candidates = append(candidates, oas.Candidate{
			Candidate: w.Text,
			Annotation: oas.NewOptString(w.Desc)})
	}
	return candidates
}

func (s MidashiService) MidashisMidashiGet(ctx context.Context, args oas.MidashisMidashiGetParams) (oas.MidashisMidashiGetRes, error) {
	slog.Info("GET", "path", fmt.Sprintf("/midashi/%s", args.Midashi))

	entries := []skkdic.Entry{}
	rs := []rune(args.Midashi)
	if isAlphabet(rs[len(rs)-1]) {
		entries = s.dict.SearchOkuriAri(args.Midashi)
	} else {
		entries = s.dict.SearchOkuriNasi(args.Midashi)
	}

	result := []oas.Candidate{}
	for _, e := range entries {
		for _, c := range entryToCandidates(e) {
			result = append(result, c)
		}
	}

	// If program could not retrive http.Request from context,
	// Assume that contentType is text/plain
	//
	// TODO: What's best?
	var contentType string = "text/plain"
	req, ok := ctx.Value("rawRequest").(*http.Request)
	if ok {
		contentType = req.Header.Get("Content-Type")
	}

	switch contentType {
	case "application/json":
		res := oas.MidashisMidashiGetOKApplicationJSON(result)
		return &res, nil
	default:
		return TextResult(result), nil
	}
}
