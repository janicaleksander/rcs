package loginstate

import (
	"time"

	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (l *LoginScene) Reset() {
	l.isLoginButtonPressed = false
	l.errorSection.loginErrorMessage = ""

}

func (l *LoginScene) Login() {
	email := l.emailInput.GetText()
	pwd := l.passwordInput.GetText()
	if len(email) <= 0 || len(pwd) <= 0 {
		l.errorSection.loginErrorMessage = "Zero length inboxInput"
		l.errorSection.errorPopup.Show()
		l.scheduler.After((3 * time.Second).Seconds(), func() {
			l.errorSection.errorPopup.Hide()
		})
		return
	}

	res, err := utils.MakeRequest(utils.NewRequest(l.cfg.Ctx, l.cfg.ServerPID, &proto.LoginUser{
		Pid: &proto.PID{
			Address: l.cfg.Ctx.PID().GetAddress(),
			Id:      l.cfg.Ctx.PID().GetID(),
		},
		Email:    email,
		Password: pwd,
	}))
	if err != nil {
		//error context deadline exceeded
		l.errorSection.loginErrorMessage = err.Error()
		l.errorSection.errorPopup.Show()
		l.scheduler.After((3 * time.Second).Seconds(), func() {
			l.errorSection.errorPopup.Hide()
		})
		return
	}

	if v, ok := res.(*proto.AcceptLogin); ok {
		//to show login is taking to much time add some circle or infobar animation
		res, err := utils.MakeRequest(utils.NewRequest(l.cfg.Ctx, l.cfg.MessageServicePID, &proto.RegisterClient{
			Id: v.Id,
			Pid: &proto.PID{
				Address: l.cfg.Ctx.PID().Address,
				Id:      l.cfg.Ctx.PID().ID,
			},
		}))

		if err != nil {
			//context deadline exceeded
			l.errorSection.loginErrorMessage = err.Error()
			l.errorSection.errorPopup.Show()
			l.scheduler.After((3 * time.Second).Seconds(), func() {
				l.errorSection.errorPopup.Hide()
			})
			return
		}

		if _, ok := res.(*proto.SuccessRegisterClient); !ok {
			l.errorSection.loginErrorMessage = "CANT LOGIN TO Service message"

			//TODO but i dont know if this should be critical error to have error below
			l.errorSection.errorPopup.Show()
			l.scheduler.After((3 * time.Second).Seconds(), func() {
				l.errorSection.errorPopup.Hide()
			})
			return
		}
		//TODO if role is 5 this else if ... others

		//TO this point we have to determine if we have error in other services
		// and w.---.messageServiceError = true
		if v.RuleLevel == 5 {
			l.stateManager.Add(statesmanager.HCMenuState)
		}

	} else {
		//error
		l.errorSection.loginErrorMessage = "Invalid credentials"
		l.errorSection.errorPopup.Show()
		l.scheduler.After((3 * time.Second).Seconds(), func() {
			l.errorSection.errorPopup.Hide()
		})
	}
}
