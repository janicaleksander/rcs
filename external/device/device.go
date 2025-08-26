package device

import "github.com/anthdm/hollywood/actor"

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
	default:
		_ = msg
	}
}
