package main

import (
	"database/sql"
	"fmt"
	"github.com/fertilewaif/avito-mx-backend-test/controllers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	PORT = 5432
)

func initDB(username, password, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@database:%d/%s?sslmode=disable",
		username, password, PORT, database)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	db, err := initDB(dbUser, dbPassword, dbName)

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	handler := controllers.NewSalesController(db)

	r.HandleFunc("/offers", handler.GetSales).Methods("GET")
	r.HandleFunc("/upload", handler.Upload).Methods("POST")
	r.HandleFunc("/get_status", handler.GetJobStatus).Methods("GET")

	http.Handle("/", r)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	handler.Close()
}
