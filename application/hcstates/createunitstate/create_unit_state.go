package createunitstate

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

// TODO add better description to component in GUI
type CreateUnitScene struct {
	cfg            *utils.SharedConfig
	stateManager   *statesmanager.StateManager
	scheduler      utils.Scheduler
	backButton     component.Button
	newUnitSection NewUnitSection
	errorSection   ErrorSection
	infoSection    InfoSection
}

// TODO CHANGE THIS NAME TO ONLY ERROR SECTION AFTER PACKAGE REFACTOR
type ErrorSection struct {
	isSetupError  bool
	isCreateError bool
	errorMessage  string
	errorPopup    component.Popup
}
type InfoSection struct {
	isInfoMessage bool
	infoMessage   string
	infoPopup     component.Popup
}

type NewUnitSection struct {
	acceptButton    component.Button
	isAcceptPressed bool
	nameInput       component.InputBox
	usersDropdown   component.ListSlider
}

func (c *CreateUnitScene) CreateUnitSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	c.cfg = cfg
	c.stateManager = state
	c.Reset()
	c.FetchUsers()
	//name of unit
	c.newUnitSection.nameInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-100),
		200, 40,
	))

	//dropdown with users
	c.newUnitSection.usersDropdown.IdxScroll = 0
	c.newUnitSection.usersDropdown.IdxActiveElement = 0
	c.newUnitSection.usersDropdown.Bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		float32(rl.GetScreenHeight()/2-60),
		240, 80,
	)
	c.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2),
		float32(rl.GetScreenHeight()-20),
		100, 20), &c.errorSection.errorMessage)

	c.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2),
		float32(rl.GetScreenHeight()-20),
		100, 20), &c.infoSection.infoMessage)

	//accept button
	c.newUnitSection.acceptButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2+50),
		200, 40,
	), "Accept", false)

	//go back from creating unit
	c.backButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		10,
		float32(rl.GetScreenHeight()-50),
		150,
		50), "Go back", false)

}

func (c *CreateUnitScene) UpdateCreateUnitState() {
	c.scheduler.Update(float64(rl.GetFrameTime()))

	//go back button
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, c.newUnitSection.usersDropdown.Bounds) {
			c.newUnitSection.usersDropdown.Focus = 1
		} else {
			c.newUnitSection.usersDropdown.Focus = 0
		}
	}

	c.newUnitSection.nameInput.Update()
	c.newUnitSection.isAcceptPressed = c.newUnitSection.acceptButton.Update()
	if c.backButton.Update() {
		c.stateManager.Add(statesmanager.GoBackState)
		return
	}
	//TODO add other from render
	if c.errorSection.isSetupError {
		c.errorSection.errorMessage = "Setup error, can't do this now!"
		return
	}

	if c.newUnitSection.isAcceptPressed {
		c.CreateUnit()
	}

}
func (c *CreateUnitScene) RenderCreateUnitState() {
	rl.DrawText(`Create unit Menu Page`, 50, 50, 20, rl.DarkGray)
	c.newUnitSection.nameInput.Render()
	c.newUnitSection.acceptButton.Render()
	c.backButton.Render()
	gui.ListViewEx(
		c.newUnitSection.usersDropdown.Bounds,
		c.newUnitSection.usersDropdown.Strings,
		&c.newUnitSection.usersDropdown.IdxScroll,
		&c.newUnitSection.usersDropdown.IdxActiveElement,
		c.newUnitSection.usersDropdown.Focus,
	)

}

//scene HC unit info  dodac guziki w polu desxc podzielnic na 4 kwadraty i np mapa, opis mzoe urzadzenia itd itd
// a scena dla dowodcy jednego unityu moze to samo ale dla 1 unitu
//moze jakies panele ze mozna sobie potem w innym oknie tworzyc wlasnie np mapa + cos id

//a tam gdzie opis soldierow to dac guziki np ze wyslij wiadomosc, albo info o itd
