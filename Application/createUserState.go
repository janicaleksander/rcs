package Application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/User"
	"github.com/janicaleksander/bcs/Utils"

	"strconv"
)

type CreateUserScene struct {
	emailInput      InputField
	passwordInput   InputField
	rePasswordInput InputField
	ruleLevelInput  InputField
	nameInput       InputField
	surnameInput    InputField
	isError         bool
	errorMessage    string
	acceptButton    Button

	isAccept      bool
	acceptMessage string
	//maybe sobe checkbox
}

func (s *CreateUserScene) Reset() {
	s.isError = false
	s.errorMessage = ""

	s.isAccept = false
	s.acceptMessage = ""
}

func (w *Window) createUserSceneSetup() {
	w.createUserScene.Reset()
	w.createUserScene.emailInput = InputField{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2-300),
			200,
			60,
		),
		text:     "",
		focus:    false,
		textSize: 64,
	}
	w.createUserScene.passwordInput = InputField{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)-200,
			200,
			60),
		text:     "",
		focus:    false,
		textSize: 64,
	}
	w.createUserScene.rePasswordInput = InputField{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)-100,
			200,
			60),
		text:     "",
		focus:    false,
		textSize: 64,
	}
	w.createUserScene.ruleLevelInput = InputField{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2),
			200,
			60,
		),
		text:     "",
		focus:    false,
		textSize: 64,
	}
	w.createUserScene.nameInput = InputField{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)+100,
			200,
			60),
		text:     "",
		focus:    false,
		textSize: 64,
	}
	w.createUserScene.surnameInput = InputField{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)+200,
			200,
			60,
		),
		text:     "",
		focus:    false,
		textSize: 64,
	}
	w.createUserScene.acceptButton = Button{
		bounds: rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2+280),
			200,
			60),
		text: "ACCEPT",
	}
}

func (w *Window) updateCreateUserState() {
	if gui.Button(w.createUserScene.acceptButton.bounds, w.createUserScene.acceptButton.text) {
		w.createUserScene.Reset()
		email := w.createUserScene.emailInput.text
		password := w.createUserScene.passwordInput.text
		rePassword := w.createUserScene.rePasswordInput.text
		ruleLevel := w.createUserScene.ruleLevelInput.text
		name := w.createUserScene.nameInput.text
		surname := w.createUserScene.surnameInput.text

		//check inboxInput
		if len(email) <= 0 || len(password) <= 0 ||
			len(rePassword) <= 0 || len(ruleLevel) <= 0 ||
			len(name) <= 0 || len(surname) <= 0 {
			w.createUserScene.isError = true
			w.createUserScene.errorMessage = "Zero length error"
			return
		}
		lvl, err := strconv.Atoi(ruleLevel)
		// TODO curr max lvl
		if lvl > 5 || err != nil {
			w.createUserScene.isError = true
			w.createUserScene.errorMessage = "Bad ruleLVL inboxInput"
			return
		}
		newUser := User.NewUser(email, password, int32(lvl), name, surname)
		resp := w.ctx.Request(w.serverPID, &Proto.CreateUser{User: newUser}, Utils.WaitTime)
		val, err := resp.Result()
		if err != nil {

			w.createUserScene.isError = true
			w.createUserScene.errorMessage = "Actor ctx error"
		}
		if _, ok := val.(*Proto.AcceptCreateUser); ok {
			w.createUserScene.isAccept = true
			w.createUserScene.acceptMessage = "Created successfully"
			w.createUserScene.emailInput.text = ""
			w.createUserScene.passwordInput.text = ""
			w.createUserScene.rePasswordInput.text = ""
			w.createUserScene.ruleLevelInput.text = ""
			w.createUserScene.nameInput.text = ""
			w.createUserScene.surnameInput.text = ""
		}
		if _, ok := val.(*Proto.DenyCreateUser); ok {
			w.createUserScene.isError = true
			w.createUserScene.errorMessage = "DB deny!"
		}
	}

}

func (w *Window) renderCreateUserState() {
	//error
	if w.createUserScene.isError {
		rl.DrawText(w.createUserScene.errorMessage, 0, 0, 32, rl.Red)
	}
	if w.createUserScene.isAccept {
		rl.DrawText(w.createUserScene.acceptMessage, 0, 0, 32, rl.Green)

	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		w.createUserScene.emailInput.focus = false
		w.createUserScene.passwordInput.focus = false
		w.createUserScene.rePasswordInput.focus = false
		w.createUserScene.ruleLevelInput.focus = false
		w.createUserScene.nameInput.focus = false
		w.createUserScene.surnameInput.focus = false

		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.emailInput.bounds) {
			w.createUserScene.emailInput.focus = true
		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.passwordInput.bounds) {
			w.createUserScene.passwordInput.focus = true

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.rePasswordInput.bounds) {
			w.createUserScene.rePasswordInput.focus = true

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.ruleLevelInput.bounds) {
			w.createUserScene.ruleLevelInput.focus = true

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.nameInput.bounds) {
			w.createUserScene.nameInput.focus = true

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.surnameInput.bounds) {
			w.createUserScene.surnameInput.focus = true
		}
	}
	gui.TextBox(w.createUserScene.emailInput.bounds, &w.createUserScene.emailInput.text, w.createUserScene.emailInput.textSize, w.createUserScene.emailInput.focus)
	gui.TextBox(w.createUserScene.passwordInput.bounds, &w.createUserScene.passwordInput.text, w.createUserScene.passwordInput.textSize, w.createUserScene.passwordInput.focus)
	gui.TextBox(w.createUserScene.rePasswordInput.bounds, &w.createUserScene.rePasswordInput.text, w.createUserScene.rePasswordInput.textSize, w.createUserScene.rePasswordInput.focus)
	gui.TextBox(w.createUserScene.ruleLevelInput.bounds, &w.createUserScene.ruleLevelInput.text, w.createUserScene.ruleLevelInput.textSize, w.createUserScene.ruleLevelInput.focus)
	gui.TextBox(w.createUserScene.nameInput.bounds, &w.createUserScene.nameInput.text, w.createUserScene.nameInput.textSize, w.createUserScene.nameInput.focus)
	gui.TextBox(w.createUserScene.surnameInput.bounds, &w.createUserScene.surnameInput.text, w.createUserScene.surnameInput.textSize, w.createUserScene.surnameInput.focus)

}
