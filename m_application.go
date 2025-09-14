package main

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/application"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func main() {
	config := struct {
		Application struct {
			Addr               string `toml:"addr"`
			MessageserviceAddr string `toml:"messageserviceAddr"`
			ServerAddr         string `toml:"serverAddr"`
		}
	}{}
	_, err := toml.DecodeFile("configproduction/application.toml", &config)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if len(strings.TrimSpace(config.Application.Addr)) == 0 ||
		len(strings.TrimSpace(config.Application.MessageserviceAddr)) == 0 ||
		len(strings.TrimSpace(config.Application.ServerAddr)) == 0 {
		utils.Logger.Error("bad application cfg file")
		return
	}
	//Setup remote access
	r := remote.New(config.Application.Addr, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	utils.Logger.Info("application is running on:", "Addr:", config.Application.Addr)
	messageServicePID := actor.NewPID(config.Application.MessageserviceAddr, "messageService/primary")
	serverPID := actor.NewPID(config.Application.ServerAddr, "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	window := application.NewWindow()
	app := application.NewWindowActor(window, serverPID, messageServicePID)
	e.Spawn(app, "app") //this is creating new app

	//window
	window.RunWindow()

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	<-window.Done
}
