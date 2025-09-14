package main

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	s "github.com/janicaleksander/bcs/server"
	"github.com/janicaleksander/bcs/utils"
)

func main() {
	config := struct {
		Server struct {
			Addr string `toml:"addr"`
		}
	}{}
	_, err := toml.DecodeFile("configproduction/server.toml", &config)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if len(strings.TrimSpace(config.Server.Addr)) == 0 {
		utils.Logger.Error("bad server cfg file")
		return
	}
	dbManager, err := db.GetDBManager(db.WithConnectionTimeout(10))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	dbase := dbManager.GetDB()
	pg := &db.Postgres{Conn: dbase}

	// in the future witt full two-side-ssl verification

	server := s.NewServer(pg)
	r := remote.New(config.Server.Addr, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	utils.Logger.Info("server is running on: ", "Addr: ", config.Server.Addr)
	e.Spawn(server, "server", actor.WithID("primary"))

	select {}
}
