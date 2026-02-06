package api

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidPrice = errors.New("invalid price")
	ErrEmptyPatch = errors.New("empty patch")
)
