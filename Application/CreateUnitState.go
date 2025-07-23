package Application

import (
	"fmt"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"time"
)

type CreateUnitScene struct {
	setupError bool

	isLoginError      bool
	loginErrorMessage string

	backButton   Button
	acceptButton Button

	//name of new unit
	nameBounds rl.Rectangle
	nameTXT    string
	nameFocus  bool

	//user choose dropdown
	users              []string
	userDropdownBounds rl.Rectangle
	userDropDownList   string
	userActiveElement  int32
	userDropDownFocus  int32
	scrollIndex        int32
}

func (w *Window) createUnitSceneSetup() {

	//get user slice from DB
	resp := w.ctx.Request(w.serverPID, &Proto.GetUserAboveLVL{}, 5*time.Second)

	val, err := resp.Result()
	if err != nil {
		w.createUnitScene.setupError = true
	}

	if v, ok := val.(*Proto.GetUserAboveLVL); ok {
		w.createUnitScene.users = make([]string, 0, 64)
		for _, user := range v.Users {
			w.createUnitScene.users = append(w.createUnitScene.users, user.Id)
		}
		w.createUnitScene.users = append(w.createUnitScene.users, "*---*---*")
	} else {
		w.createUnitScene.setupError = true
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
			float32(rl.GetScreenHeight()/2+110),
			200, 40,
		),
		text: "Accept ",
	}

	//name of unit
	w.createUnitScene.nameBounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-70),
		200, 40,
	)

	//dropdown with users
	w.createUnitScene.userDropdownBounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2+15),
		220, 80,
	)
}

func (w *Window) updateCreateUnitState() {

	if !w.createUnitScene.setupError {
		w.createUnitScene.loginErrorMessage = "Setup error, can't do this now!"
	} else if gui.Button(w.createUnitScene.acceptButton.position, w.createUnitScene.acceptButton.text) {
		fmt.Println(w.createUnitScene.users)
	}

}
func (w *Window) renderCreateUnitState() {
	rl.DrawText("Create unit Menu Page", 50, 50, 20, rl.DarkGray)

	if w.createUnitScene.isLoginError || w.createUnitScene.setupError {
		rl.DrawText(w.createUnitScene.loginErrorMessage,
			int32(w.width/2),
			int32(w.height-20),
			64,
			rl.Red)
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
