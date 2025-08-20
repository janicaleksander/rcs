package createuserstate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/component"
	"github.com/janicaleksander/bcs/utils"
)

type CreateUserScene struct {
	cfg            *utils.SharedConfig
	stateManager   *statesmanager.StateManager
	scheduler      utils.Scheduler
	backButton     component.Button
	newUserSection NewUserSection
	errorSection   ErrorSection
	infoSection    InfoSection
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
type ErrorSection struct {
	errorPopup   component.Popup
	errorMessage string
}
type InfoSection struct {
	infoPopup     component.Popup
	acceptMessage string
}

func (c *CreateUserScene) CreateUserSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	c.cfg = cfg
	c.stateManager = state
	c.Reset()
	c.scheduler.Update(float64(rl.GetFrameTime()))

	//set constse.g maxLength
	c.newUserSection.emailInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2)-100,
			float32(rl.GetScreenHeight()/2-300),
			200,
			60,
		))

	c.newUserSection.passwordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2)-100,
			float32(rl.GetScreenHeight()/2)-200,
			200,
			60))

	c.newUserSection.rePasswordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2)-100,
			float32(rl.GetScreenHeight()/2)-100,
			200,
			60))

	c.newUserSection.ruleLevelInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2)-100,
			float32(rl.GetScreenHeight()/2),
			200,
			60,
		))

	c.newUserSection.nameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2)-100,
			float32(rl.GetScreenHeight()/2)+100,
			200,
			60))

	c.newUserSection.surnameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2)-100,
			float32(rl.GetScreenHeight()/2)+200,
			200,
			60,
		))

	c.newUserSection.acceptButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2)-100,
		float32(rl.GetScreenHeight()/2+280),
		200,
		60), "Accept", false)

	c.backButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2)-150,
		float32(rl.GetScreenHeight()/2+380),
		200,
		60), "Go back", false)

	c.errorSection.errorPopup = *component.NewPopup(
		component.NewPopupConfig(),
		rl.NewRectangle(float32(rl.GetScreenWidth()/2)-150,
			float32(rl.GetScreenHeight()/2+480),
			200,
			60),
		&c.errorSection.errorMessage)
	c.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(),
		rl.NewRectangle(float32(rl.GetScreenWidth()/2)-150,
			float32(rl.GetScreenHeight()/2+480),
			200,
			60),
		&c.infoSection.acceptMessage)
}

func (c *CreateUserScene) UpdateCreateUserState() {

	c.newUserSection.emailInput.Update()
	c.newUserSection.passwordInput.Update()
	c.newUserSection.rePasswordInput.Update()
	c.newUserSection.ruleLevelInput.Update()
	c.newUserSection.nameInput.Update()
	c.newUserSection.surnameInput.Update()
	c.newUserSection.isAcceptPressed = c.newUserSection.acceptButton.Update()

	if c.backButton.Update() {
		c.stateManager.Add(statesmanager.GoBackState)
		return
	}

	if c.newUserSection.isAcceptPressed {
		c.CreateUser()
	}

}

func (c *CreateUserScene) RenderCreateUserState() {

	c.newUserSection.emailInput.Render()
	c.newUserSection.passwordInput.Render()
	c.newUserSection.rePasswordInput.Render()
	c.newUserSection.ruleLevelInput.Render()
	c.newUserSection.nameInput.Render()
	c.newUserSection.surnameInput.Render()
	c.newUserSection.acceptButton.Render()
	c.backButton.Render()
}
