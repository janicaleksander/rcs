package createdevicestate

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (d *CreateDeviceScene) Reset() {

}

func (d *CreateDeviceScene) FetchUsers() {
	res, err := utils.MakeRequest(utils.NewRequest(d.cfg.Ctx, d.cfg.ServerPID, &proto.GetUserAboveLVL{
		Lower: 0,
		Upper: 3,
	}))
	if err != nil {
		d.errorSection.error = err.Error()
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}
	d.newDeviceSection.users = make([]*proto.User, 0, 64)
	d.newDeviceSection.ownerSlider.Strings = make([]string, 0, 64)
	if v, ok := res.(*proto.UsersAboveLVL); ok {
		for _, user := range v.Users {
			d.newDeviceSection.users = append(d.newDeviceSection.users, user)
			d.newDeviceSection.ownerSlider.Strings = append(d.newDeviceSection.ownerSlider.Strings,
				user.Personal.Name, user.Personal.Surname+"\n"+user.Email)
		}
	} else {
		v, _ := res.(*proto.Error)
		d.errorSection.error = v.Content
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

}

func (d *CreateDeviceScene) FetchDeviceTypes() {
	res, err := utils.MakeRequest(utils.NewRequest(d.cfg.Ctx, d.cfg.ServerPID, &proto.GetDeviceTypes{}))
	if err != nil {
		d.errorSection.error = err.Error()
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}
	d.newDeviceSection.deviceTypes = make([]int, 0, 64)
	if v, ok := res.(*proto.DeviceTypes); ok {
		for _, t := range v.Types {
			d.newDeviceSection.deviceTypes = append(d.newDeviceSection.deviceTypes, int(t))
		}
	} else {
		v, _ := res.(*proto.Error)
		d.errorSection.error = v.Content
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

}

func (d *CreateDeviceScene) CreateDevice() {
	name := d.newDeviceSection.nameInput.GetText()
	currOwnerIdx := d.newDeviceSection.ownerSlider.IdxActiveElement
	currTypeIdx := d.newDeviceSection.typesToggle.Selected
	if len(strings.TrimSpace(name)) == 0 || currOwnerIdx == -1 {
		d.errorSection.error = "Bad inputs"
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}
	id := uuid.New()
	owner := d.newDeviceSection.users[currOwnerIdx].Id
	t := d.newDeviceSection.deviceTypes[currTypeIdx]
	device := &proto.Device{
		Id:             id.String(),
		Name:           name,
		Owner:          owner,
		LastTimeOnline: nil,
		Type:           int32(t),
	}

	res, err := utils.MakeRequest(
		utils.NewRequest(d.cfg.Ctx, d.cfg.ServerPID,
			&proto.CreateDevice{Device: device}))

	if err != nil {
		d.errorSection.error = "Ctx error"
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

	if _, ok := res.(*proto.AcceptCreateDevice); ok {
		d.infoSection.info = "Good!"
		d.infoSection.infoPopup.ShowFor(time.Second * 3)
	} else {
		d.errorSection.error = "Answer error"
		d.errorSection.errorPopup.ShowFor(time.Second * 3)
	}
}
