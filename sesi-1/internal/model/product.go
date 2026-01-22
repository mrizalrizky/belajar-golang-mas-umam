package model

// Private and Public variable
// Uppercase -> ID (Public)
// Lowercase -> id (Private)
// Data modelling pake struct
type Product struct {
	ID 		int		`json:"id"` // -> backtick ibarat untuk aliasing json keynya
	Name 	string	`json:"name" validate:"required,min=3"`
	Price 	int		`json:"price" validate:"required,gte=0"`
	Stock 	int		`json:"stock" validate:"required,gte=0"`
}
