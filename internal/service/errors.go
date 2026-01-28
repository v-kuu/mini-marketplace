package service

import "errors"

var (
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidProduct = errors.New("invalid product")
	ErrProductNotFound = errors.New("product not found")
)
