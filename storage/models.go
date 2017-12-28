package storage

import (
	// hook for using Genres as string array in sqlite3
	"github.com/lib/pq"
)

// Book describes main data structure in the app
type Book struct {
	ID     string         `gorm:"type:varchar(100);primary_key" json:"id, omitempty"`
	Title  string         `gorm:"type:varchar(100)" json:"title, omitempty"`
	Genres pq.StringArray `gorm:"type:varchar(64)" json:"genres, omitempty"`
	Pages  int            `gorm:"type:int" json:"pages, omitempty"`
	Price  float64        `gorm:"type:real" json:"price, omitempty"`
}

// Books contains book objects
type Books []Book

//Filter describes filter indicator
type BookFilter struct {
	Price string `gorm:"type:varchar(100)" json:"price, omitempty"`
}
