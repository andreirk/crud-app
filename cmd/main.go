package main

import (
	"fmt"
	"os"

	"github.com/jackietana/crud-app/internal/config"
	"github.com/jackietana/crud-app/internal/repository/psql"
	"github.com/jackietana/crud-app/internal/service"
	grpc_client "github.com/jackietana/crud-app/internal/transport/grpc"
	"github.com/jackietana/crud-app/internal/transport/rest"
	"github.com/jackietana/crud-app/pkg/database"
	"github.com/jackietana/crud-app/pkg/hash"
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
	userRepo := psql.NewUserRepo(db)
	tokenRepo := psql.NewTokenRepo(db)
	hasher := hash.NewSHA1Hasher(cfg.Salt)
	loggerClient, err := grpc_client.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}

	bookService := service.NewBookService(bookRepo)
	userService := service.NewUserService(userRepo, tokenRepo, hasher, loggerClient, cfg.Secret, cfg.Auth.TokenTTL, cfg.Auth.RefreshTTL)
	bookHandler := rest.NewHandler(bookService, userService)

	//init and run server
	r := bookHandler.InitRouter()
	log.Fatal(r.Run(fmt.Sprintf(":%d", cfg.Server.Port)))
}
