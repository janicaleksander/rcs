package device

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/utils"
)

type Device struct {
	id string // uuid
}

func NewDevice(id string) actor.Producer {
	return func() actor.Receiver {
		return &Device{
			id: id,
		}
	}
}

func (d *Device) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Device",d.id,"is initialized!")
	case actor.Started:
		utils.Logger.Info("Device",d.id," started")
	case actor.Stopped:
		utils.Logger.Info("Device",d.id," stopped")
	default:
		_ = msg
	}
}
