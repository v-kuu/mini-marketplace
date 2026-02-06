package service

import "errors"

var (
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidProduct = errors.New("invalid product")
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidPrice = errors.New("invalid price")
	ErrEmptyPatch = errors.New("empty patch")
	ErrProductNotFound = errors.New("product not found")
)
