// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/validate"
)

func (s CandidatesGetOKApplicationJSON) Validate() error {
	alias := ([]Candidate)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	if err := (validate.Array{
		MinLength:    0,
		MinLengthSet: true,
		MaxLength:    0,
		MaxLengthSet: false,
	}).ValidateLength(len(alias)); err != nil {
		return errors.Wrap(err, "array")
	}
	return nil
}
