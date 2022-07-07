package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "3002"
const maxTries = 16

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service on port", port)

	conn := connectToDB()

	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	tries := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet")
			tries += 1
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if tries >= maxTries {
			log.Panic("Could not connect to Postgres")
		}
	}
}
