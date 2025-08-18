package application

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
)

type HCMenuScene struct {
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
	userRectangle       rl.Rectangle
	createUserButton    component.Button
	isCreateUserPressed bool
	infoUserButton      component.Button
	isInfoUserPressed   bool
}

type InboxSection struct {
	inboxRectangle     rl.Rectangle
	openInboxButton    component.Button
	isOpenInboxPressed bool
}

func (w *Window) menuHCSceneSetup() {
	//profile
	profileSide := 100
	w.hcMenuScene.profileRectangle = rl.NewRectangle(
		float32(w.width)-float32(profileSide),
		25,
		float32(profileSide),
		float32(profileSide),
	)

	padding := 30
	height := 0.25 * float32(w.height)
	width := w.width - 2*padding

	//unit
	w.hcMenuScene.unitSection.unitRectangle = rl.NewRectangle(
		float32(padding),
		height,
		float32(width/4),
		0.75*float32(w.height),
	)
	//user
	w.hcMenuScene.userSection.userRectangle = rl.NewRectangle(
		w.hcMenuScene.unitSection.unitRectangle.X+w.hcMenuScene.unitSection.unitRectangle.Width+float32(padding),
		w.hcMenuScene.unitSection.unitRectangle.Y,
		w.hcMenuScene.unitSection.unitRectangle.Width,
		w.hcMenuScene.unitSection.unitRectangle.Height,
	)
	//inbox
	w.hcMenuScene.inboxSection.inboxRectangle = rl.NewRectangle(
		w.hcMenuScene.userSection.userRectangle.X+w.hcMenuScene.userSection.userRectangle.Width+float32(padding),
		w.hcMenuScene.userSection.userRectangle.Y,
		float32(width/2),
		w.hcMenuScene.unitSection.unitRectangle.Height,
	)

	//create unit
	w.hcMenuScene.unitSection.createUnitButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		10+w.hcMenuScene.unitSection.unitRectangle.X,
		10+w.hcMenuScene.unitSection.unitRectangle.Y,
		200,
		40,
	), "Create unit", false)

	//info about units
	w.hcMenuScene.unitSection.infoUnitButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		10+w.hcMenuScene.unitSection.unitRectangle.X,
		80+w.hcMenuScene.unitSection.unitRectangle.Y,
		200,
		40,
	), "Units info", false)

	// create user
	w.hcMenuScene.userSection.createUserButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			w.hcMenuScene.userSection.userRectangle.X,
			w.hcMenuScene.userSection.userRectangle.Y,
			200,
			40),
		"Create User", false)

	//info about users
	w.hcMenuScene.userSection.infoUserButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			w.hcMenuScene.userSection.userRectangle.X,
			w.hcMenuScene.userSection.userRectangle.Y+80,
			200,
			40),
		"Users info", false)

	w.hcMenuScene.inboxSection.openInboxButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.hcMenuScene.inboxSection.inboxRectangle.X,
		w.hcMenuScene.inboxSection.inboxRectangle.Y+80,
		200,
		40), "Open inbox", false)

}
func (w *Window) updateHCMenuState() {
	w.hcMenuScene.unitSection.isCreateUnitPressed = w.hcMenuScene.unitSection.createUnitButton.Update()
	w.hcMenuScene.unitSection.isInfoUnitPressed = w.hcMenuScene.unitSection.infoUnitButton.Update()
	w.hcMenuScene.userSection.isCreateUserPressed = w.hcMenuScene.userSection.createUserButton.Update()
	w.hcMenuScene.userSection.isInfoUserPressed = w.hcMenuScene.userSection.infoUserButton.Update()
	w.hcMenuScene.inboxSection.isOpenInboxPressed = w.hcMenuScene.inboxSection.openInboxButton.Update()

	if w.hcMenuScene.unitSection.isInfoUnitPressed {
		w.infoUnitSceneSetup()
		w.currentState = InfoUnitState
		w.sceneStack = append(w.sceneStack, InfoUnitState)
	}
	//button create user
	if w.hcMenuScene.unitSection.isCreateUnitPressed {
		w.createUnitSceneSetup()
		w.currentState = CreateUnitState
		w.sceneStack = append(w.sceneStack, CreateUnitState)
	}
	if w.hcMenuScene.userSection.isCreateUserPressed {
		w.createUserSceneSetup()
		w.currentState = CreateUserState
		w.sceneStack = append(w.sceneStack, CreateUserState)
	}
	//info user button
	if w.hcMenuScene.userSection.isInfoUserPressed {
		w.InfoUserSceneSetup()
		w.currentState = InfoUserState
		w.sceneStack = append(w.sceneStack, InfoUserState)
	}
	//open inbox
	if w.hcMenuScene.inboxSection.isOpenInboxPressed {
		w.setupInboxScene()
		w.currentState = InboxState
		w.sceneStack = append(w.sceneStack, InboxState)
	}
}
func (w *Window) renderHCMenuState() {
	rl.DrawText("HC Menu Page", 50, 50, 20, rl.DarkGray)
	//profile
	rl.DrawRectangle(
		int32(w.hcMenuScene.profileRectangle.X),
		int32(w.hcMenuScene.profileRectangle.Y),
		int32(w.hcMenuScene.profileRectangle.Width),
		int32(w.hcMenuScene.profileRectangle.Height), rl.Green)

	//unit
	rl.DrawRectangle(
		int32(w.hcMenuScene.unitSection.unitRectangle.X),
		int32(w.hcMenuScene.unitSection.unitRectangle.Y),
		int32(w.hcMenuScene.unitSection.unitRectangle.Width),
		int32(w.hcMenuScene.unitSection.unitRectangle.Height), rl.Green)

	//user
	rl.DrawRectangle(
		int32(w.hcMenuScene.userSection.userRectangle.X),
		int32(w.hcMenuScene.userSection.userRectangle.Y),
		int32(w.hcMenuScene.userSection.userRectangle.Width),
		int32(w.hcMenuScene.userSection.userRectangle.Height), rl.Green)

	//inbox
	rl.DrawRectangle(
		int32(w.hcMenuScene.inboxSection.inboxRectangle.X),
		int32(w.hcMenuScene.inboxSection.inboxRectangle.Y),
		int32(w.hcMenuScene.inboxSection.inboxRectangle.Width),
		int32(w.hcMenuScene.inboxSection.inboxRectangle.Height), rl.Green)

	w.hcMenuScene.unitSection.createUnitButton.Render()
	w.hcMenuScene.unitSection.infoUnitButton.Render()
	w.hcMenuScene.userSection.createUserButton.Render()
	w.hcMenuScene.userSection.infoUserButton.Render()
	w.hcMenuScene.inboxSection.openInboxButton.Render()

}
