package storage

import (
	// hook for using Genres as string array in sqlite3
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/twinj/uuid"
)

var (
	// ErrNotFound describes the state when the object is not found in the storage
	ErrNotFound = errors.New("can't find the book with given ID")
	// ErrNotValidData describes the state when the object is not valid
	ErrNotValidData = errors.New("not valid data")
	// ErrUnsupportedOperation describes the state when the object is an unsupported operation
	ErrUnsupportedOperation = errors.New("unsupported operation")
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

func (b *Book) PrepareToCreate() {
	b.ID = uuid.NewV4().String()
}
