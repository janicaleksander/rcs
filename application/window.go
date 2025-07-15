package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"runtime"
)

type GameState int

const (
	WaitToLogin int = 3
)

const (
	LoginState GameState = iota
)

type Window struct {
	currentState GameState
	running      bool

	loginSceneData LoginScene
	//menuSceneData MenuScene

	// data from actor through chan's

	Done        chan bool
	ChLoginUser chan *Proto.LoginUser
}

func NewWindow(login chan *Proto.LoginUser) *Window {
	return &Window{
		Done:        make(chan bool),
		ChLoginUser: login,
	}
}

func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(1280, 720, "BCS application")
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
	emailFocus  bool

	passwordTXT   string
	passwordFocus bool
}

func (w *Window) loginSceneSetup() {
	w.loginSceneData.loginButton = Button{
		position: rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-100),
			float32(rl.GetScreenHeight()/2),
			200, 40,
		),
		text: "Login",
	}
}
func (w *Window) updateLoginState() {
	//this is already rendering this button on the screen

	//block multiple pressed button before receive respond, for const time LoginTime

	if gui.Button(w.loginSceneData.loginButton.position, w.loginSceneData.loginButton.text) {
		//valid login credentials -> send to server loginUser, wait for response
		email := w.loginSceneData.emailTXT
		pwd := w.loginSceneData.passwordTXT
		w.ChLoginUser <- &Proto.LoginUser{
			Email:    email,
			Password: pwd,
		}
	}

}
func (w *Window) renderLoginState() {
	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)
	emailBounds := rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-140),
		200, 30,
	)
	passwordBounds := rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-100),
		float32(rl.GetScreenHeight()/2-80),
		200, 30,
	)
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, emailBounds) {
			w.loginSceneData.emailFocus = true
			w.loginSceneData.passwordFocus = false
		} else if rl.CheckCollisionPointRec(mousePos, passwordBounds) {
			w.loginSceneData.emailFocus = false
			w.loginSceneData.passwordFocus = true
		} else {
			w.loginSceneData.emailFocus = false
			w.loginSceneData.passwordFocus = false
		}
	}
	gui.TextBox(emailBounds, &w.loginSceneData.emailTXT, 64, w.loginSceneData.emailFocus)
	gui.TextBox(passwordBounds, &w.loginSceneData.passwordTXT, 64, w.loginSceneData.passwordFocus)

}
