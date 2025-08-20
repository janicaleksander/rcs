package loginstate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/statesmanager"
	component2 "github.com/janicaleksander/bcs/component"
	"github.com/janicaleksander/bcs/utils"
)

// LOGIN STATE
type LoginScene struct {
	cfg                  *utils.SharedConfig
	stateManager         *statesmanager.StateManager
	scheduler            utils.Scheduler
	loginButton          component2.Button
	emailInput           component2.InputBox
	passwordInput        component2.InputBox
	isLoginButtonPressed bool
	errorSection         ErrorSection
}
type ErrorSection struct {
	errorPopup        component2.Popup
	loginErrorMessage string
}

func (l *LoginScene) LoginSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	l.cfg = cfg
	l.stateManager = state
	l.Reset()

	l.loginButton = *component2.NewButton(component2.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2),
		200, 40,
	), "Login", false)

	l.emailInput = *component2.NewInputBox(component2.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	))

	l.passwordInput = *component2.NewInputBox(component2.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	))
	l.errorSection.errorPopup = *component2.NewPopup(component2.NewPopupConfig(), rl.NewRectangle(
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
