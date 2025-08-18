package application

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/utils"
)

// LOGIN STATE
type LoginScene struct {
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

func (w *Window) loginSceneSetup() {
	w.loginScene.Reset()

	w.loginScene.loginButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2),
		200, 40,
	), "Login", false)

	w.loginScene.emailInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	), false)

	w.loginScene.passwordInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	), false)
	w.loginScene.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2.0-100.0),
		float32(rl.GetScreenHeight()/2.0+40.0),
		200,
		100), &w.loginScene.errorSection.loginErrorMessage)
}

func (w *Window) updateLoginState() {
	w.loginScene.scheduler.Update(float64(rl.GetFrameTime()))

	w.loginScene.emailInput.Update()
	w.loginScene.passwordInput.Update()
	w.loginScene.isLoginButtonPressed = w.loginScene.loginButton.Update()
	if w.loginScene.isLoginButtonPressed {
		//i have to do check services and then mark this somehow to show that user can use it

		//and maybe use this to not make other request we have to wait if goruitne change this var to false and then???
		//thiis is to set own presence to cut all messsage service from app
		w.Login()
	}

}

// TOOD during the login is time to connect to all other services e.g messageService
// (we also are connecting to server but this is what we are doing first because without server
// we cant live but without messages services we can

// TODO (user can be offline but also error with messageSrvies is reason to set his status to offline
// even he is logged in
func (w *Window) renderLoginState() {

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		w.loginScene.emailInput.Deactivate()
		w.loginScene.passwordInput.Deactivate()

		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.loginScene.emailInput.Bounds) {
			w.loginScene.emailInput.Active()
			w.loginScene.passwordInput.Deactivate()
		} else if rl.CheckCollisionPointRec(mousePos, w.loginScene.passwordInput.Bounds) {
			w.loginScene.emailInput.Deactivate()
			w.loginScene.passwordInput.Active()
		}
	}
	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)

	w.loginScene.errorSection.errorPopup.Render()
	//TODO: secret password inboxInput box
	w.loginScene.emailInput.Render()
	w.loginScene.passwordInput.Render()
	w.loginScene.loginButton.Render()

}
