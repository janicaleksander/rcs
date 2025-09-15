package hcmenustate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

type HCMenuScene struct {
	cfg *utils.SharedConfig
	//scheduler
	stateManager *statesmanager.StateManager
	unitSection  UnitSection
	userSection  UserSection
	inboxSection InboxSection
}

type UnitSection struct {
	unitRectangle       rl.Rectangle
	createUnitButton    component.Button
	isCreateUnitPressed bool
	infoUnitButton      component.Button
	isInfoUnitPressed   bool
}

type UserSection struct {
	userRectangle         rl.Rectangle
	createUserButton      component.Button
	isCreateUserPressed   bool
	infoUserButton        component.Button
	isInfoUserPressed     bool
	createDeviceButton    component.Button
	isCreateDevicePressed bool
}

type InboxSection struct {
	inboxRectangle     rl.Rectangle
	openInboxButton    component.Button
	isOpenInboxPressed bool
}

func (h *HCMenuScene) MenuHCSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	h.cfg = cfg
	h.stateManager = state
	//profile

	padding := 30
	height := 0.25 * float32(rl.GetScreenHeight())
	width := rl.GetScreenWidth() - 2*padding

	//unit
	h.unitSection.unitRectangle = rl.NewRectangle(
		100+float32(padding),
		height,
		float32(width/4),
		0.75*float32(rl.GetScreenHeight()),
	)
	//user
	h.userSection.userRectangle = rl.NewRectangle(
		h.unitSection.unitRectangle.X+h.unitSection.unitRectangle.Width+float32(padding),
		h.unitSection.unitRectangle.Y,
		h.unitSection.unitRectangle.Width,
		h.unitSection.unitRectangle.Height,
	)
	//inbox
	h.inboxSection.inboxRectangle = rl.NewRectangle(
		h.userSection.userRectangle.X+h.userSection.userRectangle.Width+float32(padding),
		h.userSection.userRectangle.Y,
		h.unitSection.unitRectangle.Width,
		h.unitSection.unitRectangle.Height,
	)

	//create unit
	h.unitSection.createUnitButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		h.unitSection.unitRectangle.X+(h.unitSection.unitRectangle.Width/2)-float32(100),
		float32(padding)+h.unitSection.unitRectangle.Y,
		200,
		50,
	), "CREATE UNIT", false)

	//info about units
	h.unitSection.infoUnitButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		h.unitSection.createUnitButton.Bounds.X,
		float32(padding)+h.unitSection.createUnitButton.Bounds.Y+h.unitSection.createUnitButton.Bounds.Height,
		200,
		50,
	), "INFO", false)

	// create user
	h.userSection.createUserButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			h.userSection.userRectangle.X+(h.userSection.userRectangle.Width/2)-float32(100),
			float32(padding)+h.userSection.userRectangle.Y,
			200,
			50),
		"CREATE USER", false)

	//info about users
	//TODO add to user info delete device
	h.userSection.infoUserButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			h.userSection.createUserButton.Bounds.X,
			float32(padding)+h.userSection.createUserButton.Bounds.Y+h.userSection.createUserButton.Bounds.Height,
			200,
			50),
		"INFO", false)
	//create and assign device
	h.userSection.createDeviceButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			h.userSection.createUserButton.Bounds.X,
			float32(padding)+h.userSection.infoUserButton.Bounds.Y+h.userSection.infoUserButton.Bounds.Height,
			200,
			50),
		"ADD DEVICE", false)

	//open inbox
	h.inboxSection.openInboxButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		h.inboxSection.inboxRectangle.X+(h.inboxSection.inboxRectangle.Width/2)-float32(100),
		float32(padding)+h.inboxSection.inboxRectangle.Y,
		200,
		50), "OPEN INBOX", false)

}
func (h *HCMenuScene) UpdateHCMenuState() {
	h.unitSection.isCreateUnitPressed = h.unitSection.createUnitButton.Update()
	h.unitSection.isInfoUnitPressed = h.unitSection.infoUnitButton.Update()
	h.userSection.isCreateUserPressed = h.userSection.createUserButton.Update()
	h.userSection.isInfoUserPressed = h.userSection.infoUserButton.Update()
	h.userSection.isCreateDevicePressed = h.userSection.createDeviceButton.Update()
	h.inboxSection.isOpenInboxPressed = h.inboxSection.openInboxButton.Update()

	if h.unitSection.isInfoUnitPressed {
		h.stateManager.Add(statesmanager.InfoUnitState)
	} else if h.unitSection.isCreateUnitPressed {
		h.stateManager.Add(statesmanager.CreateUnitState)
	} else if h.userSection.isCreateUserPressed {
		h.stateManager.Add(statesmanager.CreateUserState)
	} else if h.userSection.isInfoUserPressed {
		h.stateManager.Add(statesmanager.InfoUserState)
	} else if h.userSection.isCreateDevicePressed {
		h.stateManager.Add(statesmanager.CreateDeviceState)
	} else if h.inboxSection.isOpenInboxPressed {
		h.stateManager.Add(statesmanager.InboxState)
	} else if h.inboxSection.isOpenInboxPressed {
		h.stateManager.Add(statesmanager.InboxState)
	}
}
func (h *HCMenuScene) RenderHCMenuState() {
	rl.ClearBackground(utils.HCMENUBG)
	rl.DrawText("MENU PAGE", int32(rl.GetScreenWidth()/2)-rl.MeasureText("MENU PAGE", 50)/2, 50, 50, rl.DarkGray)
	//unit
	rl.DrawRectangle(
		int32(h.unitSection.unitRectangle.X),
		int32(h.unitSection.unitRectangle.Y),
		int32(h.unitSection.unitRectangle.Width),
		int32(h.unitSection.unitRectangle.Height), utils.HCPARTSBG)
	rl.DrawText(
		"UNIT SECTION",
		int32(h.unitSection.unitRectangle.X)+int32(h.unitSection.unitRectangle.Width)/2-rl.MeasureText("UNIT SECTION", 20)/2,
		int32(h.unitSection.unitRectangle.Y)+int32(h.unitSection.unitRectangle.Height)-50,
		20,
		rl.Black,
	)
	//user
	rl.DrawRectangle(
		int32(h.userSection.userRectangle.X),
		int32(h.userSection.userRectangle.Y),
		int32(h.userSection.userRectangle.Width),
		int32(h.userSection.userRectangle.Height), utils.HCPARTSBG)
	rl.DrawText(
		"USER SECTION",
		int32(h.userSection.userRectangle.X)+int32(h.userSection.userRectangle.Width)/2-rl.MeasureText("USER SECTION", 20)/2,
		int32(h.userSection.userRectangle.Y)+int32(h.userSection.userRectangle.Height)-50,
		20,
		rl.Black,
	)
	//inbox
	rl.DrawRectangle(
		int32(h.inboxSection.inboxRectangle.X),
		int32(h.inboxSection.inboxRectangle.Y),
		int32(h.inboxSection.inboxRectangle.Width),
		int32(h.inboxSection.inboxRectangle.Height), utils.HCPARTSBG)
	rl.DrawText(
		"INBOX SECTION",
		int32(h.inboxSection.inboxRectangle.X)+int32(h.inboxSection.inboxRectangle.Width)/2-rl.MeasureText("INBOX SECTION", 20)/2,
		int32(h.inboxSection.inboxRectangle.Y)+int32(h.inboxSection.inboxRectangle.Height)-50,
		20,
		rl.Black,
	)

	h.unitSection.createUnitButton.Render()
	h.unitSection.infoUnitButton.Render()
	h.userSection.createUserButton.Render()
	h.userSection.infoUserButton.Render()
	h.userSection.createDeviceButton.Render()
	h.inboxSection.openInboxButton.Render()

}
