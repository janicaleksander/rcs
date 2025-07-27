package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"time"
)

// LOGIN STATE
type LoginScene struct {
	loginButton    Button
	emailBounds    rl.Rectangle
	passwordBounds rl.Rectangle
	emailTXT       string
	emailFocus     bool

	passwordTXT    string
	passwordTmpTXT string
	passwordFocus  bool

	isLoginError      bool
	loginErrorMessage string
}

func (w *Window) loginSceneSetup() {
	w.loginSceneData.loginButton = Button{
		position: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2),
			200, 40,
		),
		text: "Login",
	}
	w.loginSceneData.emailBounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	)
	w.loginSceneData.passwordBounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	)
}
func (w *Window) updateLoginState() {
	//this is already rendering this button on the screen
	//block multiple pressed button before receive respond, for const time LoginTime

	if gui.Button(w.loginSceneData.loginButton.position, w.loginSceneData.loginButton.text) {
		//valid login credentials -> send to server loginUser, wait for response
		email := w.loginSceneData.emailTXT
		pwd := w.loginSceneData.passwordTXT
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
		}, time.Second*WaitToLogin)
		val, err := resp.Result()
		//only if error this true
		w.loginSceneData.isLoginError = true
		if err != nil {
			w.loginSceneData.loginErrorMessage = err.Error()
		} else if _, ok := val.(*Proto.AcceptLogin); ok {
			v, _ := val.(*Proto.AcceptLogin)
			//TODO
			//if role is 5 this
			//else if ... others
			if v.RuleLevel == 5 {
				w.loginSceneData.isLoginError = false
				w.menuHCSceneSetup()
				w.currentState = HCMenuState
				w.sceneStack = append(w.sceneStack, HCMenuState)
			}

		} else if msg, ok := val.(*Proto.DenyLogin); ok {
			w.loginSceneData.loginErrorMessage = msg.Info
		} else {
			w.loginSceneData.loginErrorMessage = "Unknown response type"
		}

	}

}
func (w *Window) renderLoginState() {
	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)

	if w.loginSceneData.isLoginError {
		rl.DrawText(w.loginSceneData.loginErrorMessage,
			int32(rl.GetScreenWidth()/2-100),
			int32(rl.GetScreenHeight()/2+40),
			20,
			rl.Red)
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.loginSceneData.emailBounds) {
			w.loginSceneData.emailFocus = true
			w.loginSceneData.passwordFocus = false
		} else if rl.CheckCollisionPointRec(mousePos, w.loginSceneData.passwordBounds) {
			w.loginSceneData.emailFocus = false
			w.loginSceneData.passwordFocus = true
		} else {
			w.loginSceneData.emailFocus = false
			w.loginSceneData.passwordFocus = false
		}
	}
	gui.TextBox(w.loginSceneData.emailBounds, &w.loginSceneData.emailTXT, 64, w.loginSceneData.emailFocus)
	//TODO: secret password input box
	gui.TextBox(w.loginSceneData.passwordBounds, &w.loginSceneData.passwordTXT, 64, w.loginSceneData.passwordFocus)

}
