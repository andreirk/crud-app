package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackietana/cache-example"
	"github.com/jackietana/crud-app/internal/domain"
)

const ttlDuration = time.Hour * 8

type BookRepository interface {
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	CreateBook(ctx context.Context, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
	UpdateBook(ctx context.Context, id int, book domain.Book) error
}

type BookService struct {
	repo  BookRepository
	cache *cache.Cache
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{repo, cache.New()}
}

func (bs *BookService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return bs.repo.GetBooks(ctx)
}

func (bs *BookService) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	// in cache? (Y -> return) -> check db -> add in cache
	if val, err := bs.cache.Get(fmt.Sprintf("book_%d", id)); err == nil {
		if book, ok := val.(domain.Book); ok {
			log.Printf("SERVICE: book_%d retrieved from cache", id)
			return book, nil
		}
	}

	book, err := bs.repo.GetBookById(ctx, id)
	if err == nil {
		bs.cache.Set(fmt.Sprintf("book_%d", id), book, ttlDuration)
		log.Printf("SERVICE: book_%d added to cache", id)
	}

	return book, err
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
