package Application

import (
	"github.com/anthdm/hollywood/actor"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"reflect"
	"runtime"
	"time"
)

type GameState int

const (
	WaitToLogin time.Duration = 3
)

const (
	LoginState GameState = iota
)

type Window struct {
	//server PID
	serverPID *actor.PID
	ctx       *actor.Context

	currentState GameState
	running      bool

	loginSceneData LoginScene
	//menuSceneData MenuScene

	Done chan bool
}

func NewWindowActor(w *Window) actor.Producer {
	return func() actor.Receiver {
		return w
	}
}

func NewWindow() *Window {
	return &Window{
		Done: make(chan bool, 1024),
	}
}

func (w *Window) Receive(ctx *actor.Context) {
	w.ctx = ctx
	switch msg := ctx.Message().(type) {
	case actor.Started:
		Server.Logger.Info("Window actor started")
	case actor.Initialized:
		Server.Logger.Info("Actor initialized")
	case actor.Stopped:
		Server.Logger.Info("Actor stopped")
	case *Proto.NeededServerConfiguration:
		w.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	default:
		Server.Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())
	}
}
func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(1280, 720, "BCS Application")
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

	isLoginError      bool
	loginErrorMessage string
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
		pid := &Proto.PID{
			Address: w.ctx.PID().GetAddress(),
			Id:      w.ctx.PID().GetID(),
		}
		resp := w.ctx.Request(w.serverPID, &Proto.LoginUser{
			Pid:      pid,
			Email:    email,
			Password: pwd,
		}, time.Second*WaitToLogin)
		val, err := resp.Result()
		//only if error this true
		w.loginSceneData.isLoginError = true
		if err != nil {
			w.loginSceneData.loginErrorMessage = err.Error()
		} else if msg, ok := val.(*Proto.AcceptLogin); ok {
			w.loginSceneData.loginErrorMessage = msg.Info
		} else if msg, ok := val.(*Proto.DenyLogin); ok {
			w.loginSceneData.loginErrorMessage = msg.Info
		} else {
			w.loginSceneData.loginErrorMessage = "Unknown response type"
		}

	}

}
func (w *Window) renderLoginState() {
	rl.DrawText("Login Page", 50, 50, 20, rl.DarkGray)

	if w.loginSceneData.isLoginError {
		rl.DrawText(w.loginSceneData.loginErrorMessage,
			int32(rl.GetScreenWidth()/2-100),
			int32(rl.GetScreenHeight()/2+40),
			20,
			rl.Red)

	}
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
