package main

import (
	"log"
	"net/http"

	"github.com/jackietana/crud-app/internal/transport/rest"
	"github.com/jackietana/crud-app/pkg/database"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	//init handlers
	http.HandleFunc("/books", rest.HandleBooks(db))
	http.HandleFunc("/books/", rest.HandleBook(db))

	//init server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
