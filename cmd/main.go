package main

import (
	"log"
	"net/http"

	"github.com/jackietana/crud-app/internal/repository/psql"
	"github.com/jackietana/crud-app/internal/service"
	"github.com/jackietana/crud-app/internal/transport/rest"
	"github.com/jackietana/crud-app/pkg/database"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	//init dependencies
	bookRepo := psql.NewBookRepo(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := rest.NewBookHandler(bookService)

	//init and run server
	srv := &http.Server{
		Addr:    "8080",
		Handler: bookHandler.InitRouter(),
	}

	log.Println("Server started on :8080")
	log.Fatal(srv.ListenAndServe())
}
