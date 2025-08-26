package main

import (
	"log"
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/external/deviceservice"
	"github.com/janicaleksander/bcs/external/deviceservice/api"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Cant load .evn file")
		return
	}

	//Setup remote access
	r := remote.New(os.Getenv("DEVICE_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	utils.Logger.Info("device is running on:", "Addr:", os.Getenv("DEVICE_ADDR"))
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	dActor := deviceservice.NewDeviceActor()
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
	httpHandler := api.NewHandler(os.Getenv("HTTP_ADDR"), ctx, serverPID)
	hhttp := deviceservice.NewHTTPDevice(httpHandler)
	go hhttp.RunHTTPServer()
	select {}
}
