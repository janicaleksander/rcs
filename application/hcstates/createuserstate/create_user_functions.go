package createuserstate

import (
	"time"

	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/types/user"
	"github.com/janicaleksander/bcs/utils"
)

func (c *CreateUserScene) Reset() {
	c.errorSection.errorMessage = ""
	c.infoSection.acceptMessage = ""
	c.newUserSection.emailInput.Clear()
	c.newUserSection.passwordInput.Clear()
	c.newUserSection.rePasswordInput.Clear()
	c.newUserSection.nameInput.Clear()
	c.newUserSection.surnameInput.Clear()
}

func (c *CreateUserScene) CreateUser() {
	email := c.newUserSection.emailInput.GetText()
	password := c.newUserSection.passwordInput.GetText()
	rePassword := c.newUserSection.rePasswordInput.GetText()
	name := c.newUserSection.nameInput.GetText()
	surname := c.newUserSection.surnameInput.GetText()

	//check inboxInput
	if len(email) <= 0 || len(password) <= 0 ||
		len(rePassword) <= 0 ||
		len(name) <= 0 || len(surname) <= 0 {
		c.errorSection.errorMessage = "Zero length input error"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		c.errorSection.errorPopup.Show()
		return
	}
	lvl := c.newUserSection.ruleLevelToggleGroup.Selected
	if lvl > 3 || lvl < 0 {
		c.errorSection.errorPopup.Show()
		c.errorSection.errorMessage = "Bad rulelvl input"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		return
	}
	newUser := user.NewUser(email, password, int32(lvl), name, surname)
	res, err := utils.MakeRequest(
		utils.NewRequest(
			c.cfg.Ctx,
			c.cfg.ServerPID,
			&proto.CreateUser{
				User: newUser,
			},
		))

	if err != nil {
		c.errorSection.errorMessage = err.Error()
		c.errorSection.errorPopup.Show()
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
	}
	if _, ok := res.(*proto.AcceptCreateUser); ok {
		c.infoSection.acceptMessage = "Created successfully"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.infoSection.infoPopup.Hide()
		})
		c.infoSection.infoPopup.Show()
		c.newUserSection.emailInput.Clear()
		c.newUserSection.passwordInput.Clear()
		c.newUserSection.rePasswordInput.Clear()
		c.newUserSection.nameInput.Clear()
		c.newUserSection.surnameInput.Clear()
	} else {
		v, _ := res.(*proto.Error)
		c.errorSection.errorPopup.Show()
		c.errorSection.errorMessage = v.Content
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
	}

}
