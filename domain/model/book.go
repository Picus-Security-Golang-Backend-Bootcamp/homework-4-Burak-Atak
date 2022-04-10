package book

import (
	"errors"
	"strconv"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-Burak-Atak/helpers"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-Burak-Atak/infrastructure"
	"gorm.io/gorm"
)

type Books struct {
	gorm.Model
	BookName      string `json:"bookName"`
	StockCode     string `json:"stockCode"`
	Isbn          string `json:"isbn"`
	PageNumber    int    `json:"pageNumber"`
	Price         int    `json:"price"`
	StockQuantity int    `json:"stockQuantity"`
	Author        string `json:"author"`
}

type BookRepository struct {
	db *gorm.DB
}

var bookRepo *BookRepository

// Creates book repository and if db is empty reads csv file and creates books
func init() {
	db := infrastructure.NewPostgresDB("host=localhost user=postgres password=postgres dbname=library port=5432 sslmode=disable")
	bookRepo = NewRepository(db)
	bookRepo.Migration()

	allBooks := FindAll()
	if len(allBooks) == 0 {
		csvSlice := helpers.ReadCsv("books.csv")
		for _, c := range csvSlice[1:] {
			pagenumber, _ := strconv.Atoi(c[3])
			price, _ := strconv.Atoi(c[4])
			stockquantity, _ := strconv.Atoi(c[5])

			newBook := NewModel(c[0], c[1], c[2], pagenumber, price, stockquantity, c[6])
			Create(*newBook)
		}
	}
}

// Creates new book model
func NewModel(bookname string, stockcode string, isbn string, pagenumber int, price int, stockquantity int, author string) *Books {
	return &Books{
		BookName:      bookname,
		StockCode:     stockcode,
		Isbn:          isbn,
		PageNumber:    pagenumber,
		Price:         price,
		StockQuantity: stockquantity,
		Author:        author,
	}
}

// Creates book repositors
func NewRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

// Migration
func (r *BookRepository) Migration() {
	r.db.AutoMigrate(&Books{})
}

// Creates new book in database
func Create(book Books) {
	bookRepo.db.Create(&book)
}

// Deletes book in db
func Delete(b Books) {
	bookRepo.db.Delete(&b)
}

// Decrease stock quantity and saves new book
func Buy(b *Books, c int) {
	b.StockQuantity -= c
	bookRepo.db.Save(b)

}

// Finds all books in db
func FindAll() []Books {
	var books []Books
	bookRepo.db.Find(&books)

	return books
}

// Search the book by id and returns book if it exist
func SearchById(id string) (Books, error) {
	var book Books
	bookRepo.db.Find(&book, "id="+id)

	if book.BookName == "" {
		return book, errors.New("Book couldn't found")
	}
	return book, nil
}

// Searches in stock code, book name, isbn and author, returns all matched books
func SearchByInput(name string) []Books {
	var books []Books
	bookRepo.db.Find(&books,
		"LOWER(book_name) LIKE LOWER('%"+name+"%')"+
			"OR LOWER(stock_code) LIKE LOWER('%"+name+"%')"+
			"OR LOWER(isbn) LIKE LOWER('%"+name+"%')"+
			"OR LOWER(author) LIKE LOWER('%"+name+"%')")

	return books
}

// Updates book in db
func Update(b Books) {
	bookRepo.db.Save(&b)
}
