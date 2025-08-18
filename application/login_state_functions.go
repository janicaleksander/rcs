package application

import (
	"github.com/janicaleksander/bcs/utils"

	"github.com/janicaleksander/bcs/proto"
)

func (l *LoginScene) Reset() {
	l.isLoginButtonPressed = false
	l.isLoginError = false
	l.loginErrorMessage = ""

}

func (w *Window) Login() {
	email := w.loginScene.emailInput.GetText()
	pwd := w.loginScene.passwordInput.GetText()
	if len(email) <= 0 || len(pwd) <= 0 {
		w.loginScene.isLoginError = true
		w.loginScene.loginErrorMessage = "Zero length inboxInput"
		return
	}

	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.LoginUser{
		Pid: &proto.PID{
			Address: w.ctx.PID().GetAddress(),
			Id:      w.ctx.PID().GetID(),
		},
		Email:    email,
		Password: pwd,
	}))
	if err != nil {
		//error context deadline exceeded
		w.loginScene.isLoginError = true
		w.loginScene.loginErrorMessage = err.Error()
		return
	}

	if v, ok := res.(*proto.AcceptLogin); ok {
		//to show login is taking to much time add some circle or infobar animation

		res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.RegisterClient{
			Id: v.Id,
			Pid: &proto.PID{
				Address: w.ctx.PID().Address,
				Id:      w.ctx.PID().ID,
			},
		}))

		if err != nil {
			//context deadline exceeded
		}

		if _, ok := res.(*proto.SuccessRegisterClient); !ok || err != nil {
			//mark this sth
		}
		//TODO if role is 5 this else if ... others

		//TO this point we have to determine if we have error in other services
		// and w.---.messageServiceError = true
		if v.RuleLevel == 5 {
			w.menuHCSceneSetup()
			w.currentState = HCMenuState
			w.sceneStack = append(w.sceneStack, HCMenuState)
		}

	} else {
		//error
		/*
			w.loginScene.isLoginError = true
			w.loginScene.loginErrorMessage = msg.Info
			}
		*/

	}
}
