package service

type ProductCreate struct {
	Name string `json:"name"`
	Price int64 `json:"price"`
}

type ProductUpdate struct {
	Name string `json:"name"`
	Price int64 `json:"price"`
}

type ProductPatch struct {
	Name *string `json:"name,omitempty"`
	Price *int64 `json:"price,omitempty"`
}
