package api

type CreateProductRequest struct {
	Name string `json:"name"`
	Price int64 `json:"price"`
}

type UpdateProductRequest struct {
	Name string `json:"name"`
	Price int64 `json:"price"`
}

type PatchProductRequest struct {
	Name *string `json:"name,omitempty"`
	Price *int64 `json:"price,omitempty"`
}
