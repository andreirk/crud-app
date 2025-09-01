package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackietana/crud-app/internal/domain"

	_ "github.com/jackietana/crud-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type BookService interface {
	CreateBook(ctx context.Context, book domain.Book) error
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	UpdateBook(ctx context.Context, id int, book domain.Book) error
	DeleteBook(ctx context.Context, id int) error
}

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, user domain.UserSignIn) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	bookService BookService
	userService UserService
}

func NewHandler(bookService BookService, userService UserService) *Handler {
	return &Handler{bookService, userService}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(loggerMiddleware())

	{
		auth := r.Group("/auth")
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
	}

	{
		books := r.Group("/books")
		books.Use(h.authMiddleware())
		books.POST("", h.createBook)
		books.GET("/:id", h.getBookById)
		books.GET("", h.getBooks)
		books.PUT("/:id", h.updateBook)
		books.DELETE("/:id", h.deleteBook)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
