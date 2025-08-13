package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/User"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"

	"strconv"
)

type CreateUserScene struct {
	emailInput      component.InputBox
	passwordInput   component.InputBox
	rePasswordInput component.InputBox
	ruleLevelInput  component.InputBox
	nameInput       component.InputBox
	surnameInput    component.InputBox
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
	//set constse.g maxLength
	w.createUserScene.emailInput = *component.NewInputBox(
		component.NewConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2-300),
			200,
			60,
		), false)

	w.createUserScene.passwordInput = *component.NewInputBox(
		component.NewConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)-200,
			200,
			60), false)

	w.createUserScene.rePasswordInput = *component.NewInputBox(
		component.NewConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)-100,
			200,
			60), false)

	w.createUserScene.ruleLevelInput = *component.NewInputBox(
		component.NewConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2),
			200,
			60,
		), false)

	w.createUserScene.nameInput = *component.NewInputBox(
		component.NewConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)+100,
			200,
			60), false)

	w.createUserScene.surnameInput = *component.NewInputBox(
		component.NewConfig(),
		rl.NewRectangle(
			float32(w.width/2)-100,
			float32(w.height/2)+200,
			200,
			60,
		), false)

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
	w.createUserScene.emailInput.Update()
	w.createUserScene.passwordInput.Update()
	w.createUserScene.rePasswordInput.Update()
	w.createUserScene.ruleLevelInput.Update()
	w.createUserScene.nameInput.Update()
	w.createUserScene.surnameInput.Update()
	if gui.Button(w.createUserScene.acceptButton.bounds, w.createUserScene.acceptButton.text) {
		w.createUserScene.Reset()
		email := w.createUserScene.emailInput.GetText()
		password := w.createUserScene.passwordInput.GetText()
		rePassword := w.createUserScene.rePasswordInput.GetText()
		ruleLevel := w.createUserScene.ruleLevelInput.GetText()
		name := w.createUserScene.nameInput.GetText()
		surname := w.createUserScene.surnameInput.GetText()

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
		newUser := user.NewUser(email, password, int32(lvl), name, surname)
		resp := w.ctx.Request(w.serverPID, &proto.CreateUser{User: newUser}, utils.WaitTime)
		val, err := resp.Result()
		if err != nil {

			w.createUserScene.isError = true
			w.createUserScene.errorMessage = "Actor ctx error"
		}
		if _, ok := val.(*proto.AcceptCreateUser); ok {
			w.createUserScene.isAccept = true
			w.createUserScene.acceptMessage = "Created successfully"
			w.createUserScene.emailInput.Clear()
			w.createUserScene.passwordInput.Clear()
			w.createUserScene.rePasswordInput.Clear()
			w.createUserScene.ruleLevelInput.Clear()
			w.createUserScene.nameInput.Clear()
			w.createUserScene.surnameInput.Clear()
		}
		if _, ok := val.(*proto.DenyCreateUser); ok {
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
		w.createUserScene.emailInput.Deactivate()
		w.createUserScene.passwordInput.Deactivate()
		w.createUserScene.rePasswordInput.Deactivate()
		w.createUserScene.ruleLevelInput.Deactivate()
		w.createUserScene.nameInput.Deactivate()
		w.createUserScene.surnameInput.Deactivate()

		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.emailInput.Bounds) {
			w.createUserScene.emailInput.Active()
		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.passwordInput.Bounds) {
			w.createUserScene.passwordInput.Active()

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.rePasswordInput.Bounds) {
			w.createUserScene.rePasswordInput.Active()

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.ruleLevelInput.Bounds) {
			w.createUserScene.ruleLevelInput.Active()

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.nameInput.Bounds) {
			w.createUserScene.nameInput.Active()

		}
		if rl.CheckCollisionPointRec(mousePos, w.createUserScene.surnameInput.Bounds) {
			w.createUserScene.surnameInput.Active()
		}
	}
	w.createUserScene.emailInput.Render()
	w.createUserScene.passwordInput.Render()
	w.createUserScene.rePasswordInput.Render()
	w.createUserScene.ruleLevelInput.Render()
	w.createUserScene.nameInput.Render()
	w.createUserScene.surnameInput.Render()
}
