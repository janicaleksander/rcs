package application

import (
	"fmt"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"runtime"
)

type GameState int

const (
	LoginState GameState = iota
)

type Window struct {
	currentState GameState
	running      bool

	loginSceneData LoginScene
	//menuSceneData MenuScene

	// data from actor through chan's

	Done chan bool
	Test chan *Proto.Payload
}

func NewWindow(test chan *Proto.Payload) *Window {
	return &Window{
		Done: make(chan bool),
		Test: test,
	}
}

func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	rl.InitWindow(1920, 1080, "BCS application")
	rl.SetTargetFPS(60)
	w.currentState = LoginState

	//init login state
	w.loginSceneSetup()
	//
	w.running = true

}
func (w *Window) quit() {
	//maybe to quit some assets
}

func (w *Window) drawScene() {}
func (w *Window) input() {

}

func (w *Window) update() {
	w.running = !rl.WindowShouldClose()
	switch w.currentState {
	case LoginState:
		w.updateLoginState()
	}

}
func (w *Window) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	switch w.currentState {
	case LoginState:
		w.renderLoginState()
	}
	rl.EndDrawing()

}
func (w *Window) RunWindow() {

	defer func() {
		rl.CloseWindow()
		close(w.Done)
	}()

	w.setup()
	for w.running {
		w.input()
		w.update()
		w.render()
	}
	w.quit()
}

// LOGIN STATE
type LoginScene struct {
	loginButton Button
	emailTXT    string
	passwordTXT string
}

func (w *Window) loginSceneSetup() {
	w.loginSceneData.loginButton = Button{
		position: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2-100),
			200, 40,
		),
		text: "Login",
	}
}
func (w *Window) updateLoginState() {
	//this is already rendering this button on the screen
	if gui.Button(w.loginSceneData.loginButton.position, w.loginSceneData.loginButton.text) {
		fmt.Println("ok")
	}
}
func (w *Window) renderLoginState() {}
