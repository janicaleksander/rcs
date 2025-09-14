package main

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/external/connector"
	"github.com/janicaleksander/bcs/external/connector/api"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func main() {
	//Setup remote access

	config := struct {
		Connector struct {
			AddrDevice string `toml:"addrDevice"`
			AddrHTTP   string `toml:"addrHTTP"`
			ServerAddr string `toml:"serverAddr"`
		}
	}{}
	_, err := toml.DecodeFile("configproduction/connector.toml", &config)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if len(strings.TrimSpace(config.Connector.AddrDevice)) == 0 ||
		len(strings.TrimSpace(config.Connector.AddrHTTP)) == 0 ||
		len(strings.TrimSpace(config.Connector.ServerAddr)) == 0 {
		utils.Logger.Error("bad device cfg file")
		return
	}

	r := remote.New(config.Connector.AddrDevice, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	utils.Logger.Info("device is running on:", "Addr:", config.Connector.AddrDevice)
	serverPID := actor.NewPID(config.Connector.ServerAddr, "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	dActor := connector.NewServiceDeviceActor()
	pid := e.Spawn(dActor, "device")
	resp = e.Request(pid, actor.Context{}, utils.WaitTime)
	res, err := resp.Result()
	if err != nil {
		utils.Logger.Error("Cant get context")
		return
	}
	var ctx *actor.Context
	if v, ok := res.(*actor.Context); ok {
		ctx = v
	}
	httpHandler := api.NewHandler(config.Connector.AddrHTTP, ctx, serverPID)
	hhttp := connector.NewServiceHTTPDevice(httpHandler)
	go hhttp.RunHTTPServer()
	select {}
}
