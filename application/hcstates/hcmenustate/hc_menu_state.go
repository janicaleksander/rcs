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
	stateManager     *statesmanager.StateManager
	profileRectangle rl.Rectangle
	unitSection      UnitSection
	userSection      UserSection
	inboxSection     InboxSection
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
	profileSide := 100
	h.profileRectangle = rl.NewRectangle(
		float32(rl.GetScreenWidth())-float32(profileSide),
		25,
		float32(profileSide),
		float32(profileSide),
	)

	padding := 30
	height := 0.25 * float32(rl.GetScreenHeight())
	width := rl.GetScreenWidth() - 2*padding

	//unit
	h.unitSection.unitRectangle = rl.NewRectangle(
		float32(padding),
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
		float32(width/2),
		h.unitSection.unitRectangle.Height,
	)

	//create unit
	h.unitSection.createUnitButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		10+h.unitSection.unitRectangle.X,
		10+h.unitSection.unitRectangle.Y,
		200,
		40,
	), "Create unit", false)

	//info about units
	h.unitSection.infoUnitButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		10+h.unitSection.unitRectangle.X,
		80+h.unitSection.unitRectangle.Y,
		200,
		40,
	), "Units info", false)

	// create user
	h.userSection.createUserButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			h.userSection.userRectangle.X,
			h.userSection.userRectangle.Y,
			200,
			40),
		"Create User", false)

	//info about users
	//TODO add to user info delete device
	h.userSection.infoUserButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			h.userSection.userRectangle.X,
			h.userSection.userRectangle.Y+80,
			200,
			40),
		"Users info", false)
	//create and assign device
	h.userSection.createDeviceButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			h.userSection.userRectangle.X,
			h.userSection.userRectangle.Y+160,
			200,
			40),
		"Add device", false)

	h.inboxSection.openInboxButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		h.inboxSection.inboxRectangle.X,
		h.inboxSection.inboxRectangle.Y+80,
		200,
		40), "Open inbox", false)

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
	rl.DrawText("HC Menu Page", 50, 50, 20, rl.DarkGray)
	//profile
	rl.DrawRectangle(
		int32(h.profileRectangle.X),
		int32(h.profileRectangle.Y),
		int32(h.profileRectangle.Width),
		int32(h.profileRectangle.Height), rl.Green)

	//unit
	rl.DrawRectangle(
		int32(h.unitSection.unitRectangle.X),
		int32(h.unitSection.unitRectangle.Y),
		int32(h.unitSection.unitRectangle.Width),
		int32(h.unitSection.unitRectangle.Height), rl.Green)

	//user
	rl.DrawRectangle(
		int32(h.userSection.userRectangle.X),
		int32(h.userSection.userRectangle.Y),
		int32(h.userSection.userRectangle.Width),
		int32(h.userSection.userRectangle.Height), rl.Green)

	//inbox
	rl.DrawRectangle(
		int32(h.inboxSection.inboxRectangle.X),
		int32(h.inboxSection.inboxRectangle.Y),
		int32(h.inboxSection.inboxRectangle.Width),
		int32(h.inboxSection.inboxRectangle.Height), rl.Green)

	h.unitSection.createUnitButton.Render()
	h.unitSection.infoUnitButton.Render()
	h.userSection.createUserButton.Render()
	h.userSection.infoUserButton.Render()
	h.userSection.createDeviceButton.Render()
	h.inboxSection.openInboxButton.Render()

}
