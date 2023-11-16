package main

import (
	"context"
	oas "github.com/Cj-bc/skkishoe/internal/oas"
	"github.com/mattn/go-skkdic"
	"net/http"
	"strings"
)

type CandidatesService struct {
	dict *skkdic.Dict
}


// Construct result for 'text/*' Mime types
func TextResult(cs []oas.Candidate) *oas.CandidatesGetOKText {
	result := []string{}
	for _, c := range cs {
		result = append(result, c.Candidate)
	}
	reader := strings.NewReader(strings.Join(result, "/"))
	res := oas.CandidatesGetOKText{Data: reader}
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

func (s CandidatesService) CandidatesGet(ctx context.Context, args oas.CandidatesGetParams) (oas.CandidatesGetRes, error) {
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

	req, _ := ctx.Value("rawRequest").(*http.Request)

	switch req.Header.Get("Content-Type") {
	case "application/json":
		res := oas.CandidatesGetOKApplicationJSON(result)
		return &res, nil
	default:
		return TextResult(result), nil
	}
}
