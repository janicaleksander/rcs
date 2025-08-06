package Application

import (
	"reflect"
	"runtime"
	"time"

	"github.com/anthdm/hollywood/actor"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
)

type GameState int

const (
	WIDTH                  = 1280
	HEIGHT                 = 720
	WaitTime time.Duration = 3 * time.Second
)

const (
	LoginState GameState = iota
	HCMenuState
	CreateUnitState
	InfoUnitState
	CreateUserState
	InfoUserState
)

type Position struct {
	x int32
	y int32
}

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
	infoUnitScene   InfoUnitScene
	createUserScene CreateUserScene
	infoUserScene   InfoUserScene

	Done chan bool
}

func NewWindowActor(w *Window, serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		w.serverPID = serverPID
		return w
	}
}

func NewWindow() *Window {
	return &Window{
		Done: make(chan bool, 1),
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
	case *Proto.Ping:
		ctx.Respond(&Proto.Pong{})
	default:
		Server.Logger.Warn("Server got unknown errorMessage", reflect.TypeOf(msg).String())
	}
}
func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	w.width = WIDTH
	w.height = HEIGHT
	rl.InitWindow(int32(w.width), int32(w.height), "BCS Application")
	rl.SetTargetFPS(60)
	w.sceneStack = append(w.sceneStack, LoginState)

	//init login state
	w.currentState = LoginState
	w.loginSceneSetup()
	w.running = true

}
func (w *Window) quit() {
	//maybe to quit some assets
}

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
	case InfoUnitState:
		w.updateInfoUnitState()
	case CreateUserState:
		w.updateCreateUserState()
	case InfoUserState:
		w.updateInfoUserState()

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
	case InfoUnitState:
		w.renderInfoUnitState()
	case CreateUserState:
		w.renderCreateUserState()
	case InfoUserState:
		w.renderInfoUserState()
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
