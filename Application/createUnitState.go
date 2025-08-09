package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Utils"
)

// TODO add better description to components in GUI
type CreateUnitScene struct {
	isSetupError  bool
	isCreateError bool
	errorMessage  string
	errorPosition Position

	isInfoMessage bool
	infoMessage   string
	infoPosition  Position

	backButton   Button
	acceptButton Button

	//name of new unit
	nameInput InputField

	//user choose dropdown
	usersDropdown ListSlider
}

func (s *CreateUnitScene) Reset() {
	s.isSetupError = false
	s.isCreateError = false
	s.isInfoMessage = false
	s.infoMessage = ""
	s.errorMessage = ""
}

func (w *Window) createUnitSceneSetup() {
	w.createUnitScene.Reset()

	//get user slice from DB
	//TODO get proper lvl value
	resp := w.ctx.Request(w.serverPID, &Proto.GetUserAboveLVL{Lvl: -1}, Utils.WaitTime)

	val, err := resp.Result()
	if err != nil {
		w.createUnitScene.isSetupError = true
	}

	if v, ok := val.(*Proto.UsersAboveLVL); ok {
		w.createUnitScene.usersDropdown.strings = make([]string, 0, 64)
		w.createUnitScene.usersDropdown.strings = append(w.createUnitScene.usersDropdown.strings, "Choose user by his ID")
		for _, user := range v.Users {
			w.createUnitScene.usersDropdown.strings = append(w.createUnitScene.usersDropdown.strings, user.Id+"\n"+user.Email)
		}
	} else {
		w.createUnitScene.isSetupError = true
	}
	w.createUnitScene.errorPosition = Position{
		int32(w.width / 2),
		int32(w.height - 20),
	}
	w.createUnitScene.infoPosition = Position{
		int32(w.width / 2),
		int32(w.height - 20),
	}

	//go back from creating unit
	w.createUnitScene.backButton = Button{
		bounds: rl.NewRectangle(
			10,
			float32(w.height-50),
			150,
			50),
		text: "GO BACK",
	}

	//accept button
	w.createUnitScene.acceptButton = Button{
		bounds: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2+50),
			200, 40,
		),
		text: "Accept ",
	}

	//name of unit
	w.createUnitScene.nameInput.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-100),
		200, 40,
	)

	//dropdown with users
	w.createUnitScene.usersDropdown.idxScroll = 0
	w.createUnitScene.usersDropdown.idxActiveElement = 0
	w.createUnitScene.usersDropdown.bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		float32(rl.GetScreenHeight()/2-60),
		240, 80,
	)
}

func (w *Window) updateCreateUnitState() {
	if w.createUnitScene.isSetupError {
		w.createUnitScene.errorMessage = "Setup error, can't do this now!"
	}
	if !w.createUnitScene.isSetupError && gui.Button(w.createUnitScene.acceptButton.bounds, w.createUnitScene.acceptButton.text) {
		w.createUnitScene.Reset()
		name := w.createUnitScene.nameInput.text
		user := w.createUnitScene.usersDropdown.idxActiveElement
		if len(name) <= 0 || user <= 0 {
			w.createUnitScene.isCreateError = true
			w.createUnitScene.errorMessage = "Zero length error"
		} else {
			//user can be only in one unit in the same time -> error
			resp := w.ctx.Request(w.serverPID, &Proto.CreateUnit{
				Name:         name,
				IsConfigured: false,
				UserID:       w.createUnitScene.usersDropdown.strings[user],
			}, Utils.WaitTime)
			val, err := resp.Result()
			if err != nil {
				w.createUnitScene.isCreateError = true
				w.createUnitScene.errorMessage = err.Error()
			}
			if _, ok := val.(*Proto.AcceptCreateUnit); ok {
				w.createUnitScene.isInfoMessage = true
				w.createUnitScene.infoMessage = "Success!"
			}
			if v, ok := val.(*Proto.DenyCreateUnit); ok {
				w.createUnitScene.isCreateError = true
				w.createUnitScene.errorMessage = v.Info

			}
		}
	}

}
func (w *Window) renderCreateUnitState() {
	rl.DrawText(`Create unit Menu Page`, 50, 50, 20, rl.DarkGray)

	if w.createUnitScene.isCreateError || w.createUnitScene.isSetupError {
		rl.DrawText(w.createUnitScene.errorMessage,
			w.createUnitScene.errorPosition.x,
			w.createUnitScene.errorPosition.y,
			20,
			rl.Red)
	}
	if w.createUnitScene.isInfoMessage {
		rl.DrawText(w.createUnitScene.infoMessage,
			w.createUnitScene.infoPosition.x,
			w.createUnitScene.infoPosition.y,
			20,
			rl.Green)
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.createUnitScene.usersDropdown.bounds) {
			w.createUnitScene.nameInput.focus = true
			w.createUnitScene.usersDropdown.focus = 0
		} else if rl.CheckCollisionPointRec(mousePos, w.createUnitScene.usersDropdown.bounds) {
			w.createUnitScene.nameInput.focus = false
			w.createUnitScene.usersDropdown.focus = 1
		} else {
			w.createUnitScene.nameInput.focus = false
			w.createUnitScene.usersDropdown.focus = 0
		}
	}
	//rl.DrawText("Name of unit")
	gui.TextBox(w.createUnitScene.nameInput.bounds, &w.createUnitScene.nameInput.text, 64, w.createUnitScene.nameInput.focus)

	//rl.DrawText("Dropdown with users")
	gui.ListViewEx(w.createUnitScene.usersDropdown.bounds, w.createUnitScene.usersDropdown.strings, &w.createUnitScene.usersDropdown.idxScroll, &w.createUnitScene.usersDropdown.idxActiveElement, w.createUnitScene.usersDropdown.focus)

	//go back button
	if gui.Button(w.createUnitScene.backButton.bounds, w.createUnitScene.backButton.text) {
		w.goSceneBack()
	}
}

//scene HC unit info  dodac guziki w polu desxc podzielnic na 4 kwadraty i np mapa, opis mzoe urzadzenia itd itd
// a scena dla dowodcy jednego unityu moze to samo ale dla 1 unitu
//moze jakies panele ze mozna sobie potem w innym oknie tworzyc wlasnie np mapa + cos id

//a tam gdzie opis soldierow to dac guziki np ze wyslij wiadomosc, albo info o itd
