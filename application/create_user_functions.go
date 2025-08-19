package application

import (
	"strconv"
	"time"

	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/user"
	"github.com/janicaleksander/bcs/utils"
)

func (s *CreateUserScene) Reset() {
	s.errorSection.errorMessage = ""
	s.infoSection.acceptMessage = ""
	s.newUserSection.emailInput.Clear()
	s.newUserSection.passwordInput.Clear()
	s.newUserSection.rePasswordInput.Clear()
	s.newUserSection.ruleLevelInput.Clear()
	s.newUserSection.nameInput.Clear()
	s.newUserSection.surnameInput.Clear()
}

func (w *Window) CreateUser() {
	email := w.createUserScene.newUserSection.emailInput.GetText()
	password := w.createUserScene.newUserSection.passwordInput.GetText()
	rePassword := w.createUserScene.newUserSection.rePasswordInput.GetText()
	ruleLevel := w.createUserScene.newUserSection.ruleLevelInput.GetText()
	name := w.createUserScene.newUserSection.nameInput.GetText()
	surname := w.createUserScene.newUserSection.surnameInput.GetText()

	//check inboxInput
	if len(email) <= 0 || len(password) <= 0 ||
		len(rePassword) <= 0 || len(ruleLevel) <= 0 ||
		len(name) <= 0 || len(surname) <= 0 {
		w.createUserScene.errorSection.errorMessage = "Zero length error"
		w.createUserScene.scheduler.After((3 * time.Second).Seconds(), func() {
			w.createUserScene.errorSection.errorPopup.Hide()
		})
		w.createUserScene.errorSection.errorPopup.Show()
		return
	}
	lvl, err := strconv.Atoi(ruleLevel)
	// TODO curr max lvl
	if lvl > 5 || err != nil {
		w.createUserScene.errorSection.errorMessage = "Bad ruleLVL inboxInput"
		w.createUserScene.scheduler.After((3 * time.Second).Seconds(), func() {
			w.createUserScene.errorSection.errorPopup.Hide()
		})
		w.createUserScene.errorSection.errorPopup.Show()
		return
	}
	newUser := user.NewUser(email, password, int32(lvl), name, surname)

	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.CreateUser{User: newUser}))

	if err != nil {
		w.createUserScene.errorSection.errorMessage = "Actor ctx error"
		w.createUserScene.scheduler.After((3 * time.Second).Seconds(), func() {
			w.createUserScene.errorSection.errorPopup.Hide()
		})
		w.createUserScene.errorSection.errorPopup.Show()

	}
	if _, ok := res.(*proto.AcceptCreateUser); ok {
		w.createUserScene.infoSection.acceptMessage = "Created successfully"
		w.createUserScene.scheduler.After((3 * time.Second).Seconds(), func() {
			w.createUserScene.infoSection.infoPopup.Hide()
		})
		w.createUserScene.infoSection.infoPopup.Show()
		w.createUserScene.newUserSection.emailInput.Clear()
		w.createUserScene.newUserSection.passwordInput.Clear()
		w.createUserScene.newUserSection.rePasswordInput.Clear()
		w.createUserScene.newUserSection.ruleLevelInput.Clear()
		w.createUserScene.newUserSection.nameInput.Clear()
		w.createUserScene.newUserSection.surnameInput.Clear()
	} else {
		w.createUserScene.errorSection.errorMessage = "DB deny!"
		w.createUserScene.scheduler.After((3 * time.Second).Seconds(), func() {
			w.createUserScene.errorSection.errorPopup.Hide()
		})
		w.createUserScene.errorSection.errorPopup.Show()

	}

}
