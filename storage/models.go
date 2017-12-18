package storage

type Book struct {
	ID	string 		`json:"id, omitempty"`
	Title	string		`json:"title, omitempty"`
	Genres	[]string	`json:"genres, omitempty"`
	Pages	int		`json:"pages, omitempty"`
	Price	float64		`json:"price, omitempty"`
}

type Books []Book

type Filter struct {
	Price	string		`json:"price, omitempty"`

}
