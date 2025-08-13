package rest

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackietana/crud-app/internal/domain"
)

type BookService interface {
	GetBooks(ctx context.Context) ([]domain.Book, error)
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	CreateBook(ctx context.Context, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, id int, book domain.Book) error
}

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(bookService BookService) *BookHandler {
	return &BookHandler{bookService}
}

func (bh *BookHandler) InitRouter() *gin.Engine {
	r := gin.Default()

	{
		books := r.Group("/books")
		books.GET("", bh.getBooks)
		books.GET("/:id", bh.getBookById)
		books.POST("", bh.createBook)
		books.DELETE("/:id", bh.deleteBook)
		books.PUT("/:id", bh.updateBook)
	}

	return r
}

func (bh *BookHandler) getBooks(c *gin.Context) {
	books, err := bh.bookService.GetBooks(context.TODO())
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, books)

	log.Println("HANDLER: GET method in /books")
}

func (bh *BookHandler) getBookById(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := bh.bookService.GetBookById(context.TODO(), id)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
	log.Printf("HANDLER: GET method in /books/%d", id)
}

func (bh *BookHandler) createBook(c *gin.Context) {
	var book domain.Book
	if err := c.BindJSON(&book); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err := bh.bookService.CreateBook(context.TODO(), book)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.String(http.StatusCreated, "Book successfully created")
	log.Println("HANDLER: POST method in /books")
}

func (bh *BookHandler) deleteBook(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.DeleteBook(context.TODO(), id)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, "Book successfully removed")
	log.Printf("HANDLER: DELETE method in /books/%d", id)
}

func (bh *BookHandler) updateBook(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	var book domain.Book
	if err := c.BindJSON(&book); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.UpdateBook(context.TODO(), id, book)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, "Book successfully updated")
	log.Printf("HANDLER: PUT method in /books/%d", id)
}

func getId(c *gin.Context) (int, error) {
	id := c.Param("id")

	return strconv.Atoi(id)
}
