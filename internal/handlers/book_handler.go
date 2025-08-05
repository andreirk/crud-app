package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackietana/crud-app/internal/models"
	"github.com/jackietana/crud-app/internal/repository"
)

func HandleBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBooks(db)(w, r)
		case http.MethodPost:
			createBook(db)(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func HandleBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBookById(db)(w, r)
		case http.MethodDelete:
			deleteBook(db)(w, r)
		case http.MethodPut:
			updateBook(db)(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func getBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := repository.GetBooks(db)
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
}

func getBookById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		book, err := repository.GetBookById(db, id)
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
}

func createBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var book models.BookCreate
		if err = json.Unmarshal(reqBytes, &book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateData(book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = repository.CreateBook(db, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Book successfully created"))

		log.Println("HANDLER: POST method in /books")
	}
}

func deleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = repository.DeleteBook(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Book successfully removed"))

		log.Printf("HANDLER: DELETE method in /books/%d", id)
	}
}

func updateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var book models.BookCreate
		if err = json.Unmarshal(reqBytes, &book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = validateData(book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = repository.UpdateBook(db, id, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Book successfully updated"))

		log.Printf("HANDLER: PUT method in /books/%d", id)
	}
}

func getId(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return 0, errors.New("ID not provided")
	}

	return strconv.Atoi(parts[2])
}

func validateData(b models.BookCreate) error {
	if b.Name == "" || b.Description == "" || b.Author == "" || b.IsFree == nil || len(b.Genres) == 0 {
		return errors.New("all fields must be filled")
	}

	return nil
}
