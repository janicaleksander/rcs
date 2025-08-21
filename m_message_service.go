package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/messageservice"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Error("Can't load .env file")
		return
	}
	r := remote.New(os.Getenv("MESSAGE_SERVICE_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error("Error with engine configuration")
		return
	}
	dbManager, err := db.GetDBManager(db.WithConnectionTimeout(10))
	if err != nil {
		utils.Logger.Error("Error with loading .env file")
		return
	}
	dbase := dbManager.GetDB()
	pg := &db.Postgres{Conn: dbase}
	//TODO I don't know if i need a connection between MSSVC and server
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	utils.Logger.Info("messageservice is running on: ", "Address: ", os.Getenv("MESSAGE_SERVICE_ADDR"))
	messageService := messageservice.NewMessageService(serverPID, pg)
	e.Spawn(messageService, "messageService", actor.WithID("primary")) //this is creating new app

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	select {}
}
