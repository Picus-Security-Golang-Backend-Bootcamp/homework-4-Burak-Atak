package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	bookModel "github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-Burak-Atak/domain/model"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-Burak-Atak/helpers"
	"github.com/gorilla/mux"
)

// List all books
func List(w http.ResponseWriter, r *http.Request) {

	books := bookModel.FindAll()
	w.Header().Set("Content-type", "application/slice")

	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		output := []byte("Couldn't found any book")
		w.Write(output)
	} else {
		resp, _ := json.Marshal(books)
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// Search a book
func Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	books := bookModel.SearchByInput(vars["query"])
	w.Header().Set("Content-type", "application/slice")
	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		output := []byte("Couldn't found any book")
		w.Write(output)
	} else {
		resp, _ := json.Marshal(books)
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// Buy a book
func Buy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	number, err2 := strconv.Atoi(vars["count"])
	if err != nil || err2 != nil {
		output := []byte("Wrong input.")
		w.Write(output)
	} else {
		book, err := bookModel.SearchById(vars["id"])

		if err != nil {
			res := fmt.Sprintf("%v", err)
			output := []byte(res)
			w.Write(output)
		} else {
			if number > 0 && number <= book.StockQuantity {
				bookModel.Buy(&book, number)
				res := fmt.Sprintf("You bought %d books. New stock is %d", number, book.StockQuantity)
				output := []byte(res)
				w.Write(output)
			} else {
				res := fmt.Sprintf("You can't buy %d books. Only %d books are available.", number, book.StockQuantity)
				output := []byte(res)
				w.Write(output)
				fmt.Printf("Enter a number beetween 0-%d", book.StockQuantity)
			}
		}
	}
}

// Delete a book
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		output := []byte("Wrong input.")
		w.Write(output)
	} else {
		book, err := bookModel.SearchById(vars["id"])
		if err != nil {
			res := fmt.Sprintf("%v", err)
			output := []byte(res)
			w.Write(output)
		} else {
			bookModel.Delete(book)
			output := []byte("You deleted the book.")
			w.Write(output)
		}
	}
}

// Create a new book
func Create(w http.ResponseWriter, r *http.Request) {
	var book bookModel.Books

	err := helpers.DecodeJSONBody(w, r, &book)
	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	bookModel.Create(book)
}
