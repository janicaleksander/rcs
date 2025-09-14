package main

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/messageservice"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func main() {
	config := struct {
		Messageservice struct {
			Addr       string `toml:"addr"`
			ServerAddr string `toml:"serverAddr"`
		}
	}{}
	_, err := toml.DecodeFile("configproduction/messageservice.toml", &config)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if len(strings.TrimSpace(config.Messageservice.Addr)) == 0 ||
		len(strings.TrimSpace(config.Messageservice.ServerAddr)) == 0 {
		utils.Logger.Error("bad server cfg file")
		return
	}
	r := remote.New(config.Messageservice.Addr, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error("Error with engine configuration")
		return
	}
	dbManager, err := db.GetDBManager(db.WithConnectionTimeout(10))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	dbase := dbManager.GetDB()
	pg := &db.Postgres{Conn: dbase}
	//TODO I don't know if i need a connection between MSSVC and server
	serverPID := actor.NewPID(config.Messageservice.ServerAddr, "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	utils.Logger.Info("messageservice is running on: ", "Address: ", config.Messageservice.Addr)
	messageService := messageservice.NewMessageService(serverPID, pg)
	e.Spawn(messageService, "messageService", actor.WithID("primary")) //this is creating new app

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	select {}
}
