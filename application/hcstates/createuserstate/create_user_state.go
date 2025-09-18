package createuserstate

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

type CreateUserScene struct {
	cfg            *utils.SharedConfig
	stateManager   *statesmanager.StateManager
	backButton     component.Button
	newUserSection NewUserSection
	errorSection   ErrorSection
	infoSection    InfoSection
}
type NewUserSection struct {
	emailInput           component.InputBox
	passwordInput        component.InputBox
	rePasswordInput      component.InputBox
	ruleLevelToggleGroup component.ToggleGroup
	nameInput            component.InputBox
	surnameInput         component.InputBox
	acceptButton         component.Button
	isAcceptPressed      bool
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
	xPos := float32(rl.GetScreenWidth()/2) - 100
	yPos := float32(rl.GetScreenHeight()/2 - 300)
	c.newUserSection.emailInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			xPos,
			yPos,
			200,
			40,
		))
	yPos += 100
	c.newUserSection.passwordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			xPos,
			yPos,
			200,
			40))

	yPos += 100
	c.newUserSection.rePasswordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			xPos,
			yPos,
			200,
			40))

	c.newUserSection.ruleLevelToggleGroup = component.ToggleGroup{
		Selected: 0,
		Labels:   []string{"Mobile user", "Unit member", "Unit commander", "Head Commander"},
		Bounds: []rl.Rectangle{
			rl.NewRectangle(432, 370, 100, 40),
			rl.NewRectangle(532, 370, 100, 40),
			rl.NewRectangle(632, 370, 100, 40),
			rl.NewRectangle(732, 370, 100, 40),
		},
	}
	yPos += 200
	c.newUserSection.nameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			xPos,
			yPos,
			200,
			40))

	yPos += 100
	c.newUserSection.surnameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			xPos,
			yPos,
			200,
			40,
		))
	yPos += 80
	c.newUserSection.acceptButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		xPos,
		yPos,
		200,
		50), "ACCEPT", false)

	c.backButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(10),
		float32(rl.GetScreenHeight()-65),
		150,
		50), "GO BACK", false)

	popupRect := rl.NewRectangle(
		float32(rl.GetScreenWidth()/2)-100,
		float32(rl.GetScreenHeight()/2+280),
		200,
		50)
	c.errorSection.errorPopup = *component.NewPopup(
		component.NewPopupConfig(component.WithBgColor(utils.POPUPERRORBG)),
		popupRect,
		&c.errorSection.errorMessage)
	c.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(component.WithBgColor(utils.POPUPINFOBG)),
		popupRect,
		&c.infoSection.acceptMessage)
}

func (c *CreateUserScene) UpdateCreateUserState() {
	c.newUserSection.emailInput.Update()
	c.newUserSection.passwordInput.Update()
	c.newUserSection.rePasswordInput.Update()
	for i := range c.newUserSection.ruleLevelToggleGroup.Labels {
		toggleState := c.newUserSection.ruleLevelToggleGroup.Selected == i
		if gui.Toggle(
			c.newUserSection.ruleLevelToggleGroup.Bounds[i],
			c.newUserSection.ruleLevelToggleGroup.Labels[i],
			toggleState,
		) {
			c.newUserSection.ruleLevelToggleGroup.Selected = i
		}
	}
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
	rl.ClearBackground(utils.CREATEUSERBG)
	rl.DrawText("EMAIL", int32(c.newUserSection.emailInput.Bounds.X)-300, int32(c.newUserSection.emailInput.Bounds.Y)+int32(c.newUserSection.emailInput.Bounds.Height)/5, 25, rl.LightGray)
	rl.DrawText("PASSWORD", int32(c.newUserSection.passwordInput.Bounds.X)-300, int32(c.newUserSection.passwordInput.Bounds.Y)+int32(c.newUserSection.emailInput.Bounds.Height)/5, 25, rl.LightGray)
	rl.DrawText("RE-PASSWORD", int32(c.newUserSection.rePasswordInput.Bounds.X)-300, int32(c.newUserSection.rePasswordInput.Bounds.Y)+int32(c.newUserSection.emailInput.Bounds.Height)/5, 25, rl.LightGray)
	rl.DrawText("RULE LVL", int32(c.newUserSection.rePasswordInput.Bounds.X)-300, int32(c.newUserSection.ruleLevelToggleGroup.Bounds[0].Y)+int32(c.newUserSection.emailInput.Bounds.Height)/5, 25, rl.LightGray)
	rl.DrawText("NAME", int32(c.newUserSection.nameInput.Bounds.X)-300, int32(c.newUserSection.nameInput.Bounds.Y)+int32(c.newUserSection.emailInput.Bounds.Height)/5, 25, rl.LightGray)
	rl.DrawText("SURNAME", int32(c.newUserSection.surnameInput.Bounds.X)-300, int32(c.newUserSection.surnameInput.Bounds.Y)+int32(c.newUserSection.emailInput.Bounds.Height)/5, 25, rl.LightGray)
	c.newUserSection.emailInput.Render()
	c.newUserSection.passwordInput.Render()
	c.newUserSection.rePasswordInput.Render()
	c.newUserSection.nameInput.Render()
	c.newUserSection.surnameInput.Render()
	c.newUserSection.acceptButton.Render()
	c.errorSection.errorPopup.Render()
	c.infoSection.infoPopup.Render()
	c.backButton.Render()
}
