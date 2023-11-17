package main

import (
	"context"
	"net/http"
	"strings"

	oas "github.com/Cj-bc/skkishoe/internal/oas"
	"github.com/mattn/go-skkdic"
)

type CandidatesService struct {
	dict *skkdic.Dict
}


// Construct result for 'text/*' Mime types
func TextResult(cs []oas.Candidate) *oas.CandidatesGetOKText {
	cands := []string{}
	annos := []string{}
	for _, c := range cs {
		cands = append(cands, c.Candidate)
		annos = append(annos, c.Annotation.Or(""))
	}
	reader := strings.NewReader(strings.Join(cands, "/")+"\n"+strings.Join(annos, "/"))
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
