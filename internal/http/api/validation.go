package api

import (
	"strings"
)

func validateCreate(req CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidName
	}
	if req.Price <= 0 {
		return ErrInvalidPrice
	}
	
	return nil
}

func validateUpdate(req UpdateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidName
	}
	if req.Price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}

func validatePatch(req PatchProductRequest) error {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return ErrInvalidName
	}
	if req.Price != nil && *req.Price <= 0 {
		return ErrInvalidPrice
	}
	if req.Name == nil && req.Price == nil {
		return ErrEmptyPatch
	}
	return nil
}
