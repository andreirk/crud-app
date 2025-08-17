package main

import (
	"fmt"
	"os"

	"github.com/jackietana/crud-app/internal/config"
	"github.com/jackietana/crud-app/internal/repository/psql"
	"github.com/jackietana/crud-app/internal/service"
	"github.com/jackietana/crud-app/internal/transport/rest"
	"github.com/jackietana/crud-app/pkg/database"
	log "github.com/sirupsen/logrus"
)

const (
	CONF_DIR  = "configs"
	CONF_FILE = "main"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// @title CRUD-app
// @version 1.0
// @description CRUD-application providing Web API to data in PostgreSQL.

// @host localhost:8080
// @BasePath /

func main() {
	// init configuration
	cfg, err := config.New(CONF_DIR, CONF_FILE)
	if err != nil {
		log.Fatal(err)
	}

	// init db connection
	db, err := database.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//init dependencies
	bookRepo := psql.NewBookRepo(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := rest.NewBookHandler(bookService)

	//init and run server
	r := bookHandler.InitRouter()
	log.Fatal(r.Run(fmt.Sprintf(":%d", cfg.Server.Port)))
}
