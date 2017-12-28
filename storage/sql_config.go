package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*gorm.DB, error) {
	// Opening file
	db, err := gorm.Open("sqlite3", "storage/data.db")
	if err != nil {
		return nil, err
	}
	// Display SQL queries
	if err = db.LogMode(true).Error; err != nil {
		return nil, err
	}
	// Creating the table
	if !db.HasTable(&Book{}) {
		if err = db.CreateTable(&Book{}).Set("gorm:table_options", "ENGINE=InnoDB").Error; err != nil {
			return nil, err
		}
	}

	db = db.Debug()

	return db, nil
}
