package External

import "github.com/janicaleksander/bcs/Device"

type External struct {
	devices []*Device.Device
}

func NewExternal() *External {
	return &External{
		make([]*Device.Device, 1024),
	}
}
