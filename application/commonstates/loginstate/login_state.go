package loginstate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

// LOGIN STATE
type LoginScene struct {
	cfg          *utils.SharedConfig
	stateManager *statesmanager.StateManager
	scheduler    utils.Scheduler
	loginSection LoginSection
	errorSection ErrorSection
}
type LoginSection struct {
	loginButton          component.Button
	emailInput           component.InputBox
	passwordInput        component.InputBox
	isLoginButtonPressed bool
}
type ErrorSection struct {
	errorPopup        component.Popup
	loginErrorMessage string
}

func (l *LoginScene) LoginSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	l.cfg = cfg
	l.stateManager = state
	l.Reset()

	l.loginSection.loginButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2),
		200, 50,
	), "LOGIN", false)

	l.loginSection.emailInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	))

	l.loginSection.passwordInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	))
	l.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(component.WithBgColor(utils.POPUPERRORBG)), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2.0-100.0),
		float32(rl.GetScreenHeight()/2.0+60.0),
		200,
		40), &l.errorSection.loginErrorMessage)
}

func (l *LoginScene) UpdateLoginState() {
	l.scheduler.Update(float64(rl.GetFrameTime()))
	l.loginSection.emailInput.Update()
	l.loginSection.passwordInput.Update()
	l.loginSection.isLoginButtonPressed = l.loginSection.loginButton.Update()

	if l.loginSection.isLoginButtonPressed {
		l.Login()
	}

}

func (l *LoginScene) RenderLoginState() {
	rl.ClearBackground(utils.LOGINBGCOLOR)
	rl.DrawText("LOGIN PAGE", int32(rl.GetScreenWidth()/2)-rl.MeasureText("LOGIN PAGE", 80)/2, 50, 80, rl.DarkGray)
	rl.DrawText("remote command system", int32(rl.GetScreenWidth()/2)-rl.MeasureText("LOGIN PAGE", 50)/4, 110, 50, rl.DarkGray)

	l.errorSection.errorPopup.Render()
	l.loginSection.emailInput.Render()
	l.loginSection.passwordInput.Render()
	l.loginSection.loginButton.Render()

}
