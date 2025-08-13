package service

import (
	"context"

	"github.com/jackietana/crud-app/internal/domain"
)

type BookRepository interface {
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	CreateBook(ctx context.Context, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, id int, book domain.Book) error
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{repo}
}

func (bs *BookService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return bs.repo.GetBooks(ctx)
}

func (bs *BookService) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	return bs.repo.GetBookById(ctx, id)
}

func (bs *BookService) CreateBook(ctx context.Context, book domain.Book) error {
	return bs.repo.CreateBook(ctx, book)
}

func (bs *BookService) DeleteBook(ctx context.Context, id int) error {
	return bs.repo.DeleteBook(ctx, id)
}

func (bs *BookService) UpdateBook(ctx context.Context, id int, book domain.Book) error {
	return bs.repo.UpdateBook(ctx, id, book)
}
