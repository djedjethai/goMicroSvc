package main

import (
	"fmt"
	"github.com/djedjethai/mail/pkg/handler/handlehttp"
	"github.com/djedjethai/mail/pkg/sender"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = 80

func main() {
	log.Println("Starting mail service on port: ", webPort)

	// get storage

	// set the mailer
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := sender.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	// inj store in svc
	sender := sender.NewSender(&m)

	// inject svc in handler
	mux := handlehttp.Handler(sender)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
