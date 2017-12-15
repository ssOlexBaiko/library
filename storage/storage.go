package storage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"github.com/twinj/uuid"
	"encoding/json"
)

type Book struct {
	ID	string 		`json:"id, omitempty"`
	Title	string		`json:"title, omitempty"`
	Genres	[]string	`json:"genres, omitempty"`
	Pages	int		`json:"pages, omitempty"`
	Price	float64		`json:"price, omitempty"`
}

type Books []Book

func GetBooks() (Books, error) {
	var books Books
	filepath, _ := filepath.Abs("storage/storage.json")
	file, err := ioutil.ReadFile(filepath)
    	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil, err
    	}
	err = json.Unmarshal(file, &books)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil, err
	}

	return books, nil
}

func CreateBook(book Book) error {
	filepath, _ := filepath.Abs("storage/storage.json")
	books, err := GetBooks()
	if err != nil {
		return err
	}
	book.ID = uuid.NewV4().String()
	books = append(books, book)
	booksBytes, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, booksBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
