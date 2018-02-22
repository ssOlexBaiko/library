package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type library struct {
	mu      sync.RWMutex // Use the force Luke!
	storage *os.File     // Here you can put opened os.File object. After that you will be able to implement concurrent safe operations with file storage
}

func (l library) Close() error {
	return l.storage.Close()
}

// NewLibrary constructor for library struct.
// Constructors are often used for initialize some data structures (map, slice, chan...)
// or when you need some data preparation
// or when you want to start some watchers (goroutines). In this case you also have to think about Close() method.
func NewLibrary(path string) (*library, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)
	// https://joeshaw.org/dont-defer-close-on-writable-files/
	// defer file.Close()
	if err != nil {
		return nil, err
	}
	return &library{storage: file}, nil
}

func (l *library) wantedIndex(id string, books Books) (int, error) {
	for index, book := range books {
		if id == book.ID {
			return index, nil
		}
	}
	return 0, ErrNotFound
}

//GetBooks returns all book objects
func (l *library) GetBooks() (Books, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.getBooks()

}
func (l *library) getBooks() (Books, error) {
	var books Books

	_, err := l.storage.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	file, err := ioutil.ReadAll(l.storage)
	if err != nil {
		return nil, err
	}
	return books, json.Unmarshal(file, &books)
}

// CreateBook adds book object into db
func (l *library) CreateBook(book Book) error {
	l.mu.Lock()
	defer l.mu.Unlock()

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

	book.PrepareToCreate()
	books, err := l.getBooks()
	if err != nil {
		return err
	}

	books = append(books, book)
	byteBooks, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}
	// hook for "clearing" file, seeking to it's zero position,
	// and then writing updated Books to file
	err = l.storage.Truncate(0)
	if err != nil {
		return err
	}
	_, err = l.storage.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = l.storage.Write(byteBooks)
	return err
}

// GetBook returns book object with specified id
func (l *library) GetBook(id string) (Book, error) {
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
func (l *library) RemoveBook(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	books, err := l.getBooks()
	if err != nil {
		return err
	}

	index, err := l.wantedIndex(id, books)
	if err != nil {
		return err
	}
	books = append(books[:index], books[index+1:]...)
	byteBooks, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}
	// hook for "clearing" file, seeking to it's position to zero,
	// and then writing updated Books
	err = l.storage.Truncate(0)
	if err != nil {
		return err
	}
	_, err = l.storage.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = l.storage.Write(byteBooks)
	return err
}

// ChangeBook updates book object with specified id
func (l *library) ChangeBook(changedBook Book) (Book, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// заглушка для повернення помилки
	var b Book
	books, err := l.getBooks()
	if err != nil {
		return b, err
	}

	index, err := l.wantedIndex(changedBook.ID, books)
	if err != nil {
		return b, err
	}

	// What would happen if user omits some field?
	book := &books[index]
	book.Price = changedBook.Price
	book.Title = changedBook.Title
	book.Pages = changedBook.Pages
	book.Genres = changedBook.Genres
	byteBooks, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return b, err
	}
	// hook for "clearing" file, seeking to it's position to zero,
	// and then writing updated Books
	err = l.storage.Truncate(0)
	if err != nil {
		return b, err
	}
	_, err = l.storage.Seek(0, 0)
	if err != nil {
		return b, err
	}
	_, err = l.storage.Write(byteBooks)

	return *book, err
}

// PriceFilter returns filtered book objects
func (l *library) PriceFilter(filter BookFilter) (Books, error) {
	var wantedBooks Books

	if len(filter.Price) <= 1 {
		return nil, ErrNotValidData
	}
	operator := string(filter.Price[0])
	if operator != "<" && operator != ">" {
		return nil, ErrUnsupportedOperation
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
