package external

import "github.com/janicaleksander/bcs/device"

type External struct {
	devices []*device.Device
}

func NewExternal() *External {
	return &External{
		make([]*device.Device, 1024),
	}
}
