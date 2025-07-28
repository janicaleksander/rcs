package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type HCMenuScene struct {
	//backButton       Button
	profileRectangle rl.Rectangle
	unitRectangle    rl.Rectangle
	userRectangle    rl.Rectangle
	inboxRectangle   rl.Rectangle

	createUnitButton Button
	infoUnitButton   Button

	createUserButton Button
	infoUserButton   Button
}

func (w *Window) menuHCSceneSetup() {
	//profile
	profileSide := 100
	w.hcMenuSceneData.profileRectangle = rl.NewRectangle(
		float32(w.width)-float32(profileSide),
		25,
		float32(profileSide),
		float32(profileSide),
	)

	padding := 30
	height := 0.25 * float32(w.height)
	width := w.width - 2*padding

	//unit
	w.hcMenuSceneData.unitRectangle = rl.NewRectangle(
		float32(padding),
		height,
		float32(width/4),
		0.75*float32(w.height),
	)
	//user
	w.hcMenuSceneData.userRectangle = rl.NewRectangle(
		w.hcMenuSceneData.unitRectangle.X+w.hcMenuSceneData.unitRectangle.Width+float32(padding),
		w.hcMenuSceneData.unitRectangle.Y,
		w.hcMenuSceneData.unitRectangle.Width,
		w.hcMenuSceneData.unitRectangle.Height,
	)
	//inbox
	w.hcMenuSceneData.inboxRectangle = rl.NewRectangle(
		w.hcMenuSceneData.userRectangle.X+w.hcMenuSceneData.userRectangle.Width+float32(padding),
		w.hcMenuSceneData.userRectangle.Y,
		float32(width/2),
		w.hcMenuSceneData.unitRectangle.Height,
	)

	//create unit
	w.hcMenuSceneData.createUnitButton = Button{
		bounds: rl.NewRectangle(
			10+w.hcMenuSceneData.unitRectangle.X,
			10+w.hcMenuSceneData.unitRectangle.Y,
			200,
			40,
		),
		text: "Create unit",
	}

	//info about units
	w.hcMenuSceneData.infoUnitButton = Button{
		bounds: rl.NewRectangle(
			10+w.hcMenuSceneData.unitRectangle.X,
			80+w.hcMenuSceneData.unitRectangle.Y,
			200,
			40,
		),
		text: "Units info",
	}

	// create user
	w.hcMenuSceneData.createUserButton = Button{
		bounds: rl.Rectangle{},
		text:   "",
	}

	//info about users
	w.hcMenuSceneData.infoUserButton = Button{
		bounds: rl.Rectangle{},
		text:   "",
	}

}
func (w *Window) updateHCMenuState() {

}
func (w *Window) renderHCMenuState() {
	rl.DrawText("HC Menu Page", 50, 50, 20, rl.DarkGray)
	//profile
	rl.DrawRectangle(
		int32(w.hcMenuSceneData.profileRectangle.X),
		int32(w.hcMenuSceneData.profileRectangle.Y),
		int32(w.hcMenuSceneData.profileRectangle.Width),
		int32(w.hcMenuSceneData.profileRectangle.Height), rl.Green)

	//unit
	rl.DrawRectangle(
		int32(w.hcMenuSceneData.unitRectangle.X),
		int32(w.hcMenuSceneData.unitRectangle.Y),
		int32(w.hcMenuSceneData.unitRectangle.Width),
		int32(w.hcMenuSceneData.unitRectangle.Height), rl.Green)

	//user
	rl.DrawRectangle(
		int32(w.hcMenuSceneData.userRectangle.X),
		int32(w.hcMenuSceneData.userRectangle.Y),
		int32(w.hcMenuSceneData.userRectangle.Width),
		int32(w.hcMenuSceneData.userRectangle.Height), rl.Green)

	//inbox
	rl.DrawRectangle(
		int32(w.hcMenuSceneData.inboxRectangle.X),
		int32(w.hcMenuSceneData.inboxRectangle.Y),
		int32(w.hcMenuSceneData.inboxRectangle.Width),
		int32(w.hcMenuSceneData.inboxRectangle.Height), rl.Green)

	//button create unit
	if gui.Button(w.hcMenuSceneData.createUnitButton.bounds, w.hcMenuSceneData.createUnitButton.text) {
		w.createUnitSceneSetup()
		w.currentState = CreateUnitState
		w.sceneStack = append(w.sceneStack, CreateUnitState)
	}

	//button info units
	if gui.Button(w.hcMenuSceneData.infoUnitButton.bounds, w.hcMenuSceneData.infoUnitButton.text) {
		w.infoUnitSceneSetup()
		w.currentState = InfoUnitState
		w.sceneStack = append(w.sceneStack, InfoUnitState)
	}

}
