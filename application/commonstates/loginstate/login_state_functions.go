package loginstate

import (
	"strings"
	"time"

	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (l *LoginScene) Reset() {
	l.loginSection.isLoginButtonPressed = false
	l.errorSection.loginErrorMessage = ""
	l.loginSection.loginButton.Active()
}

//user has to connect to the server -> critical error but this is on main file lvl

// user has to connect to message service -> not a critical error ->
// -> we have to disconnect from inbox section and send message in user info
func (l *LoginScene) Login() {
	email := l.loginSection.emailInput.GetText()
	pwd := l.loginSection.passwordInput.GetText()
	if len(strings.TrimSpace(email)) <= 0 || len(strings.TrimSpace(pwd)) <= 0 {
		l.errorSection.loginErrorMessage = "zero length inbox input"
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
		l.errorSection.loginErrorMessage = err.Error()
		l.errorSection.errorPopup.Show()
		l.scheduler.After((3 * time.Second).Seconds(), func() {
			l.errorSection.errorPopup.Hide()
		})
		return
	}

	if v, ok := res.(*proto.AcceptUserLogin); !ok {
		l.errorSection.loginErrorMessage = "Invalid credentials"
		l.errorSection.errorPopup.Show()
		l.scheduler.After((3 * time.Second).Seconds(), func() {
			l.errorSection.errorPopup.Hide()
		})
	} else {
		res, err = utils.MakeRequest(utils.NewRequest(
			l.cfg.Ctx,
			l.cfg.MessageServicePID, &proto.RegisterClientInMessageService{
				Id: v.UserID,
				Pid: &proto.PID{
					Address: l.cfg.Ctx.PID().Address,
					Id:      l.cfg.Ctx.PID().ID,
				},
			}))

		if err != nil {
			l.errorSection.loginErrorMessage = err.Error()
			l.errorSection.errorPopup.Show()
			l.scheduler.After((3 * time.Second).Seconds(), func() {
				l.errorSection.errorPopup.Hide()
			})
			return
		}

		if _, ok = res.(*proto.AcceptRegisterClient); !ok {
			l.errorSection.loginErrorMessage = "can't login to message service"
			l.errorSection.errorPopup.Show()
			l.scheduler.After((3 * time.Second).Seconds(), func() {
				l.errorSection.errorPopup.Hide()
			})
			return
		}
		//TODO
		if v.RuleLevel == 3 {
			l.stateManager.Add(statesmanager.HCMenuState)
		}

	}
}
