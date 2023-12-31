// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// MidashisMidashiGet implements GET /midashis/{midashi} operation.
//
// Returns candidates for specified midashi.
//
// GET /midashis/{midashi}
func (UnimplementedHandler) MidashisMidashiGet(ctx context.Context, params MidashisMidashiGetParams) (r MidashisMidashiGetRes, _ error) {
	return r, ht.ErrNotImplemented
}
