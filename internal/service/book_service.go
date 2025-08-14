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

var (
	allBooksCached bool
	cachedBookIDs  = make(map[string]string, 0)
)

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
	var books = make([]domain.Book, 0)

	if allBooksCached {
		for _, id := range cachedBookIDs {
			if id != "" {
				item, _ := bs.cache.Get(id)
				books = append(books, item.(domain.Book))
			}
		}

		log.Println("SERVICE: all cached books retrieved")
		return books, nil
	}

	books, err := bs.repo.GetBooks(ctx)

	for _, book := range books {
		bookID := fmt.Sprintf("book_%d", book.ID)

		if item, _ := bs.cache.Get(bookID); item == nil {
			bs.cache.Set(bookID, book, ttlDuration)
			cachedBookIDs[bookID] = bookID
			log.Printf("SERVICE: %s added to cache", bookID)
		}
	}

	allBooksCached = true

	return books, err
}

func (bs *BookService) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	bookID := fmt.Sprintf("book_%d", id)

	if val, err := bs.cache.Get(bookID); err == nil {
		if book, ok := val.(domain.Book); ok {
			log.Printf("SERVICE: %s retrieved from cache", bookID)
			return book, nil
		}
	}

	book, err := bs.repo.GetBookById(ctx, id)
	if err == nil {
		bs.cache.Set(bookID, book, ttlDuration)
		cachedBookIDs[bookID] = bookID
		log.Printf("SERVICE: book_%d added to cache", id)
	}

	return book, err
}

func (bs *BookService) CreateBook(ctx context.Context, book domain.Book) error {
	allBooksCached = false
	return bs.repo.CreateBook(ctx, book)
}

func (bs *BookService) DeleteBook(ctx context.Context, id int) error {
	bookId := fmt.Sprintf("book_%d", id)

	if _, err := bs.cache.Get(bookId); err == nil {
		bs.cache.Delete(bookId)
		delete(cachedBookIDs, bookId)
		log.Printf("SERVICE: %s removed from cache", bookId)
	}

	return bs.repo.DeleteBook(ctx, id)
}

func (bs *BookService) UpdateBook(ctx context.Context, id int, book domain.Book) error {
	bookId := fmt.Sprintf("book_%d", id)

	if _, err := bs.cache.Get(bookId); err == nil {
		bs.cache.Delete(bookId)
		allBooksCached = false
		log.Printf("SERVICE: %s removed from cache", bookId)
	}

	return bs.repo.UpdateBook(ctx, id, book)
}
