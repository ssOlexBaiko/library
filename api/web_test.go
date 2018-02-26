package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/jinzhu/gorm"
	"github.com/ssOlexBaiko/library/storage"
	"github.com/ssOlexBaiko/library/api/web"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/suite"
	"github.com/gorilla/mux"
	"log"
)

// add flag for setting path to the storage and for using sql db
var (
	testLibPath = flag.String("libPath", "", "set path the storage file; default: test_data/test_storage.json")
	testSQLPath = flag.String("testSQLPath", "", "use sql db instead of json file; default: test_data/data.db")
)

func getRouter() *mux.Router {
	var store web.Storage
	var err error
	flag.Parse()
	if len(*testSQLPath) != 0  {
		store, err = storage.NewSQLLibrary(*testSQLPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		store, err = storage.NewLibrary(*testLibPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := web.NewRouter(
		web.NewHandler(store),
	)
	return router
}

func getTestBooks(t *testing.T) (storage.Books, error) {
	//t.Helper() //is available in go1.9 release
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()

	handler := getRouter()
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

// go test -bench=. api/web_test.go -cpuprofile=cpu.out -libPath=test_data/test_storage.json
// go tool pprof --svg api.test cpu.out > cpu1.svg
func BenchmarkGetTestBooks(b *testing.B) {
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		b.Fatal(err)
	}

	handler := getRouter()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func TestIndexHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter()
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

func TestBooksIndexHandler(t *testing.T) {
	test := assert.New(t)
	_, err := getTestBooks(t)
	test.NoError(err, "test failed")
}

//type ExampleSuite struct {
//	suite.Suite
//	db *gorm.DB
//}
//
//func (s *ExampleSuite) BeforeTest() {
//	// TODO: !!!!!!!!!
//	s.db, _ = storage.InitDB()
//}

func TestGetBookHandler(t *testing.T) {
	test := assert.New(t)
	books, err := getTestBooks(t)
	test.NoError(err, "test failed")

	url := "/books/" + books[0].ID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		test.FailNow(err.Error())
	}

	rr := httptest.NewRecorder()

	handler := getRouter()
	handler.ServeHTTP(rr, req)

	test.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	var book storage.Book
	err = json.NewDecoder(rr.Body).Decode(&book)
	test.NoError(err, "handler returned wrong data")

}

func TestBookCreateHandler(t *testing.T) {
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

	handler := getRouter()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check data!
	books, err := getTestBooks(t)
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

func TestRemoveBookHandler(t *testing.T) {
	books, err := getTestBooks(t)
	if err != nil {
		t.Errorf("test failed: %v", err)
	}

	url := "/books/" + books[len(books)-1].ID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := getRouter()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestChangeBookHandler(t *testing.T) {
	books, err := getTestBooks(t)
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

	handler := getRouter()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestBookFilterHandler(t *testing.T) {
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

	handler := getRouter()
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
