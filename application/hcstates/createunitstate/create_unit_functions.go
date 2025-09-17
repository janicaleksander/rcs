package createunitstate

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (c *CreateUnitScene) Reset() {
	c.errorSection.isSetupError = false
	c.errorSection.isCreateError = false
	c.infoSection.isInfoMessage = false
	c.infoSection.infoMessage = ""
	c.errorSection.errorMessage = ""
	c.newUnitSection.acceptButton.Active()
}
func (c *CreateUnitScene) FetchUsers() {
	res, err := utils.MakeRequest(utils.NewRequest(c.cfg.Ctx, c.cfg.ServerPID, &proto.GetUserAboveLVL{
		Lower: 0,
		Upper: 1,
	}))
	if err != nil {
		c.errorSection.isSetupError = true
		return
	}
	if v, ok := res.(*proto.UsersAboveLVL); ok {
		c.newUnitSection.usersDropdown.Strings = make([]string, 0, 16)
		c.newUnitSection.usersDropdown.Strings = append(c.newUnitSection.usersDropdown.Strings,
			"Choose user by his ID")
		for _, user := range v.Users {
			c.newUnitSection.usersDropdown.Strings = append(c.newUnitSection.usersDropdown.Strings,
				user.Id+"\n"+user.Email)
		}
	} else {
		v, _ := res.(*proto.Error)
		c.errorSection.errorMessage = v.Content
		c.errorSection.isSetupError = true
	}

}
func (c *CreateUnitScene) CreateUnit() {
	name := c.newUnitSection.nameInput.GetText()
	user := c.newUnitSection.usersDropdown.IdxActiveElement
	if len(strings.TrimSpace(name)) <= 0 || user <= 0 {
		c.errorSection.errorPopup.Show()
		c.errorSection.errorMessage = "Zero length error"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
	} else {
		//user can be only in one unit in the same time -> error
		res, err := utils.MakeRequest(utils.NewRequest(c.cfg.Ctx, c.cfg.ServerPID, &proto.CreateUnit{
			Unit: &proto.Unit{
				Id:   uuid.New().String(),
				Name: name,
			},
			UserID: c.newUnitSection.usersDropdown.Strings[user][:36],
		}))

		if err != nil {
			c.errorSection.errorPopup.Show()
			c.errorSection.errorMessage = err.Error()
			c.scheduler.After((3 * time.Second).Seconds(), func() {
				c.errorSection.errorPopup.Hide()
			})
			return
		}
		if _, ok := res.(*proto.AcceptCreateUnit); ok {
			c.infoSection.infoPopup.Show()
			c.infoSection.infoMessage = "Success!"
			c.scheduler.After((3 * time.Second).Seconds(), func() {
				c.infoSection.infoPopup.Hide()
			})
		} else {
			v, _ := res.(*proto.Error)
			c.errorSection.errorMessage = v.Content
			c.errorSection.errorPopup.Show()
			c.scheduler.After((3 * time.Second).Seconds(), func() {
				c.errorSection.errorPopup.Hide()
			})

		}

	}

}
