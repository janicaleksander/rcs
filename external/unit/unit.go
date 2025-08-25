package unit

import (
	"reflect"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/external/deviceservice"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO device is going to be created and assign by unit commander/headcommander
// and then we can login by this device
// user_email -> isInUnit -> do have a device in this unit -> login in

// if we delete user, we have to think what to do with assign device(maybe move somewhere?)
type Unit struct {
	id      string // uuid
	devices []*proto.Device
}

func NewUnit(id string) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			id:      id,
			devices: make([]*proto.Device, 0, 64),
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
		u.devices = append(u.devices, msg.Device)
		pid := ctx.SpawnChild(deviceservice.NewDeviceActor(), "device", actor.WithID(msg.Device.Id))
		ctx.Respond(&proto.SuccessSpawnDevice{
			UserID: msg.Device.Owner,
			DevicePID: &proto.PID{
				Address: pid.Address,
				Id:      pid.ID,
			},
		})
	default:
		utils.Logger.Warn("Unit got unknown message", "Type", reflect.TypeOf(msg).String())
		_ = msg
	}

}
