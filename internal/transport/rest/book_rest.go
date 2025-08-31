package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackietana/crud-app/internal/domain"
	log "github.com/sirupsen/logrus"
)

// @Summary Create book
// @Description create new book
// @Tags books
// @Accept json
// @Produce json
// @Security TokenAuth
// @Success 201 {string} string "Book successfully created"
// @Router /books [post]
func (h *Handler) createBook(c *gin.Context) {
	var book domain.Book
	if err := c.BindJSON(&book); err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"issue":   "bindJson error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.bookService.CreateBook(context.TODO(), book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.String(http.StatusCreated, "Book successfully created")
	log.Info("Handler: createBook")
}

// @Summary Get specific book
// @Description get book by id
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Security TokenAuth
// @Success 200 {object} domain.Book
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [get]
func (h *Handler) getBookById(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"issue":   "getId error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetBookById(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
	log.Info("Handler: getBookById")
}

// @Summary List books
// @Description get all books
// @Tags books
// @Produce json
// @Security TokenAuth
// @Success 200 {object} []domain.Book
// @Router /books [get]
func (h *Handler) getBooks(c *gin.Context) {
	books, err := h.bookService.GetBooks(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBooks",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, books)

	log.Info("Handler: getBooks")
}

// @Summary Update book
// @Description update existing book
// @Tags books
// @Accept json
// @Produce plain
// @Param id path int true "Book ID"
// @Security TokenAuth
// @Success 200 {string} string "Book successfully updated"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [put]
func (h *Handler) updateBook(c *gin.Context) {
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

	err = h.bookService.UpdateBook(context.TODO(), id, book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, "Book successfully updated")
	log.Info("Handler: updateBook")
}

// @Summary Delete book
// @Description delete book by id
// @Tags books
// @Produce plain
// @Param id path int true "Book ID"
// @Security TokenAuth
// @Success 200 {string} string "Book successfully removed"
// @Failure 404 {string} string "book not found"
// @Router /books/{id} [delete]
func (h *Handler) deleteBook(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"issue":   "getId error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.bookService.DeleteBook(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"issue":   "service error",
		}).Error(err)
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, "Book successfully removed")
	log.Info("Handler: deleteBook")
}

func getId(c *gin.Context) (int, error) {
	id := c.Param("id")

	return strconv.Atoi(id)
}
