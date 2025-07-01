package main

import (
	db "github.com/janicaleksander/bcs/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	_, err = db.NewPostgres(
		os.Getenv("DBNAME"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("SSLMODE"),
		db.WithConnectionTimeout(10))
	// in the future witt full two-side-ssl verification

	if err != nil {
		panic(err)
	}

}
