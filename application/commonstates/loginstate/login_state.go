package loginstate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

// LOGIN STATE
type LoginScene struct {
	cfg                  *utils.SharedConfig
	stateManager         *statesmanager.StateManager
	scheduler            utils.Scheduler
	loginButton          component.Button
	emailInput           component.InputBox
	passwordInput        component.InputBox
	isLoginButtonPressed bool
	errorSection         ErrorSection
}
type ErrorSection struct {
	errorPopup        component.Popup
	loginErrorMessage string
}

func (l *LoginScene) LoginSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	l.cfg = cfg
	l.stateManager = state
	l.Reset()

	l.loginButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2),
		200, 40,
	), "Login", false)

	l.emailInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	))

	l.passwordInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	))
	l.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2.0-100.0),
		float32(rl.GetScreenHeight()/2.0+40.0),
		200,
		100), &l.errorSection.loginErrorMessage)
}

func (l *LoginScene) UpdateLoginState() {
	l.scheduler.Update(float64(rl.GetFrameTime()))
	l.emailInput.Update()
	l.passwordInput.Update()
	l.isLoginButtonPressed = l.loginButton.Update()

	if l.isLoginButtonPressed {
		//i have to do check services and then mark this somehow to show that user can use it

		//and maybe use this to not make other request we have to wait if goruitne change this var to false and then???
		//thiis is to set own presence to cut all messsage service from app
		l.Login()
	}

}

// TOOD during the login is time to connect to all other services e.g messageService
// (we also are connecting to server but this is what we are doing first because without server
// we cant live but without messages services we can

// TODO (user can be offline but also error with messageSrvies is reason to set his status to offline
// even he is logged in
func (l *LoginScene) RenderLoginState() {

	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)

	l.errorSection.errorPopup.Render()
	//TODO: secret password inboxInput box
	l.emailInput.Render()
	l.passwordInput.Render()
	l.loginButton.Render()

}
