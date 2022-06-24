package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/djedjethai/broker/pkg/domain/jsdomain"
	"github.com/djedjethai/broker/pkg/http/rest"
	"github.com/djedjethai/broker/pkg/service"
)

const webPort = "80"

// type Config struct{}

func main() {
	// app := Config{}

	log.Println("Starting broker service on port: ", webPort)

	// get the memory datas(need to create them)
	var svc service.Service

	repo := new(jsdomain.Repository)
	svc = service.NewService(repo)

	// populate the repo
	svc.AddData(sampleData)

	mux := rest.Handler(svc)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
