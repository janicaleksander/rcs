package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
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
	w.loginSceneData.Reset()
	w.loginSceneData.loginButton = Button{
		bounds: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2),
			200, 40,
		),
		text: "Login",
	}
	w.loginSceneData.emailInput.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	)
	w.loginSceneData.passwordInput.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	)
	w.loginSceneData.errorPosition = Position{
		x: int32(rl.GetScreenWidth()/2 - 100),
		y: int32(rl.GetScreenHeight()/2 + 40),
	}
}
func (w *Window) updateLoginState() {
	if gui.Button(w.loginSceneData.loginButton.bounds, w.loginSceneData.loginButton.text) {
		email := w.loginSceneData.emailInput.text
		pwd := w.loginSceneData.passwordInput.text
		if len(email) <= 0 || len(pwd) <= 0 {
			w.loginSceneData.isLoginError = true
			w.loginSceneData.loginErrorMessage = "Zero length input"
			return
		}
		pid := &Proto.PID{
			Address: w.ctx.PID().GetAddress(),
			Id:      w.ctx.PID().GetID(),
		}
		resp := w.ctx.Request(w.serverPID, &Proto.LoginUser{
			Pid:      pid,
			Email:    email,
			Password: pwd,
		}, WaitTime)
		val, err := resp.Result()
		if err != nil {
			w.loginSceneData.isLoginError = true
			w.loginSceneData.loginErrorMessage = err.Error()
		} else if v, ok := val.(*Proto.AcceptLogin); ok {
			//TODO
			//if role is 5 this
			//else if ... others
			if v.RuleLevel == 5 {
				w.menuHCSceneSetup()
				w.currentState = HCMenuState
				w.sceneStack = append(w.sceneStack, HCMenuState)
			}

		} else if msg, ok := val.(*Proto.DenyLogin); ok {
			w.loginSceneData.isLoginError = true
			w.loginSceneData.loginErrorMessage = msg.Info
		} else {
			w.loginSceneData.isLoginError = true
			w.loginSceneData.loginErrorMessage = "Unknown response type"
		}

	}

}
func (w *Window) renderLoginState() {
	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)
	if w.loginSceneData.isLoginError {
		rl.DrawText(w.loginSceneData.loginErrorMessage,
			w.loginSceneData.errorPosition.x,
			w.loginSceneData.errorPosition.y,
			20,
			rl.Red)
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.loginSceneData.emailInput.bounds) {
			w.loginSceneData.emailInput.focus = true
			w.loginSceneData.passwordInput.focus = false
		} else if rl.CheckCollisionPointRec(mousePos, w.loginSceneData.passwordInput.bounds) {
			w.loginSceneData.emailInput.focus = false
			w.loginSceneData.passwordInput.focus = true
		} else {
			w.loginSceneData.emailInput.focus = false
			w.loginSceneData.passwordInput.focus = false
		}
	}
	gui.TextBox(w.loginSceneData.emailInput.bounds, &w.loginSceneData.emailInput.text, 64, w.loginSceneData.emailInput.focus)

	//TODO: secret password input box
	gui.TextBox(w.loginSceneData.passwordInput.bounds, &w.loginSceneData.passwordInput.text, 64, w.loginSceneData.passwordInput.focus)

}
