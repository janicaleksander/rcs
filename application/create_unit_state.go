package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO add better description to component in GUI
type CreateUnitScene struct {
	scheduler      utils.Scheduler
	backButton     component.Button
	newUnitSection NewUnitSection
	errorSection   ErrorSection2
	infoSection    InfoSection
}

// TODO CHANGE THIS NAME TO ONLY ERROR SECTION AFTER PACKAGE REFACTOR
type ErrorSection2 struct {
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
	usersDropdown   ListSlider
}

func (s *CreateUnitScene) Reset() {
	s.errorSection.isSetupError = false
	s.errorSection.isCreateError = false
	s.infoSection.isInfoMessage = false
	s.infoSection.infoMessage = ""
	s.errorSection.errorMessage = ""
}

func (w *Window) createUnitSceneSetup() {
	w.createUnitScene.Reset()

	//TODO get proper lvl value
	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetUserAboveLVL{Lvl: -1}))
	if err != nil {
		//context deadline exceeded
		//do sth with that
		w.createUnitScene.errorSection.isSetupError = true
	}

	if v, ok := res.(*proto.UsersAboveLVL); ok {
		w.createUnitScene.newUnitSection.usersDropdown.strings = make([]string, 0, 64)
		w.createUnitScene.newUnitSection.usersDropdown.strings = append(w.createUnitScene.newUnitSection.usersDropdown.strings,
			"Choose user by his ID")
		for _, user := range v.Users {
			w.createUnitScene.newUnitSection.usersDropdown.strings = append(w.createUnitScene.newUnitSection.usersDropdown.strings,
				user.Id+"\n"+user.Email)
		}
	} else {
		w.createUnitScene.errorSection.isSetupError = true
	}

	//name of unit
	w.createUnitScene.newUnitSection.nameInput = *component.NewInputBox(component.NewInputBoxConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-100),
		200, 40,
	))

	//dropdown with users
	w.createUnitScene.newUnitSection.usersDropdown.idxScroll = 0
	w.createUnitScene.newUnitSection.usersDropdown.idxActiveElement = 0
	w.createUnitScene.newUnitSection.usersDropdown.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		float32(rl.GetScreenHeight()/2-60),
		240, 80,
	)
	w.createUnitScene.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(w.width/2),
		float32(w.height-20),
		100, 20), &w.createUnitScene.errorSection.errorMessage)

	w.createUnitScene.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(w.width/2),
		float32(w.height-20),
		100, 20), &w.createUnitScene.infoSection.infoMessage)

	//accept button
	w.createUnitScene.newUnitSection.acceptButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2+50),
		200, 40,
	), "Accept", false)

	//go back from creating unit
	w.createUnitScene.backButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		10,
		float32(w.height-50),
		150,
		50), "Go back", false)

}

func (w *Window) updateCreateUnitState() {
	w.createUnitScene.scheduler.Update(float64(rl.GetFrameTime()))

	//go back button
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.createUnitScene.newUnitSection.usersDropdown.bounds) {
			w.createUnitScene.newUnitSection.usersDropdown.focus = 1
		} else {
			w.createUnitScene.newUnitSection.usersDropdown.focus = 0
		}
	}

	w.createUnitScene.newUnitSection.nameInput.Update()
	w.createUnitScene.newUnitSection.isAcceptPressed = w.createUnitScene.newUnitSection.acceptButton.Update()
	if w.createUnitScene.backButton.Update() {
		w.goSceneBack()
		return
	}
	//TODO add other from render
	if w.createUnitScene.errorSection.isSetupError {
		w.createUnitScene.errorSection.errorMessage = "Setup error, can't do this now!"
		return
	}

	if w.createUnitScene.newUnitSection.isAcceptPressed {
		w.CreateUnit()
	}

}
func (w *Window) renderCreateUnitState() {
	rl.DrawText(`Create unit Menu Page`, 50, 50, 20, rl.DarkGray)
	w.createUnitScene.newUnitSection.nameInput.Render()
	w.createUnitScene.newUnitSection.acceptButton.Render()
	w.createUnitScene.backButton.Render()
	gui.ListViewEx(
		w.createUnitScene.newUnitSection.usersDropdown.bounds,
		w.createUnitScene.newUnitSection.usersDropdown.strings,
		&w.createUnitScene.newUnitSection.usersDropdown.idxScroll,
		&w.createUnitScene.newUnitSection.usersDropdown.idxActiveElement,
		w.createUnitScene.newUnitSection.usersDropdown.focus,
	)

}

//scene HC unit info  dodac guziki w polu desxc podzielnic na 4 kwadraty i np mapa, opis mzoe urzadzenia itd itd
// a scena dla dowodcy jednego unityu moze to samo ale dla 1 unitu
//moze jakies panele ze mozna sobie potem w innym oknie tworzyc wlasnie np mapa + cos id

//a tam gdzie opis soldierow to dac guziki np ze wyslij wiadomosc, albo info o itd
