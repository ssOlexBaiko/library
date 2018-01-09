package storage

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type sqlLibrary struct {
	sqlStorage *gorm.DB
}

func NewSQLLibrary(sqlStorage *gorm.DB) *sqlLibrary {
	return &sqlLibrary{
		sqlStorage,
	}
}

func (l *sqlLibrary) GetBooks() (Books, error) {
	var books Books
	// SELECT * FROM books
	return books, l.sqlStorage.Find(&books).Error
}

func (l *sqlLibrary) CreateBook(book Book) error {
	book.PrepareToCreate()
	return l.sqlStorage.Create(&book).Error
}

func (l *sqlLibrary) GetBook(id string) (Book, error) {
	var book Book

	err := l.sqlStorage.Where("id = ?", id).First(&book).Error
	if err == gorm.ErrRecordNotFound {
		return book, ErrNotFound
	}

	return book, err
}

func (l *sqlLibrary) RemoveBook(id string) error {
	err := l.sqlStorage.Where("id = ?", id).Delete(&Book{}).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	return err
}

func (l *sqlLibrary) ChangeBook(changedBook Book) (Book, error) {
	var book Book
	err := l.sqlStorage.Update(&changedBook).Error
	if err == gorm.ErrRecordNotFound {
		return book, ErrNotFound
	}

	return changedBook, err
}

func (l *sqlLibrary) PriceFilter(filter BookFilter) (Books, error) {
	var books Books
	if len(filter.Price) <= 1 {
		return nil, ErrNotValidData
	}
	operator := string(filter.Price[0])
	if operator != "<" && operator != ">" {
		return nil, ErrUnsupportedOperation
	}

	price, err := strconv.ParseFloat(filter.Price[1:], 64)
	if err != nil {
		return nil, err
	}
	if operator == ">" {
		err := l.sqlStorage.Find(&books, "price > ?", price).Error
		if err != nil {
			return nil, err
		} else {
			err := l.sqlStorage.Find(&books, "price < ?", price).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return books, err
}
