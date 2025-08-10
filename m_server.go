package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	s "github.com/janicaleksander/bcs/server"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		utils.Logger.Error("Error with loading .env file")
		return
	}
	dbManager, err := db.GetDBManager(db.WithConnectionTimeout(10))
	if err != nil {
		utils.Logger.Error("Error with loading .env file")
		return
	}
	dbase := dbManager.GetDB()
	pg := &db.Postgres{Conn: dbase}

	// in the future witt full two-side-ssl verification

	server := s.NewServer(os.Getenv("SERVER_ADDR"), pg)
	r := remote.New(os.Getenv("SERVER_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	utils.Logger.Info("server is running on: ", "Addr: ", os.Getenv("SERVER_ADDR"))
	e.Spawn(server, "server", actor.WithID("primary"))

	select {}
}
