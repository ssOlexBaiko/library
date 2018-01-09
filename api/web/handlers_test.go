package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	//"github.com/jinzhu/gorm"
	"github.com/ssOlexBaiko/library/storage"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/suite"
	"log"
	"os"
	"path/filepath"
)

// add flag for setting path to the storage and for using sql db
var (
	testLibPath = flag.String("libPath", "test_data/test_storage.json", "set path the storage file")
	sqlUse      = flag.Bool("sqlUse", false, "use sql db instead of json file")
)

func TestWeb(t *testing.T) {
	// hook for passing the same open file for test cases
	// and closing the file when all tests will pass
	path, err := filepath.Abs(*testLibPath)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(path, os.O_RDWR, 0660)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	testIndexHandler(t, file)
	testBooksIndexHandler(t, file)
	testGetBookHandler(t, file)
	testBookCreateHandler(t, file)
	testRemoveBookHandler(t, file)
	testChangeBookHandler(t, file)
	testBookFilterHandler(t, file)
}

func getRouter(file *os.File) *mux.Router {
	var store Storage
	flag.Parse()
	if *sqlUse {
		sqlStorage, err := storage.InitDB()
		if err != nil {
			log.Fatal(err)
		}

		store = storage.NewSQLLibrary(sqlStorage)
	} else {

		store = storage.NewLibrary(file)
	}

	router := NewRouter(
		NewHandler(store),
	)
	return router
}

func getTestBooks(t *testing.T, file *os.File) (storage.Books, error) {
	//t.Helper() //is available in go1.9 release
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		return nil, errors.New("BooksIndex handler returned wrong status code")
	}

	var books storage.Books
	err = json.NewDecoder(rr.Body).Decode(&books)
	if err != nil {
		return nil, errors.New("BooksIndex handler returned wrong data")
	}
	return books, nil
}

func testIndexHandler(t *testing.T, file *os.File) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, this is the library resource"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func testBooksIndexHandler(t *testing.T, file *os.File) {
	test := assert.New(t)
	_, err := getTestBooks(t, file)
	test.NoError(err, "test failed")
}

//type ExampleSuite struct {
//	suite.Suite
//	db *gorm.DB
//}

//func (s *ExampleSuite) BeforeTest() {
//	// TODO: !!!!!!!!!
//	s.db, _ = storage.InitDB()
//}

func testGetBookHandler(t *testing.T, file *os.File) {
	test := assert.New(t)
	books, err := getTestBooks(t, file)
	test.NoError(err, "test failed")

	url := "/books/" + books[0].ID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		test.FailNow(err.Error())
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)

	test.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	var book storage.Book
	err = json.NewDecoder(rr.Body).Decode(&book)
	test.NoError(err, "handler returned wrong data")

}

func testBookCreateHandler(t *testing.T, file *os.File) {
	testBook := storage.Book{
		Title:  "TestBook",
		Genres: []string{"test1", "test2"},
		Pages:  777,
		Price:  777,
	}

	book, err := json.Marshal(testBook)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewReader(book)
	req, err := http.NewRequest("POST", "/books", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check data!
	books, err := getTestBooks(t, file)
	addedBook := false
	for _, b := range books {
		if b.Title == testBook.Title {
			addedBook = true
		}
	}
	if !addedBook {
		t.Errorf("handler didn't add the book")
	}
}

func testRemoveBookHandler(t *testing.T, file *os.File) {
	books, err := getTestBooks(t, file)
	if err != nil {
		t.Errorf("test failed: %v", err)
	}

	url := "/books/" + books[len(books)-1].ID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func testChangeBookHandler(t *testing.T, file *os.File) {
	books, err := getTestBooks(t, file)
	if err != nil {
		t.Errorf("test failed: %v", err)
	}

	testBook := storage.Book{Title: "test", ID: books[0].ID}
	book, err := json.Marshal(testBook)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewReader(book)
	url := "/books/" + books[0].ID
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func testBookFilterHandler(t *testing.T, file *os.File) {
	price := storage.BookFilter{Price: "<77"}
	filter, err := json.Marshal(price)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewReader(filter)
	req, err := http.NewRequest("POST", "/books/filter", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter(file)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TODO:
//func TestBookFilterHandler(t *testing.T) {
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name string
//		args args
//		expectedBooks Books
//	}{
//		{"First", args{rr, req}}
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			BookFilterHandler(tt.args.w, tt.args.r)
//		})
//	}
//}
