package main

import (
	"log"
	"os"

	"github.com/janicaleksander/bcs/external/devicehttp"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
	d := devicehttp.New(os.Getenv("DEVICEHTTP_ADDR"))
	d.RunHTTPServer()
}
