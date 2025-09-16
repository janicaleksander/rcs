package unit

import (
	"context"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/external/device"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO device is going to be created and assign by unit commander/headcommander
// and then we can login by this device
// user_email -> isInUnit -> do have a device in this unit -> login in

// if we delete user, we have to think what to do with assign device(maybe move somewhere?)
type Unit struct {
	id        string
	serverPID *actor.PID
	devices   map[string]*actor.PID //device id to his PID
	storage   database.Storage
}

func NewUnit(id string, serverPID *actor.PID, storage database.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			id:        id,
			serverPID: serverPID,
			devices:   make(map[string]*actor.PID, 64),
			storage:   storage,
		}
	}
}

// actor with remote
func (u *Unit) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Unit has initialized")
	case actor.Started:
		utils.Logger.Info("Unit has started")
		ctx.Send(u.serverPID, &proto.LoginUnit{
			Pid: &proto.PID{
				Address: ctx.PID().Address,
				Id:      ctx.PID().ID,
			},
			UnitID: u.id,
		})
	case actor.Stopped:
		utils.Logger.Info("Unit has stopped")
	case *proto.Ping:
		ctx.Respond(&proto.Pong{})
	case *proto.SpawnAndRunDevice:
		pid := ctx.SpawnChild(device.NewDevice(msg.Device.Id, ctx.PID()), "device", actor.WithID(msg.Device.Id))
		u.devices[msg.Device.Id] = pid
		ctx.Respond(&proto.AcceptSpawnAndRunDevice{
			UserID:   msg.Device.Owner,
			DeviceID: msg.Device.Id,
			DevicePID: &proto.PID{
				Address: pid.Address,
				Id:      pid.ID,
			},
		})

	case *proto.UpdateLocationReq:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := u.storage.UpdateLocation(c, msg)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptUpdateLocationReq{})
		}
	case *proto.UserTaskReq:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t, err := u.storage.GetTask(c, msg.TaskID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UserTaskRes{Task: t})
		}
	case *proto.UserTasksReq:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tasks, err := u.storage.GetUserTasks(c, msg.DeviceID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UserTasksRes{Tasks: tasks})
		}

	default:
		utils.Logger.Warn("Unit got unknown message", "Type", reflect.TypeOf(msg).String())
		_ = msg
	}

}
