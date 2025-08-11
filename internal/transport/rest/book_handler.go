package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackietana/crud-app/internal/domain"
)

type BookService interface {
	GetBooks(ctx context.Context) ([]domain.Book, error)
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	CreateBook(ctx context.Context, book domain.BookCreate) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, id int, book domain.BookCreate) error
}

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(bookService BookService) *BookHandler {
	return &BookHandler{bookService}
}

func (bh *BookHandler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", bh.getBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", bh.getBookById).Methods(http.MethodGet)
		books.HandleFunc("", bh.createBook).Methods(http.MethodPost)
		books.HandleFunc("/{id:[0-9]+}", bh.deleteBook).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", bh.updateBook).Methods(http.MethodPut)
	}

	return r
}

func (bh *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := bh.bookService.GetBooks(context.TODO())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	log.Println("HANDLER: GET method in /books")
}

func (bh *BookHandler) getBookById(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := bh.bookService.GetBookById(context.TODO(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	log.Printf("HANDLER: GET method in /books/%d", id)
}

func (bh *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var book domain.BookCreate
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validateData(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.CreateBook(context.TODO(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book successfully created"))

	log.Println("HANDLER: POST method in /books")
}

func (bh *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.DeleteBook(context.TODO(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book successfully removed"))

	log.Printf("HANDLER: DELETE method in /books/%d", id)
}

func (bh *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var book domain.BookCreate
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validateData(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bh.bookService.UpdateBook(context.TODO(), id, book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book successfully updated"))

	log.Printf("HANDLER: PUT method in /books/%d", id)
}

func getId(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return 0, errors.New("ID not provided")
	}

	return strconv.Atoi(parts[2])
}

func validateData(b domain.BookCreate) error {
	if b.Name == "" || b.Description == "" || b.Author == "" || b.IsFree == nil || len(b.Genres) == 0 {
		return errors.New("all fields must be filled")
	}

	return nil
}
