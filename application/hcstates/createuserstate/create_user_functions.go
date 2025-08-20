package createuserstate

import (
	"strconv"
	"time"

	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/user"
	"github.com/janicaleksander/bcs/utils"
)

func (c *CreateUserScene) Reset() {
	c.errorSection.errorMessage = ""
	c.infoSection.acceptMessage = ""
	c.newUserSection.emailInput.Clear()
	c.newUserSection.passwordInput.Clear()
	c.newUserSection.rePasswordInput.Clear()
	c.newUserSection.ruleLevelInput.Clear()
	c.newUserSection.nameInput.Clear()
	c.newUserSection.surnameInput.Clear()
}

func (c *CreateUserScene) CreateUser() {
	email := c.newUserSection.emailInput.GetText()
	password := c.newUserSection.passwordInput.GetText()
	rePassword := c.newUserSection.rePasswordInput.GetText()
	ruleLevel := c.newUserSection.ruleLevelInput.GetText()
	name := c.newUserSection.nameInput.GetText()
	surname := c.newUserSection.surnameInput.GetText()

	//check inboxInput
	if len(email) <= 0 || len(password) <= 0 ||
		len(rePassword) <= 0 || len(ruleLevel) <= 0 ||
		len(name) <= 0 || len(surname) <= 0 {
		c.errorSection.errorMessage = "Zero length error"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		c.errorSection.errorPopup.Show()
		return
	}
	lvl, err := strconv.Atoi(ruleLevel)
	// TODO curr max lvl
	if lvl > 5 || err != nil {
		c.errorSection.errorMessage = "Bad ruleLVL inboxInput"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		c.errorSection.errorPopup.Show()
		return
	}
	newUser := user.NewUser(email, password, int32(lvl), name, surname)

	res, err := utils.MakeRequest(utils.NewRequest(c.cfg.Ctx, c.cfg.ServerPID, &proto.CreateUser{User: newUser}))

	if err != nil {
		c.errorSection.errorMessage = "Actor ctx error"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		c.errorSection.errorPopup.Show()

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
		c.newUserSection.ruleLevelInput.Clear()
		c.newUserSection.nameInput.Clear()
		c.newUserSection.surnameInput.Clear()
	} else {
		c.errorSection.errorMessage = "DB deny!"
		c.scheduler.After((3 * time.Second).Seconds(), func() {
			c.errorSection.errorPopup.Hide()
		})
		c.errorSection.errorPopup.Show()

	}

}
