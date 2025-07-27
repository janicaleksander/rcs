package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"time"
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
	nameBounds rl.Rectangle
	nameTXT    string
	nameFocus  bool

	//user choose dropdown
	users              []string
	userDropdownBounds rl.Rectangle
	userActiveElement  int32
	userDropDownFocus  int32
	scrollIndex        int32
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
	resp := w.ctx.Request(w.serverPID, &Proto.GetUserAboveLVL{}, 5*time.Second)

	val, err := resp.Result()
	if err != nil {
		w.createUnitScene.isSetupError = true
	}

	if v, ok := val.(*Proto.GetUserAboveLVL); ok {
		w.createUnitScene.users = make([]string, 0, 64)
		w.createUnitScene.users = append(w.createUnitScene.users, "Choose user by his ID")
		for _, user := range v.Users {
			w.createUnitScene.users = append(w.createUnitScene.users, user.Id)
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
		position: rl.NewRectangle(
			10,
			float32(w.height-50),
			150,
			50),
		text: "GO BACK",
	}

	//accept button
	w.createUnitScene.acceptButton = Button{
		position: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2+50),
			200, 40,
		),
		text: "Accept ",
	}

	//name of unit
	w.createUnitScene.nameBounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-100),
		200, 40,
	)

	//dropdown with users
	w.createUnitScene.scrollIndex = 0
	w.createUnitScene.userActiveElement = 0
	w.createUnitScene.userDropdownBounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		float32(rl.GetScreenHeight()/2-60),
		240, 80,
	)
}

func (w *Window) updateCreateUnitState() {
	if w.createUnitScene.isSetupError {
		w.createUnitScene.errorMessage = "Setup error, can't do this now!"
	}
	if !w.createUnitScene.isSetupError && gui.Button(w.createUnitScene.acceptButton.position, w.createUnitScene.acceptButton.text) {
		w.createUnitScene.Reset()
		name := w.createUnitScene.nameTXT
		user := w.createUnitScene.userActiveElement
		if len(name) <= 0 || user <= 0 {
			w.createUnitScene.isCreateError = true
			w.createUnitScene.errorMessage = "Zero length error"
		} else {
			//user can be only in one unit in the same time -> error
			resp := w.ctx.Request(w.serverPID, &Proto.CreateUnit{
				Name:         name,
				IsConfigured: false,
				UserID:       w.createUnitScene.users[user],
			}, time.Second*5)
			// TODO check if its not out of bounds error
			//TODO do sth after success
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
		if rl.CheckCollisionPointRec(mousePos, w.createUnitScene.nameBounds) {
			w.createUnitScene.nameFocus = true
			w.createUnitScene.userDropDownFocus = 0
		} else if rl.CheckCollisionPointRec(mousePos, w.createUnitScene.userDropdownBounds) {
			w.createUnitScene.nameFocus = false
			w.createUnitScene.userDropDownFocus = 1
		} else {
			w.createUnitScene.nameFocus = false
			w.createUnitScene.userDropDownFocus = 0
		}
	}
	//rl.DrawText("Name of unit")
	gui.TextBox(w.createUnitScene.nameBounds, &w.createUnitScene.nameTXT, 64, w.createUnitScene.nameFocus)

	//rl.DrawText("Dropdown with users")
	gui.ListViewEx(w.createUnitScene.userDropdownBounds, w.createUnitScene.users, &w.createUnitScene.scrollIndex, &w.createUnitScene.userActiveElement, w.createUnitScene.userDropDownFocus)

	//go back button
	if gui.Button(w.createUnitScene.backButton.position, w.createUnitScene.backButton.text) {
		w.goSceneBack()
	}
}
