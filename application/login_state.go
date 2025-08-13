package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

// LOGIN STATE
type LoginScene struct {
	loginButton   Button
	emailInput    InputField
	passwordInput InputField

	isLoginError      bool
	loginErrorMessage string
	errorPosition     Position
}

func (l *LoginScene) Reset() {
	l.isLoginError = false
	l.loginErrorMessage = ""
}
func (w *Window) loginSceneSetup() {
	w.loginSceneScene.Reset()
	w.loginSceneScene.loginButton = Button{
		bounds: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2),
			200, 40,
		),
		text: "Login",
	}
	w.loginSceneScene.emailInput.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	)
	w.loginSceneScene.passwordInput.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	)
	w.loginSceneScene.errorPosition = Position{
		x: int32(rl.GetScreenWidth()/2 - 100),
		y: int32(rl.GetScreenHeight()/2 + 40),
	}
}

//maybe move this to setup?

func (w *Window) updateLoginState() {
	if gui.Button(w.loginSceneScene.loginButton.bounds, w.loginSceneScene.loginButton.text) {
		//i have to do check services and then mark this somehow to show that user can use it

		//and maybe use this to not make other request we have to wait if goruitne change this var to false and then???
		//thiis is to set own presence to cut all messsage service from app

		email := w.loginSceneScene.emailInput.text
		pwd := w.loginSceneScene.passwordInput.text
		if len(email) <= 0 || len(pwd) <= 0 {
			w.loginSceneScene.isLoginError = true
			w.loginSceneScene.loginErrorMessage = "Zero length inboxInput"
			return
		}
		resp := w.ctx.Request(w.serverPID, &proto.LoginUser{
			Pid: &proto.PID{
				Address: w.ctx.PID().GetAddress(),
				Id:      w.ctx.PID().GetID(),
			},
			Email:    email,
			Password: pwd,
		}, utils.WaitTime)
		val, err := resp.Result()
		if err != nil {
			w.loginSceneScene.isLoginError = true
			w.loginSceneScene.loginErrorMessage = err.Error()
		} else if v, ok := val.(*proto.AcceptLogin); ok {
			//to show login is taking to much time add some circle or infobar animation
			resp := w.ctx.Request(w.messageServicePID, &proto.RegisterClient{
				Id: v.Id,
				Pid: &proto.PID{
					Address: w.ctx.PID().Address,
					Id:      w.ctx.PID().ID,
				},
			}, utils.WaitTime)
			res, err := resp.Result()
			if _, ok := res.(*proto.SuccessRegisterClient); !ok || err != nil {
				//mark this sth
				panic(err)
			}
			//TODO if role is 5 this else if ... others

			//TO this point we have to determine if we have error in other services
			// and w.---.messageServiceError = true
			if v.RuleLevel == 5 {
				w.menuHCSceneSetup()
				w.currentState = HCMenuState
				w.sceneStack = append(w.sceneStack, HCMenuState)
			}

		} else if msg, ok := val.(*proto.DenyLogin); ok {
			w.loginSceneScene.isLoginError = true
			w.loginSceneScene.loginErrorMessage = msg.Info
		} else {
			w.loginSceneScene.isLoginError = true
			w.loginSceneScene.loginErrorMessage = "Unknown response type"
		}

	}

}

// TOOD during the login is time to connect to all other services e.g messageService
// (we also are connecting to server but this is what we are doing first because without server
// we cant live but without messages services we can

// TODO (user can be offline but also error with messageSrvies is reason to set his status to offline
// even he is logged in
func (w *Window) renderLoginState() {
	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)
	if w.loginSceneScene.isLoginError {
		rl.DrawText(w.loginSceneScene.loginErrorMessage,
			w.loginSceneScene.errorPosition.x,
			w.loginSceneScene.errorPosition.y,
			20,
			rl.Red)
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.loginSceneScene.emailInput.bounds) {
			w.loginSceneScene.emailInput.focus = true
			w.loginSceneScene.passwordInput.focus = false
		} else if rl.CheckCollisionPointRec(mousePos, w.loginSceneScene.passwordInput.bounds) {
			w.loginSceneScene.emailInput.focus = false
			w.loginSceneScene.passwordInput.focus = true
		} else {
			w.loginSceneScene.emailInput.focus = false
			w.loginSceneScene.passwordInput.focus = false
		}
	}
	gui.TextBox(w.loginSceneScene.emailInput.bounds, &w.loginSceneScene.emailInput.text, 64, w.loginSceneScene.emailInput.focus)

	//TODO: secret password inboxInput box
	gui.TextBox(w.loginSceneScene.passwordInput.bounds, &w.loginSceneScene.passwordInput.text, 64, w.loginSceneScene.passwordInput.focus)

}
