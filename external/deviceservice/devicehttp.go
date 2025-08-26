package deviceservice

import (
	"github.com/janicaleksander/bcs/external/deviceservice/api"
)

type DeviceHTTP struct {
	handler *api.Handler
	//devicePID  *actor.PID
	//parent PID in the future
}

func NewHTTPDevice(handler *api.Handler) *DeviceHTTP {
	return &DeviceHTTP{
		handler: handler,
		//devicePID:  pid,
	}
}

// TODO idk if panic or return err
func (d *DeviceHTTP) RunHTTPServer() {
	d.handler.SetupRouter()
	d.handler.RunHTTP()
}
