package main

import (
	"context"
	oas "github.com/Cj-bc/skkishoe/internal/oas"

)

type CandidatesService struct {
	dict map[string][]oas.Candidate
}

func (s CandidatesService) CandidatesGet(ctx context.Context, args oas.CandidatesGetParams) ([]oas.Candidate, error) {
	if cands, ok := s.dict[args.Midashi]; ok {
		return cands, nil
	} else {
		return []oas.Candidate{}, nil
	}
}
