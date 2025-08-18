package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackietana/cache-example"
	"github.com/jackietana/crud-app/internal/domain"
	log "github.com/sirupsen/logrus"
)

const ttlDuration = time.Hour * 8

var (
	allBooksCached bool
	cachedBookIDs  = make(map[string]string, 0)
)

type CacheHandler struct {
	cache *cache.Cache
}

func NewCacheHandler() *CacheHandler {
	return &CacheHandler{cache.New()}
}

func (ch *CacheHandler) GetCachedBooks() ([]domain.Book, error) {
	var books = make([]domain.Book, 0)

	if allBooksCached {
		for _, id := range cachedBookIDs {
			if id != "" {
				item, err := ch.cache.Get(id)
				if err != nil {
					return nil, err
				}

				if book, ok := item.(domain.Book); ok {
					books = append(books, book)
				}
			}
		}

		log.Info("Cacher: GetCachedBooks")
		return books, nil
	}

	return nil, errors.New("not all books are in cache")
}

func (ch *CacheHandler) AddBook(book domain.Book) {
	bookID := fmt.Sprintf("book_%d", book.ID)

	if item, _ := ch.cache.Get(bookID); item == nil {
		ch.cache.Set(bookID, book, ttlDuration)
		cachedBookIDs[bookID] = bookID
		log.WithField("id", book.ID).Info("Cacher: AddBook")
	}
}

func (ch *CacheHandler) AddBooks(books []domain.Book) {
	for _, book := range books {
		ch.AddBook(book)
	}

	allBooksCached = true
}

func (ch CacheHandler) GetCachedBook(id int) (domain.Book, error) {
	bookID := fmt.Sprintf("book_%d", id)

	if val, err := ch.cache.Get(bookID); err == nil {
		if book, ok := val.(domain.Book); ok {
			log.WithField("id", id).Info("Cacher: GetCachedBook")
			return book, nil
		}
	}

	return domain.Book{}, errors.New(bookID + "not found")
}

func (ch *CacheHandler) DeleteCachedBook(id int) {
	bookId := fmt.Sprintf("book_%d", id)

	if _, err := ch.cache.Get(bookId); err == nil {
		ch.cache.Delete(bookId)
		delete(cachedBookIDs, bookId)
		log.WithField("id", id).Info("Cacher: DeleteCachedBook")
	}
}

func (ch *CacheHandler) UpdateCachedBook(id int, book domain.Book) {
	bookId := fmt.Sprintf("book_%d", id)

	if val, err := ch.cache.Get(bookId); err == nil {
		if cachedBook, ok := val.(domain.Book); ok {
			book.ID = cachedBook.ID
			book.PublishedAt = cachedBook.PublishedAt

			ch.cache.Set(bookId, book, ttlDuration)
			log.WithField("id", id).Info("Cacher: UpdateCachedBook")
		}
	}
}

func (ch *CacheHandler) UpdateCacher() {
	allBooksCached = false
}
