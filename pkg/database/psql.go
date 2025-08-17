package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jackietana/crud-app/internal/config"
	_ "github.com/lib/pq"
)

func ConnectDB(p *config.Postgres) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Pass, p.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfull connection to PostgreSQL")

	return db, nil
}
