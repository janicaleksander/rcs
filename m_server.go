package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	s "github.com/janicaleksander/bcs/server"
	"github.com/joho/godotenv"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		s.Logger.Error("Error with loading .env file")
		return
	}

	dbase, err := db.NewPostgres(
		os.Getenv("DBNAME"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("SSLMODE"),
		db.WithConnectionTimeout(10))
	// in the future witt full two-side-ssl verification

	if err != nil {
		s.Logger.Error(err.Error())
		return
	}

	server := s.NewServer(os.Getenv("SERVER_ADDR"), dbase)
	r := remote.New(os.Getenv("SERVER_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		s.Logger.Error(err.Error())
		return
	}

	s.Logger.Info("server is running on: ", "Addr: ", os.Getenv("SERVER_ADDR"))
	e.Spawn(server, "server", actor.WithID("primary"))

	select {}
}
