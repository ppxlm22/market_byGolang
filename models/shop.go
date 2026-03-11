package models

type Category struct {
	ID int
	Name string
}
type Products struct {
	ID int
	Name string
	Price float64
	Stock int
	CategoryID int
}