package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	log.Println("Starting mail service on port", webPort)

	app := Config{
		Mailer: Mail{
			Domain:      os.Getenv("MAIL_DOMAIN"),
			Host:        os.Getenv("MAIL_HOST"),
			Port:        first(strconv.Atoi(os.Getenv("MAIL_PORT"))),
			Username:    os.Getenv("MAIL_USERNAME"),
			Password:    os.Getenv("MAIL_PASSWORD"),
			FromName:    os.Getenv("MAIL_NAME"),
			FromAddress: os.Getenv("MAIL_ADDRESS"),
		},
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func first[T, U any](val T, _ U) T {
	return val
}
