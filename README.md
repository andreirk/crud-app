# crud-app
CRUD-application providing Web API to data in PostgreSQL.  
Data structure:
```json
{
    "id": 3,
    "name": "Book name",
    "description": "Book description",
    "author": "Book author",
    "is_free": false,
    "genres": ["Genre A", "Genre B"],
    "published_at": "2020-01-01T09:30:00.00000Z"
}
```
Endpoints: /books (GET and POST) and /books/id (GET, PUT and DELETE).
> /books GET: retrieve all available books  
> /books POST: create a new book  
> /books/id GET: retrieve a book by id  
> /books/id PUT: update an existing book by id  
> /books/id DELETE: delete an existing book by id

### Quick Start:
1. Install Go language
2. Install PostgreSQL and create a database
3. To create database run these commands in the terminal:
```bash
psql -U postgres
CREATE DATABASE db_name;
```
4. To quickly configure your DB, run this command in the terminal while in the root folder of the project:
```bash
psql -U postgres -d db_name -f migrations/001_create_users.sql
```
5. To install all dependencies run these commands:
```bash
go mod init your-app-name
go get github.com/lib/pq
```
6. To run the application run this command:
```bash
go run cmd/main.go
```
```go
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
```