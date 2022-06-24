package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/djedjethai/authentication/pkg/checking"
	"github.com/djedjethai/authentication/pkg/http/rest"
	"github.com/djedjethai/authentication/pkg/lib/config"
	"github.com/djedjethai/authentication/pkg/storage/postgres"
	// "github.com/djedjethai/authentication/pkg/storage/postgres"
)

const webPort = "80"

var app config.AppConfig

func main() {
	log.Println("starting authentication service")

	// start the config
	config := config.AppConfig{}

	dsnString := os.Getenv("DSN") // import from env
	fmt.Println("the dsn: ", dsnString)

	// set a pointer to the Storage struct(which have the DB connection)
	s, _ := NewStorage(dsnString)
	app = config

	repo := postgres.NewPostgresRepo(s.DB, &app)

	checker := checking.NewService(repo)

	mux := rest.Handler(checker)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
