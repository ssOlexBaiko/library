package storage

// Book describes main data structure in the app
type Book struct {
	ID     string   `json:"id, omitempty"`
	Title  string   `json:"title, omitempty"`
	Genres []string `json:"genres, omitempty"`
	Pages  int      `json:"pages, omitempty"`
	Price  float64  `json:"price, omitempty"`
}

// Books contains book objects
type Books []Book

//Filter describes filter indicator
type BookFilter struct {
	Price string `json:"price, omitempty"`
}
