package model

type ItemResponse struct {
	Id          UUID
	Name        string
	Description string
	ItemImage   string
	Count       int
}
