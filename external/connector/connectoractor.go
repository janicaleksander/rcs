package connector

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO 025/08/26 12:23:28 ERROR Actor name already claimed pid=server/primary/unit/f3dcdf75-7555-40b1-8e00-5b5be0a2e039/device/2aef0730-1ccd-487d-add3-27ddb1660a41
type DeviceActor struct {
	ctx         *actor.Context
	connections map[string]*actor.PID // device  who is using particular device to device PID
}

func NewServiceDeviceActor() actor.Producer {
	return func() actor.Receiver {
		return &DeviceActor{
			connections: make(map[string]*actor.PID, 64),
		}
	}
}

func (d *DeviceActor) Receive(ctx *actor.Context) {
	d.ctx = ctx
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Device actor initialized")
	case actor.Started:
		utils.Logger.Info("Device actor started")
	case actor.Stopped:
		utils.Logger.Info("Device actor stopped")
	case actor.Context:
		ctx.Respond(ctx)
	case *proto.ConnectHDeviceToADevice:
		d.connections[msg.DeviceID] = actor.NewPID(msg.DevicePID.Address, msg.DevicePID.Id)
	case *proto.UpdateLocationReq:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, d.connections[msg.DeviceID], msg))
		if err != nil {
			d.ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			d.ctx.Respond(res)
		}
	default:
		utils.Logger.Info("Unrecognized message!")
		_ = msg
	}

}

//FLOW: device is pinging server during login to spawn new child od Device actor and
//return unit PID and this child PID
