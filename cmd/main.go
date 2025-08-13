package main

import (
	"log"

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
	r := bookHandler.InitRouter()
	log.Fatal(r.Run(":8080"))
}
