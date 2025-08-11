package psql

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jackietana/crud-app/internal/domain"
	"github.com/lib/pq"
)

func GetBooks(db *sql.DB) ([]domain.Book, error) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]domain.Book, 0)

	for rows.Next() {
		b := domain.Book{}
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Author, &b.IsFree, pq.Array(&b.Genres), &b.PublishedAt); err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Println("REPO: SELECT * FROM books")

	return books, nil
}

func GetBookById(db *sql.DB, id int) (domain.Book, error) {
	var b domain.Book
	err := db.QueryRow("SELECT * FROM books WHERE id=$1", id).
		Scan(&b.ID, &b.Name, &b.Description, &b.Author, &b.IsFree, pq.Array(&b.Genres), &b.PublishedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return b, errors.New("book not found")
		}

		return b, err
	}

	log.Printf("REPO: SELECT * FROM books WHERE id=%d", id)

	return b, err
}

func CreateBook(db *sql.DB, b domain.BookCreate) error {
	strExec := "INSERT INTO books (name, description, author, is_free, genres) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Exec(strExec, b.Name, b.Description, b.Author, b.IsFree, pq.Array(b.Genres))

	log.Printf("INSERT INTO books (name, description, author, is_free, genres) VALUES (%s, %s, %s, %t, %s)",
		b.Name, b.Description, b.Author, *b.IsFree, b.Genres)

	return err
}

func DeleteBook(db *sql.DB, id int) error {
	exists, err := bookExistsByID(db, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.New("book not found")
	}

	_, err = db.Exec("DELETE FROM books WHERE id=$1", id)

	log.Printf("REPO: DELETE FROM books WHERE id=%d", id)

	return err
}

func UpdateBook(db *sql.DB, id int, b domain.BookCreate) error {
	exists, err := bookExistsByID(db, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.New("book not found")
	}

	strExec := "UPDATE books SET name=$1, description=$2, author=$3, is_free=$4, genres=$5 WHERE id=$6"
	_, err = db.Exec(strExec, b.Name, b.Description, b.Author, b.IsFree, pq.Array(b.Genres), id)

	log.Printf("REPO: UPDATE books SET name=%s, description=%s, author=%s, is_free=%t, genres=%s WHERE id=%d",
		b.Name, b.Description, b.Author, *b.IsFree, b.Genres, id)

	return err
}

func bookExistsByID(db *sql.DB, id int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
