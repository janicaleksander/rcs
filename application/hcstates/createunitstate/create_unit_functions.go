package createunitstate

import (
	"time"

	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (c *CreateUnitScene) Reset() {
	c.errorSection.isSetupError = false
	c.errorSection.isCreateError = false
	c.infoSection.isInfoMessage = false
	c.infoSection.infoMessage = ""
	c.errorSection.errorMessage = ""
}
func (c *CreateUnitScene) FetchUsers() {
	//TODO get proper lvl value
	res, err := utils.MakeRequest(utils.NewRequest(c.cfg.Ctx, c.cfg.ServerPID, &proto.GetUserAboveLVL{Lvl: -1}))
	if err != nil {
		//context deadline exceeded
		//do sth with that
		c.errorSection.isSetupError = true
	}

	if v, ok := res.(*proto.UsersAboveLVL); ok {
		c.newUnitSection.usersDropdown.Strings = make([]string, 0, 64)
		c.newUnitSection.usersDropdown.Strings = append(c.newUnitSection.usersDropdown.Strings,
			"Choose user by his ID")
		for _, user := range v.Users {
			c.newUnitSection.usersDropdown.Strings = append(c.newUnitSection.usersDropdown.Strings,
				user.Id+"\n"+user.Email)
		}
	} else {
		c.errorSection.isSetupError = true
	}

}
func (c *CreateUnitScene) CreateUnit() {
	name := c.newUnitSection.nameInput.GetText()
	user := c.newUnitSection.usersDropdown.IdxActiveElement
	if len(name) <= 0 || user <= 0 {
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		c.errorSection.errorPopup.Show()
		c.errorSection.errorMessage = "Zero length error"
	} else {
		//user can be only in one unit in the same time -> error
		res, err := utils.MakeRequest(utils.NewRequest(c.cfg.Ctx, c.cfg.ServerPID, &proto.CreateUnit{
			Name:   name,
			UserID: c.newUnitSection.usersDropdown.Strings[user],
		}))

		if err != nil {
			//context deadline exceeded
			c.scheduler.After((3 * time.Second).Seconds(), func() {
				c.errorSection.errorPopup.Hide()
			})
			c.errorSection.errorPopup.Show()
			c.errorSection.errorMessage = err.Error()
			return
		}
		if _, ok := res.(*proto.AcceptCreateUnit); ok {
			c.scheduler.After((3 * time.Second).Seconds(), func() {
				c.infoSection.infoPopup.Hide()
			})
			c.infoSection.infoPopup.Show()
			c.infoSection.infoMessage = "Success!"
		} else {
			c.scheduler.After((3 * time.Second).Seconds(), func() {
				c.errorSection.errorPopup.Hide()
			})
			c.errorSection.errorPopup.Show()
			c.errorSection.errorMessage = "Error!"

		}

	}

}
