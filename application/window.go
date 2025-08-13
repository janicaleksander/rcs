package application

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/anthdm/hollywood/actor"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

type GameState int

const (
	WIDTH  = 1280
	HEIGHT = 720
)

const (
	LoginState GameState = iota
	HCMenuState
	CreateUnitState
	InfoUnitState
	CreateUserState
	InfoUserState
	InboxState
)

type Window struct {
	sceneStack []GameState
	width      int
	height     int
	//server PID
	serverPID         *actor.PID
	messageServicePID *actor.PID //TODO
	ctx               *actor.Context

	currentState GameState
	running      bool

	messageServiceError bool
	//and other errors this is good idea
	loginSceneScene LoginScene
	hcMenuScene     HCMenuScene
	createUnitScene CreateUnitScene
	infoUnitScene   InfoUnitScene
	createUserScene CreateUserScene
	infoUserScene   InfoUserScene
	inboxScene      InboxScene

	Done chan bool
}

func NewWindowActor(w *Window, serverPID, messageServicePID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		w.serverPID = serverPID
		w.messageServicePID = messageServicePID
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
	case actor.Initialized:
		utils.Logger.Info("Actor initialized")
	case actor.Started:
		utils.Logger.Info("Window actor started")
	case actor.Stopped:
		utils.Logger.Info("Actor stopped")
	case *proto.Ping:
		ctx.Respond(&proto.Pong{})
	case *proto.DeliverMessage:
		fmt.Println("ODEBRALEM,", msg.Message)
		w.inboxScene.messageChan <- msg.Message
	default:
		utils.Logger.Warn("server got unknown errorMessage", reflect.TypeOf(msg).String())
	}
}
func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	w.width = WIDTH
	w.height = HEIGHT
	rl.InitWindow(int32(w.width), int32(w.height), "BCS application")
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
	case InboxState:
		w.updateInboxState()

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
	case InboxState:
		w.renderInboxState()
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

type Position struct {
	x int32
	y int32
}
