package application

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/utils"
)

type CreateUserScene struct {
	scheduler      utils.Scheduler
	backButton     component.Button
	newUserSection NewUserSection
	errorSection   ErrorSection3
	infoSection    InfoSection2
}
type NewUserSection struct {
	emailInput      component.InputBox
	passwordInput   component.InputBox
	rePasswordInput component.InputBox
	ruleLevelInput  component.InputBox
	nameInput       component.InputBox
	surnameInput    component.InputBox
	acceptButton    component.Button
	isAcceptPressed bool
}
type ErrorSection3 struct {
	errorPopup   component.Popup
	errorMessage string
}
type InfoSection2 struct {
	infoPopup     component.Popup
	acceptMessage string
}

func (w *Window) createUserSceneSetup() {
	w.createUserScene.Reset()
	w.createUserScene.scheduler.Update(float64(rl.GetFrameTime()))

	//set constse.g maxLength
	w.createUserScene.newUserSection.emailInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2-300),
			200,
			60,
		))

	w.createUserScene.newUserSection.passwordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)-200,
			200,
			60))

	w.createUserScene.newUserSection.rePasswordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)-100,
			200,
			60))

	w.createUserScene.newUserSection.ruleLevelInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2),
			200,
			60,
		))

	w.createUserScene.newUserSection.nameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)+100,
			200,
			60))

	w.createUserScene.newUserSection.surnameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)+200,
			200,
			60,
		))

	w.createUserScene.newUserSection.acceptButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(w.width/2)-100,
		float32(w.height/2+280),
		200,
		60), "Accept", false)

	w.createUserScene.backButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(w.width/2)-150,
		float32(w.height/2+380),
		200,
		60), "Go back", false)

	w.createUserScene.errorSection.errorPopup = *component.NewPopup(
		component.NewPopupConfig(),
		rl.NewRectangle(float32(w.width/2)-150,
			float32(w.height/2+480),
			200,
			60),
		&w.createUserScene.errorSection.errorMessage)
	w.createUserScene.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(),
		rl.NewRectangle(float32(w.width/2)-150,
			float32(w.height/2+480),
			200,
			60),
		&w.createUserScene.infoSection.acceptMessage)
}

func (w *Window) updateCreateUserState() {

	w.createUserScene.newUserSection.emailInput.Update()
	w.createUserScene.newUserSection.passwordInput.Update()
	w.createUserScene.newUserSection.rePasswordInput.Update()
	w.createUserScene.newUserSection.ruleLevelInput.Update()
	w.createUserScene.newUserSection.nameInput.Update()
	w.createUserScene.newUserSection.surnameInput.Update()
	w.createUserScene.newUserSection.isAcceptPressed = w.createUserScene.newUserSection.acceptButton.Update()

	if w.createUserScene.backButton.Update() {
		w.goSceneBack()
		return
	}

	if w.createUserScene.newUserSection.isAcceptPressed {
		w.CreateUser()
	}

}

func (w *Window) renderCreateUserState() {

	w.createUserScene.newUserSection.emailInput.Render()
	w.createUserScene.newUserSection.passwordInput.Render()
	w.createUserScene.newUserSection.rePasswordInput.Render()
	w.createUserScene.newUserSection.ruleLevelInput.Render()
	w.createUserScene.newUserSection.nameInput.Render()
	w.createUserScene.newUserSection.surnameInput.Render()
	w.createUserScene.newUserSection.acceptButton.Render()
	w.createUserScene.backButton.Render()
}
