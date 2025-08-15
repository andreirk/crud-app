package service

import (
	"context"

	"github.com/jackietana/crud-app/internal/domain"
	"github.com/jackietana/crud-app/pkg"
)

type BookRepository interface {
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	CreateBook(ctx context.Context, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, id int, book domain.Book) error
}

type BookService struct {
	repo   BookRepository
	cacher *pkg.CacheHandler
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{repo, pkg.NewCacheHandler()}
}

func (bs *BookService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	books, err := bs.cacher.GetCachedBooks()
	if err == nil {
		return books, err
	}

	books, err = bs.repo.GetBooks(ctx)
	bs.cacher.AddBooks(books)

	return books, err
}

func (bs *BookService) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	book, err := bs.cacher.GetCachedBook(id)
	if err == nil {
		return book, err
	}

	book, err = bs.repo.GetBookById(ctx, id)
	if err == nil {
		bs.cacher.AddBook(book)
	}

	return book, err
}

func (bs *BookService) CreateBook(ctx context.Context, book domain.Book) error {
	bs.cacher.UpdateCacher()

	return bs.repo.CreateBook(ctx, book)
}

func (bs *BookService) DeleteBook(ctx context.Context, id int) error {
	bs.cacher.DeleteCachedBook(id)

	return bs.repo.DeleteBook(ctx, id)
}

func (bs *BookService) UpdateBook(ctx context.Context, id int, book domain.Book) error {
	bs.cacher.UpdateCachedBook(id)

	return bs.repo.UpdateBook(ctx, id, book)
}
