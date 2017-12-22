package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/twinj/uuid"
)

var (
	// ErrNotFound describe the state when the object is not found in the storage
	ErrNotFound = errors.New("can't find the book with given ID")
)

type Lib struct {
	storage string
	//storage io.ReadWriteCloser
}

func (l *Lib) writeData(books Books) error {
	path, err := filepath.Abs(l.storage)
	if err != nil {
		return err
	}

	booksBytes, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, booksBytes, 0644)
}

func (l *Lib) wantedIndex(id string, books Books) (int, error) {
	for index, book := range books {
		if id == book.ID {
			return index, nil
		}
	}
	return 0, ErrNotFound
}

//GetBooks returns all book objects
func (l *Lib) GetBooks() (Books, error) {
	var books Books

	path, err := filepath.Abs(l.storage)
	if err != nil {
		return nil, err
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return books, json.Unmarshal(file, &books)
}

// CreateBook adds book object into db
func (l *Lib) CreateBook(book Book) error {
	err := errors.New("not all fields are populated")
	switch {
	case book.Genres == nil:
		return err
	case book.Pages == 0:
		return err
	case book.Price == 0:
		return err
	case book.Title == "":
		return err
	}

	books, err := l.GetBooks()
	if err != nil {
		return err
	}

	book.ID = uuid.NewV4().String()
	books = append(books, book)
	return l.writeData(books)

}

// GetBook returns book object with specified id
func (l *Lib) GetBook(id string) (Book, error) {
	var b Book
	books, err := l.GetBooks()
	if err != nil {
		return b, err
	}

	for _, book := range books {
		if id == book.ID {
			return book, nil
		}
	}
	return b, ErrNotFound
}

// RemoveBook removes book object with specified id
func (l *Lib) RemoveBook(id string) error {
	books, err := l.GetBooks()
	if err != nil {
		return err
	}

	index, err := l.wantedIndex(id, books)
	if err != nil {
		return err
	}
	books = append(books[:index], books[index+1:]...)
	return l.writeData(books)
}

// ChangeBook updates book object with specified id
func (l *Lib) ChangeBook(id string, changedBook Book) error {
	books, err := l.GetBooks()
	if err != nil {
		return err
	}

	index, err := l.wantedIndex(id, books)
	if err != nil {
		return err
	}

	book := &books[index]
	book.Price = changedBook.Price
	book.Title = changedBook.Title
	book.Pages = changedBook.Pages
	book.Genres = changedBook.Genres
	err = l.writeData(books)
	return err
}

// PriceFilter returns filtered book objects
func (l *Lib) PriceFilter(filter BookFilter) (Books, error) {
	var wantedBooks Books

	if len(filter.Price) <= 1 {
		return nil, errors.New("Not valid data")
	}
	operator := string(filter.Price[0])
	if operator != "<" && operator != ">" {
		err := errors.New("unsupported operation")
		return nil, err
	}

	books, err := l.GetBooks()
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(filter.Price[1:], 64)
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		if operator == ">" {
			if book.Price > price {
				wantedBooks = append(wantedBooks, book)
			}
		} else {
			if book.Price < price {
				wantedBooks = append(wantedBooks, book)
			}
		}
	}
	return wantedBooks, nil
}
