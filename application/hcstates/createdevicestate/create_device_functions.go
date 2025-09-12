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
	res, err := utils.MakeRequest(utils.NewRequest(d.cfg.Ctx, d.cfg.ServerPID, &proto.GetUserAboveLVL{Lvl: -1}))
	if err != nil {
		// context deadline exceeded
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
		// TODO error
	}

}

func (d *CreateDeviceScene) FetchDeviceTypes() {
	res, err := utils.MakeRequest(utils.NewRequest(d.cfg.Ctx, d.cfg.ServerPID, &proto.GetDeviceTypes{}))
	if err != nil {
		//context deadline error
		return
	}
	d.newDeviceSection.deviceTypes = make([]string, 0, 64)
	d.newDeviceSection.typeSlider.Strings = make([]string, 0, 64)
	if v, ok := res.(*proto.DeviceTypes); ok {
		for _, t := range v.Types {
			d.newDeviceSection.deviceTypes = append(d.newDeviceSection.deviceTypes, t)
			d.newDeviceSection.typeSlider.Strings = append(d.newDeviceSection.typeSlider.Strings, t)
		}
	} else {
		//error
	}

}

func (d *CreateDeviceScene) CreateDevice() {
	//add to device
	name := d.newDeviceSection.nameInput.GetText()
	currOwnerIdx := d.newDeviceSection.ownerSlider.IdxActiveElement
	currTypeIdx := d.newDeviceSection.typeSlider.IdxActiveElement
	if len(strings.TrimSpace(name)) == 0 || currTypeIdx == -1 || currOwnerIdx == -1 {
		d.errorSection.error = "Bad inputs"
		d.errorSection.errorPopup.Show()
		d.scheduler.After((3 * time.Second).Seconds(), func() {
			d.errorSection.errorPopup.Hide()
		})
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
		Type:           t,
	}

	res, err := utils.MakeRequest(
		utils.NewRequest(d.cfg.Ctx, d.cfg.ServerPID,
			&proto.CreateDevice{Device: device}))

	if err != nil {
		d.errorSection.error = "Ctx error"
		d.errorSection.errorPopup.Show()
		d.scheduler.After((3 * time.Second).Seconds(), func() {
			d.errorSection.errorPopup.Hide()
		})
	}

	if _, ok := res.(*proto.AcceptCreateDevice); ok {
		d.infoSection.info = "Good!"
		d.infoSection.infoPopup.Show()
		d.scheduler.After((3 * time.Second).Seconds(), func() {
			d.infoSection.infoPopup.Hide()
		})
	} else {
		d.errorSection.error = "Answer error"
		d.errorSection.errorPopup.Show()
		d.scheduler.After((3 * time.Second).Seconds(), func() {
			d.errorSection.errorPopup.Hide()
		})
	}
}
