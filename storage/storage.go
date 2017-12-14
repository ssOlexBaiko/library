package storage

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"path/filepath"
)

type Book struct {
	ID	string
	Title	string
	Genres	[]string
	Pages	int
	Price	float64
}

type Books []Book

func InitDB() Books {
	books := Books{
		Book{
			ID:		"1349A807-87CA-446C-9740-480238489517",
			Title:		"Book title1",
			Genres:		[]string{"detective", "comedy"},
			Pages:		321,
			Price:		12.43,
		},
		Book{
			ID:		"C97376B9-6C2E-41E5-9DBE-2E82C0EF114B",
			Title:		"Book title2",
			Genres:		[]string{"adventure"},
			Pages:		234,
			Price:		25.43,
		},
		Book{
			ID:		"FFAD23EB-8FF4-4E09-82D2-AA33EBE3997F",
			Title:		"Book title3",
			Genres:		[]string{"historical"},
			Pages:		321,
			Price:		999.00,
		},
	}
	return books
}

func GetBooks() Books {
	var books Books
	filepath, _ := filepath.Abs("storage/storage.json")
	file, err := ioutil.ReadFile(filepath)
    	if err != nil {
        	fmt.Printf("File error: %v\n", err)
        	os.Exit(1)
    	}
	json.Unmarshal(file, &books)
	return books
}
