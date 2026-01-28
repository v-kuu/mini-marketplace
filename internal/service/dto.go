package service

type ProductPatch struct {
	Name *string `json:"name"`
	Price *int64 `json:"price"`
}
