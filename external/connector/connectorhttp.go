package connector

import (
	"github.com/janicaleksander/bcs/external/connector/api"
)

// TODO think about what if user have more than device
// what do i need do then in flutter login page??

// or block ability to have more than one with the same type device
// and maybe we will check it internally -> mobile app we can login if type matches
// smartwatch app - we cant if we have only type android etc
type DeviceHTTP struct {
	handler *api.Handler
	//devicePID  *actor.PID
	//parent PID in the future
}

func NewServiceHTTPDevice(handler *api.Handler) *DeviceHTTP {
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
