package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackietana/crud-app/internal/domain"
	log "github.com/sirupsen/logrus"

	_ "github.com/jackietana/crud-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// @Summary List books
// @Description get all books
// @Tags books
// @Produce json
// @Success 200 {object} []domain.Book
// @Router /books [get]
func (bh *BookHandler) getBooks(c *gin.Context) {
	books, err := bh.bookService.GetBooks(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBooks",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, books)

	log.WithField("handler", "getBooks").Info()
}

// @Summary Get specific book
// @Description get book by id
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [get]
func (bh *BookHandler) getBookById(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"issue":   "getId error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := bh.bookService.GetBookById(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
	log.WithField("handler", "getBookById").Info()
}

// @Summary Create book
// @Description create new book
// @Tags books
// @Accept json
// @Produce json
// @Success 201 {string} string "Book successfully created"
// @Router /books [post]
func (bh *BookHandler) createBook(c *gin.Context) {
	var book domain.Book
	if err := c.BindJSON(&book); err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"issue":   "bindJson error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err := bh.bookService.CreateBook(context.TODO(), book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.String(http.StatusCreated, "Book successfully created")
	log.WithField("handler", "createBook").Info()
}

// @Summary Delete book
// @Description delete book by id
// @Tags books
// @Produce plain
// @Param id path int true "Book ID"
// @Success 200 {string} string "Book successfully removed"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [delete]
func (bh *BookHandler) deleteBook(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"issue":   "getId error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.DeleteBook(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, "Book successfully removed")
	log.WithField("handler", "deleteBook").Info()
}

// @Summary Update book
// @Description update existing book
// @Tags books
// @Accept json
// @Produce plain
// @Param id path int true "Book ID"
// @Success 200 {string} string "Book successfully updated"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [put]
func (bh *BookHandler) updateBook(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"issue":   "getId error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	var book domain.Book
	if err := c.BindJSON(&book); err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"issue":   "bindJson error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.UpdateBook(context.TODO(), id, book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, "Book successfully updated")
	log.WithField("handler", "updateBook").Info()
}

func getId(c *gin.Context) (int, error) {
	id := c.Param("id")

	return strconv.Atoi(id)
}
