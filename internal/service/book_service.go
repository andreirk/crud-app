package service

import (
	"context"

	"github.com/jackietana/crud-app/internal/domain"
)

type BooksRepository interface {
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	CreateBook(ctx context.Context, book domain.BookCreate) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, id int, book domain.BookCreate) error
}

type BooksService struct {
	repo BooksRepository
}

func NewBooksService(repo BooksRepository) *BooksService {
	return &BooksService{repo}
}

func (bs *BooksService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return bs.repo.GetBooks(ctx)
}

func (bs *BooksService) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	return bs.repo.GetBookById(ctx, id)
}

func (bs *BooksService) CreateBook(ctx context.Context, book domain.BookCreate) error {
	return bs.repo.CreateBook(ctx, book)
}

func (bs *BooksService) DeleteBook(ctx context.Context, id int) error {
	return bs.repo.DeleteBook(ctx, id)
}

func (bs *BooksService) UpdateBook(ctx context.Context, id int, book domain.BookCreate) error {
	return bs.repo.UpdateBook(ctx, id, book)
}
