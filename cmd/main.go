package main

import (
	"log"
	"net/http"

	"github.com/jackietana/crud-app/internal/database"
	"github.com/jackietana/crud-app/internal/handlers"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	//init handlers
	http.HandleFunc("/books", handlers.HandleBooks(db))
	http.HandleFunc("/books/", handlers.HandleBook(db))

	//init server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
