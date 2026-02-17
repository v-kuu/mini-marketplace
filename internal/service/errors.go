package service

import "errors"

var (
	ErrInvalidProduct = errors.New("invalid product")
	ErrProductNotFound = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
)
