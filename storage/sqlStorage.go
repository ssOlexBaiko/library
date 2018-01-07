package storage

import (
	"errors"
	"github.com/jinzhu/gorm"
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

func (l *sqlLibrary) ChangeBook(changedBook Book) error {
	err := l.sqlStorage.Update(&changedBook).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	return err
}

func (l *sqlLibrary) PriceFilter(filter BookFilter) (Books, error) {
	return nil, errors.New("NotImplemented")
}
