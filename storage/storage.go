package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/twinj/uuid"
	"encoding/json"
	"log"
)

type Book struct {
	ID	string
	Title	string
	Genres	[]string
	Pages	int
	Price	float64
}

type Books []Book

func GetBooks() []byte {
	filepath, _ := filepath.Abs("storage/storage.json")
	books, err := ioutil.ReadFile(filepath)
    	if err != nil {
        	fmt.Printf("File error: %v\n", err)
        	os.Exit(1)
    	}

	return books
}

func CreateBook(book Book) error {
	filepath, _ := filepath.Abs("storage/storage.json")
	book.ID = fmt.Sprint(uuid.NewV4())
	bookJson, _ := json.Marshal(book)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(bookJson); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return nil
}
