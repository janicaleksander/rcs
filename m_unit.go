package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/external/unit"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Cant load .evn file")
		return
	}
	unitAddrFlag := flag.String("address", "", "Type here IP address of application")
	unitIDFlag := flag.String("unitID", "", "Type here ID unit to what you want to connect")
	flag.Parse()
	if len(strings.TrimSpace(*unitAddrFlag)) <= 0 || len(strings.TrimSpace(*unitIDFlag)) <= 0 {
		utils.Logger.Error("Type value of flag")
		return
	}
	//Setup remote access
	r := remote.New(*unitAddrFlag, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	utils.Logger.Info("unit is running on:", "Addr:", *unitAddrFlag)
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	dbManager, err := db.GetDBManager(db.WithConnectionTimeout(10))
	if err != nil {
		utils.Logger.Error("Error with loading .env file")
		return
	}
	dbase := dbManager.GetDB()
	pg := &db.Postgres{Conn: dbase}
	_ = e.Spawn(unit.NewUnit(*unitIDFlag, serverPID, pg), "device", actor.WithID(*unitIDFlag))
	select {}
}
