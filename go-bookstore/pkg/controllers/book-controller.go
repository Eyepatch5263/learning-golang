package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"github.com/eyepatch5263/go-bookstore/pkg/models"
	"github.com/eyepatch5263/go-bookstore/pkg/utils"
)

var NewBook models.Book

// GetBooks handles GET requests to retrieve all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	newBooks:=models.GetAllBooks()
	res,_:=json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetBook handles GET requests to retrieve a book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the URL parameters
	vars:=mux.Vars(r)
	bookId:=vars["bookId"]
	// Convert the book ID from string to int64
	Id,err:=strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing book ID:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bookDetails, _ := models.GetBookById(Id)
	if bookDetails == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No book found with ID %d", Id)
		return
	}
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res) 
} 

// CreateBook handles POST requests to create a new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the book details
	var newBook models.Book
	err := utils.ParseBody(r, &newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create the book in the database
	b := newBook.CreateBook()
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// DeleteBook handles DELETE requests to remove a book by ID
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars:=mux.Vars(r)
	bookId:=vars["bookId"]
	Id,err:=strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing book ID:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book:=models.DeleteBook(Id)
	res,_:=json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UpdateBook handles PUT requests to update an existing book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	Id, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing book ID:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var book models.Book
	err = utils.ParseBody(r, &book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bookDetails, db := models.GetBookById(Id)
	if bookDetails == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No book found with ID %d", Id)
		return
	}
	if book.Title != "" {
		bookDetails.Title = book.Title
	}
	if book.Author != "" {
		bookDetails.Author = book.Author
	}
	if book.Publication != "" {
		bookDetails.Publication = book.Publication
	}
	if book.Price != "" {
		bookDetails.Price = book.Price
	}
	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}