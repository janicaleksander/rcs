package server

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO make a two maybe maps one for connected by app commander lvl person
// and second to 0 1 2 users/soldiers connected by device->unit->server
// this is problem e.g. in sending message from 5lvl from app to 0lvl to device in unit
// but this 5 see this 0 in his lists of person in system (maybe)
type Server struct {
	storage            db.Storage
	connections        map[string]*actor.PID // uuid to PID
	reverseConnections map[string]string     //PIDstring to -> uuid
}

func NewServer(storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:            storage,
			connections:        make(map[string]*actor.PID),
			reverseConnections: make(map[string]string),
		}
	}
}

// TODO In the future legal way to disconnect
// TODO manage if err !=nil some log errors
func (s *Server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("server initialized")
	case actor.Started:
		utils.Logger.Info("server has started")
		ctx.SendRepeat(ctx.PID(), &proto.HeartbeatTick{}, 10*time.Second)
		//safety thing if user is not in unit but have a unit id -> error
	case actor.Stopped:
		utils.Logger.Info("server has stopped")
	case *proto.HeartbeatTick:
		s.startHeartbeat(ctx)
	case *proto.IsServerRunning:
		ctx.Respond(&proto.IsServerRunning{})
	case *proto.Disconnect: // after this switch state to loginScene
		pid, ok := s.connections[msg.Id]
		if ok {
			delete(s.connections, msg.Id)
			delete(s.reverseConnections, pid.String())
		}

	case *proto.IsOnline:
		if _, ok := s.connections[msg.UserID]; ok {
			ctx.Respond(&proto.Online{})
		} else {
			ctx.Respond(&proto.Offline{})
		}
	case *proto.LoginUnit:
		pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //unit PID
		s.connections[msg.UnitID] = pid                  //pid to uuid
		s.reverseConnections[pid.String()] = msg.UnitID
		fmt.Println(s.connections)
	case *proto.LoginUser:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		id, role, err := s.storage.LoginUser(c, msg.Email, msg.Password)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
			s.connections[id] = pid                          //pid to uuid
			s.reverseConnections[pid.String()] = id
			ctx.Respond(&proto.AcceptUserLogin{UserID: id, RuleLevel: int32(role)})
		}
		//TODO idk if this getlogged works

	case *proto.GetLoggedInUUID: //returning an uuid of current logged-in user
		pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
		id := s.reverseConnections[pid.String()]
		ctx.Respond(&proto.LoggedInUUID{Id: id})
	case *proto.GetUserAboveLVL:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		users, err := s.storage.GetUsersWithLVL(c, int(msg.Lower), int(msg.Upper))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UsersAboveLVL{Users: users})

		}
	case *proto.CreateUnit:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.storage.InsertUnit(c, msg.Unit, msg.UserID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptCreateUnit{})
		}
	case *proto.GetAllUnits:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		units, err := s.storage.GetAllUnits(c)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AllUnits{Units: units})

		}

	case *proto.GetUsersInUnit:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		unitID := msg.UnitID
		users, err := s.storage.GetUsersInUnit(c, unitID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UsersInUnit{Users: users})

		}
	case *proto.CreateUser:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.storage.InsertUser(c, msg.User)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptCreateUser{})
		}
	case *proto.IsUserInUnit:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		isInUnit, unitID, err := s.storage.IsUserInUnit(c, msg.Id)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else if isInUnit {
			ctx.Respond(&proto.UserIsInUnit{UnitID: unitID})
		} else {
			ctx.Respond(&proto.Error{Content: err.Error()})
		}
	case *proto.AssignUserToUnit:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.storage.AssignUserToUnit(c, msg.UserID, msg.UnitID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptAssignUserToUnit{})
		}
	case *proto.DeleteUserFromUnit:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.storage.DeleteUserFromUnit(c, msg.UserID, msg.UnitID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptDeleteUserFromUnit{})
		}
	case *proto.HTTPSpawnDevice:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		userID, _, err := s.storage.LoginUser(c, msg.Email, msg.Password)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
			utils.Logger.Error(err.Error())
		} else {
			ok, unitID, err := s.storage.IsUserInUnit(c, userID)
			if err != nil {
				utils.Logger.Error(err.Error())
				//TODO
			} else {
				if ok {
					ok, devices, err := s.storage.DoesUserHaveDevice(c, userID)
					if err != nil {
						utils.Logger.Error(err.Error())
						//TODO
					} else {
						if ok {
							for _, device := range devices {
								res, err := utils.MakeRequest(
									utils.NewRequest(
										ctx,
										s.connections[unitID],
										&proto.SpawnAndRunDevice{Device: device}),
								)
								if err != nil {
									//TODO
									utils.Logger.Error(err.Error())

								}
								ctx.Respond(res)
							}
						}
					}
				}
				//this unit need to have a device with owner of this id user
				// units needs to spawn a child ->get this PID and respond
			}
		}
	case *proto.GetPins:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pins, err := s.storage.GetPins(c)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.Pins{Pins: pins})
		}
	case *proto.GetCurrentTask:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		currentTask, err := s.storage.GetCurrentTask(c, msg.DeviceID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(currentTask)
		}
	case *proto.GetDeviceTypes:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		types, err := s.storage.GetDeviceTypes(c)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.DeviceTypes{Types: types})
		}
	case *proto.CreateDevice:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.storage.InsertDevice(c, msg.Device)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptCreateDevice{})
		}
	case *proto.GetUserInformation:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		information, err := s.storage.GetUserInformation(c, msg.UserID)
		if err != nil {
			fmt.Println(err)
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UserInformations{UserInformation: information})
		}
	case *proto.GetUnitInformation:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		usrInformation, err := s.storage.GetUnitInformation(c, msg.UnitID)
		if err != nil {
			fmt.Println(err)
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(usrInformation)
		}
	default:
		utils.Logger.Warn("server got unknown message", reflect.TypeOf(msg).String())

	}
}

func (s *Server) startHeartbeat(ctx *actor.Context) {
	for ID, PID := range s.connections {
		go func(pid *actor.PID, id string) {
			resp := ctx.Request(pid, &proto.Ping{}, utils.WaitTime)
			res, err := resp.Result()
			_, ok := res.(*proto.Pong)
			if !ok || err != nil {
				utils.Logger.Error("User is not responding for some time:", id, ctx.PID().String())
				ctx.Send(ctx.PID(), &proto.Disconnect{Id: id})
			}
		}(PID, ID)
	}
}
