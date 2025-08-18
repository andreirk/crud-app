package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackietana/crud-app/internal/domain"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepository {
	return &BookRepository{db}
}

func (br *BookRepository) GetBooks(ctx context.Context) ([]domain.Book, error) {
	rows, err := br.db.QueryContext(ctx, "SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]domain.Book, 0)

	for rows.Next() {
		b := domain.Book{}
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Author, &b.IsFree, pq.Array(&b.Genres),
			&b.PublishedAt); err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Info("Repository: GetBooks")

	return books, nil
}

func (br *BookRepository) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	var b domain.Book
	err := br.db.QueryRowContext(ctx, "SELECT * FROM books WHERE id=$1", id).
		Scan(&b.ID, &b.Name, &b.Description, &b.Author, &b.IsFree, pq.Array(&b.Genres), &b.PublishedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return b, errors.New("book not found")
		}

		return b, err
	}

	log.WithField("id", id).Info("Repository: GetBookById")

	return b, err
}

func (br *BookRepository) CreateBook(ctx context.Context, b domain.Book) error {
	strExec := "INSERT INTO books (name, description, author, is_free, genres) VALUES ($1, $2, $3, $4, $5)"
	_, err := br.db.ExecContext(ctx, strExec, b.Name, b.Description, b.Author, b.IsFree, pq.Array(b.Genres))

	log.Info("Repository: CreateBook")

	return err
}

func (br *BookRepository) DeleteBook(ctx context.Context, id int) error {
	exists, err := br.bookExistsByID(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.New("book not found")
	}

	_, err = br.db.ExecContext(ctx, "DELETE FROM books WHERE id=$1", id)

	log.WithField("id", id).Info("Repository: DeleteBook")

	return err
}

func (br *BookRepository) UpdateBook(ctx context.Context, id int, b domain.Book) error {
	exists, err := br.bookExistsByID(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.New("book not found")
	}

	strExec := "UPDATE books SET name=$1, description=$2, author=$3, is_free=$4, genres=$5 WHERE id=$6"
	_, err = br.db.ExecContext(ctx, strExec, b.Name, b.Description, b.Author, b.IsFree, pq.Array(b.Genres), id)

	log.WithField("id", id).Info("Repository: UpdateBook")

	return err
}

func (br *BookRepository) bookExistsByID(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := br.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM books WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
