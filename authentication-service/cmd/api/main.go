package main

import (
	"authentication/data"
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

const webPort = "8082"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// decided to add smth new
	// ho-ho-ho. new funcs
	// really mate
	// let's do it
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
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

	// // Выполняем запрос к таблице users, чтобы получить строку с id = 2
	// var (
	// 	id       int
	// 	email    string
	// 	password string
	// )
	// query := `SELECT id, email, password FROM users WHERE id = $1`
	// err = db.QueryRow(query, 2).Scan(&id, &email, &password)
	// if err != nil {
	// 	return nil, err
	// }

	// // Выводим данные для проверки (можно удалить или заменить на другую логику)
	// log.Printf("User: ID=%d, Email=%s, Password=%s\n", id, email, password)

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	// dsn := "host=postgres port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
