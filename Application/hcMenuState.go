package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
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

	openInboxButton Button
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
	w.hcMenuScene.unitRectangle = rl.NewRectangle(
		float32(padding),
		height,
		float32(width/4),
		0.75*float32(w.height),
	)
	//user
	w.hcMenuScene.userRectangle = rl.NewRectangle(
		w.hcMenuScene.unitRectangle.X+w.hcMenuScene.unitRectangle.Width+float32(padding),
		w.hcMenuScene.unitRectangle.Y,
		w.hcMenuScene.unitRectangle.Width,
		w.hcMenuScene.unitRectangle.Height,
	)
	//inbox
	w.hcMenuScene.inboxRectangle = rl.NewRectangle(
		w.hcMenuScene.userRectangle.X+w.hcMenuScene.userRectangle.Width+float32(padding),
		w.hcMenuScene.userRectangle.Y,
		float32(width/2),
		w.hcMenuScene.unitRectangle.Height,
	)

	//create unit
	w.hcMenuScene.createUnitButton = Button{
		bounds: rl.NewRectangle(
			10+w.hcMenuScene.unitRectangle.X,
			10+w.hcMenuScene.unitRectangle.Y,
			200,
			40,
		),
		text: "Create unit",
	}

	//info about units
	w.hcMenuScene.infoUnitButton = Button{
		bounds: rl.NewRectangle(
			10+w.hcMenuScene.unitRectangle.X,
			80+w.hcMenuScene.unitRectangle.Y,
			200,
			40,
		),
		text: "Units info",
	}

	// create user
	w.hcMenuScene.createUserButton = Button{
		bounds: rl.NewRectangle(w.hcMenuScene.userRectangle.X, w.hcMenuScene.userRectangle.Y, 200, 40),
		text:   "Create User",
	}

	//info about users
	w.hcMenuScene.infoUserButton = Button{
		bounds: rl.NewRectangle(w.hcMenuScene.userRectangle.X, w.hcMenuScene.userRectangle.Y+80, 200, 40),
		text:   "Users info",
	}
	w.hcMenuScene.openInboxButton = Button{
		bounds: rl.NewRectangle(w.hcMenuScene.inboxRectangle.X, w.hcMenuScene.inboxRectangle.Y+80, 200, 40),
		text:   "Open inbox",
	}

}
func (w *Window) updateHCMenuState() {

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
		int32(w.hcMenuScene.unitRectangle.X),
		int32(w.hcMenuScene.unitRectangle.Y),
		int32(w.hcMenuScene.unitRectangle.Width),
		int32(w.hcMenuScene.unitRectangle.Height), rl.Green)

	//user
	rl.DrawRectangle(
		int32(w.hcMenuScene.userRectangle.X),
		int32(w.hcMenuScene.userRectangle.Y),
		int32(w.hcMenuScene.userRectangle.Width),
		int32(w.hcMenuScene.userRectangle.Height), rl.Green)

	//inbox
	rl.DrawRectangle(
		int32(w.hcMenuScene.inboxRectangle.X),
		int32(w.hcMenuScene.inboxRectangle.Y),
		int32(w.hcMenuScene.inboxRectangle.Width),
		int32(w.hcMenuScene.inboxRectangle.Height), rl.Green)

	//button create unit
	if gui.Button(w.hcMenuScene.createUnitButton.bounds, w.hcMenuScene.createUnitButton.text) {
		w.updatePresence(w.ctx.PID(), &Proto.PresencePlace{
			Place: &Proto.PresencePlace_Outbox{
				Outbox: &Proto.Outbox{},
			},
		})
		w.createUnitSceneSetup()
		w.currentState = CreateUnitState
		w.sceneStack = append(w.sceneStack, CreateUnitState)
	}

	//button info units
	if gui.Button(w.hcMenuScene.infoUnitButton.bounds, w.hcMenuScene.infoUnitButton.text) {
		w.updatePresence(w.ctx.PID(), &Proto.PresencePlace{
			Place: &Proto.PresencePlace_Outbox{
				Outbox: &Proto.Outbox{}}})
		w.infoUnitSceneSetup()
		w.currentState = InfoUnitState
		w.sceneStack = append(w.sceneStack, InfoUnitState)
	}
	//button create user
	if gui.Button(w.hcMenuScene.createUserButton.bounds, w.hcMenuScene.createUserButton.text) {
		w.updatePresence(w.ctx.PID(), &Proto.PresencePlace{
			Place: &Proto.PresencePlace_Outbox{
				Outbox: &Proto.Outbox{}}})
		w.createUserSceneSetup()
		w.currentState = CreateUserState
		w.sceneStack = append(w.sceneStack, CreateUserState)
	}
	//info user button
	if gui.Button(w.hcMenuScene.infoUserButton.bounds, w.hcMenuScene.infoUserButton.text) {
		w.updatePresence(w.ctx.PID(), &Proto.PresencePlace{
			Place: &Proto.PresencePlace_Outbox{
				Outbox: &Proto.Outbox{}}})

		w.InfoUserSceneSetup()
		w.currentState = InfoUserState
		w.sceneStack = append(w.sceneStack, InfoUserState)
	}
	//open inbox
	if gui.Button(w.hcMenuScene.openInboxButton.bounds, w.hcMenuScene.openInboxButton.text) {
		//TODO
		w.updatePresence(w.ctx.PID(), &Proto.PresencePlace{
			Place: &Proto.PresencePlace_Inbox{
				Inbox: &Proto.Inbox{}}})

		w.setupInboxScene()
		w.currentState = InboxState
		w.sceneStack = append(w.sceneStack, InboxState)
	}

}
