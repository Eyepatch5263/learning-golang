package models

import (
	"github.com/jinzhu/gorm"
	"github.com/eyepatch5263/go-bookstore/pkg/config"
	"fmt"
)

var db* gorm.DB

type Book struct {
	gorm.Model
	
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
	Price       string `json:"price"`
}

// Initialize the database connection and migrate the Book model
func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
    result := db.Create(&b)
    if result.Error != nil {
        fmt.Printf("Error creating book: %v\n", result.Error)
    }
    return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var book Book
	db:=db.Where("ID=?",Id).Find(&book)
	return &book, db
}

func DeleteBook(Id int64) Book {
	var book Book
	db.Where("ID=?", Id).Delete(&book)
	return book 
}