package unit

import (
	"context"
	"reflect"

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
	id      string                // uuid
	devices map[string]*actor.PID //device id to his PID
	storage database.Storage
}

func NewUnit(id string, storage database.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			id:      id,
			devices: make(map[string]*actor.PID, 64),
			storage: storage,
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
	case actor.Stopped:
		utils.Logger.Info("Unit has stopped")
	case *proto.Ping:
		ctx.Respond(&proto.Pong{})
	case *proto.SpawnAndRunDevice:
		pid := ctx.SpawnChild(device.NewDevice(msg.Device.Id, ctx.PID()), "device", actor.WithID(msg.Device.Id))
		u.devices[msg.Device.Id] = pid
		ctx.Respond(&proto.SuccessSpawnDevice{
			UserID:   msg.Device.Owner,
			DeviceID: msg.Device.Id,
			DevicePID: &proto.PID{
				Address: pid.Address,
				Id:      pid.ID,
			},
		})

	case *proto.UpdateLocationReq:
		c := context.Background()
		err := u.storage.UpdateLocation(c, msg)
		if err != nil {
			ctx.Respond(&proto.FailureUpdateLocationReq{})
		} else {
			ctx.Respond(&proto.SuccessUpdateLocationReq{})
		}

	default:
		utils.Logger.Warn("Unit got unknown message", "Type", reflect.TypeOf(msg).String())
		_ = msg
	}

}
