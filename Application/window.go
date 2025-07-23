package Application

import (
	"github.com/anthdm/hollywood/actor"
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
	HCMenuState
	CreateUnitState
)

type Window struct {
	sceneStack []GameState
	width      int
	height     int
	//server PID
	serverPID *actor.PID
	ctx       *actor.Context

	currentState GameState
	running      bool

	loginSceneData  LoginScene
	hcMenuSceneData HCMenuScene
	createUnitScene CreateUnitScene

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
	w.width = 1280
	w.height = 720
	rl.InitWindow(1280, 720, "BCS Application")
	rl.SetTargetFPS(60)
	w.sceneStack = append(w.sceneStack, LoginState)
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
	case HCMenuState:
		w.updateHCMenuState()
	case CreateUnitState:
		w.updateCreateUnitState()
	}

}
func (w *Window) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	switch w.currentState {
	case LoginState:
		w.renderLoginState()
	case HCMenuState:
		w.renderHCMenuState()
	case CreateUnitState:
		w.renderCreateUnitState()
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

func (w *Window) goSceneBack() {
	if len(w.sceneStack) > 1 {
		w.sceneStack = w.sceneStack[:len(w.sceneStack)-1]
		w.currentState = w.sceneStack[len(w.sceneStack)-1]
	}
}
