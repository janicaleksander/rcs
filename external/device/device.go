package device

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

type Device struct {
	unitPID *actor.PID
	id      string // uuid
}

func NewDevice(id string, unitPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Device{
			id:      id,
			unitPID: unitPID,
		}
	}
}

func (d *Device) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Device", d.id, "is initialized!")
	case actor.Started:
		utils.Logger.Info("Device", d.id, " started")
	case actor.Stopped:
		utils.Logger.Info("Device", d.id, " stopped")
	case *proto.UpdateLocationReq:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, d.unitPID, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.UserTaskReq:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, d.unitPID, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.UserTasksReq:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, d.unitPID, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.UpdateCurrentTaskReq:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, d.unitPID, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.DeleteTaskReq:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, d.unitPID, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}

	default:
		_ = msg
	}
}
