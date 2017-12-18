package storage

import (
	"io/ioutil"
	"path/filepath"
	"github.com/twinj/uuid"
	"encoding/json"
	"errors"
	"strconv"
)

func writeData(books Books) error {
	path, err := filepath.Abs("storage/storage.json")
	if err != nil {
		return err
	}
	booksBytes, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, booksBytes, 0644)
	return err
}

func wantedIndex(id string, books Books) (int, error) {
	for index, book := range books {
		if id == book.ID {
			return index, nil
		}
	}
	err := errors.New("can't find the book with given ID")
	return 0, err
}

//GetBooks returns all book objects
func GetBooks() (Books, error) {
	var books Books
	path, err := filepath.Abs("storage/storage.json")
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(path)
    	if err != nil {
		return nil, err
    	}
	err = json.Unmarshal(file, &books)
	if err != nil {
		return nil, err
	}
	return books, nil
}

//CreateBook adds book object into db
func CreateBook(book Book) error {
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
	default:
		books, err := GetBooks()
		if err != nil {
			return err
		}
		book.ID = uuid.NewV4().String()
		books = append(books, book)
		err = writeData(books)
		if err != nil {
			return err
		}
		return nil
	}
}

//GetBook returns book object with specified id
func GetBook(id string) (Book, error) {
	var b Book
	books, err := GetBooks()
	if err != nil {
		return b, err
	}
	for _, book := range books {
		if id == book.ID {
			return book, nil
		}
	}
	err  = errors.New("can't find the book with given ID")
	return b, err
}

//RemoveBook removes book object with specified id
func RemoveBook(id string) error {
	books, err := GetBooks()
	if err != nil {
		return err
	}
	index, err := wantedIndex(id, books)
	if err != nil {
		return err
	}
	books = append(books[:index], books[index+1:]...)
	err = writeData(books)
	return err
}

//ChangeBook updates book object with specified id
func ChangeBook(id string, changedBook Book) error {
	books, err := GetBooks()
	if err != nil {
		return err
	}
	index, err := wantedIndex(id, books)
	if err != nil {
		return err
	}
	book := &books[index]
	if changedBook.Price != 0 {
		book.Price = changedBook.Price
	}
	if changedBook.Title != "" {
		book.Title = changedBook.Title
	}
	if changedBook.Pages != 0 {
		book.Pages = changedBook.Pages
	}
	if changedBook.Genres != nil {
		book.Genres = changedBook.Genres
	}
	err = writeData(books)
	return err
}

//PriceFilter returns filtered book objects
func PriceFilter(filter Filter) (Books, error) {
	var wantedBooks Books
	operator := string(filter.Price[0])
	if operator != "<" && operator != ">" {
		err := errors.New("unsupported operation")
		return nil, err
	}
	books, err := GetBooks()
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
