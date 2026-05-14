package model

type Product struct {
	ID string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Price int64 `json:"price" bson:"price"`
}
