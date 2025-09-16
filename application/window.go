package application

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/anthdm/hollywood/actor"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/commonstates/createtaskstate"
	"github.com/janicaleksander/bcs/application/commonstates/inboxstate"
	"github.com/janicaleksander/bcs/application/commonstates/loginstate"
	"github.com/janicaleksander/bcs/application/hcstates/createdevicestate"
	"github.com/janicaleksander/bcs/application/hcstates/createunitstate"
	"github.com/janicaleksander/bcs/application/hcstates/createuserstate"
	"github.com/janicaleksander/bcs/application/hcstates/hcmenustate"
	"github.com/janicaleksander/bcs/application/hcstates/infounitstate"
	"github.com/janicaleksander/bcs/application/hcstates/infouserstate"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

const (
	WIDTH  = 1280
	HEIGHT = 720
)

type Window struct {
	stateManager statesmanager.StateManager
	width        int
	height       int
	//server PID
	sharedCfg utils.SharedConfig
	running   bool

	messageServiceError bool
	//and other errors this is good idea
	loginScene        loginstate.LoginScene
	hcMenuScene       hcmenustate.HCMenuScene
	createUnitScene   createunitstate.CreateUnitScene
	infoUnitScene     infounitstate.InfoUnitScene
	createUserScene   createuserstate.CreateUserScene
	infoUserScene     infouserstate.InfoUserScene
	inboxScene        inboxstate.InboxScene
	createDeviceScene createdevicestate.CreateDeviceScene
	createTaskScene   createtaskstate.CreateTaskScene
	Flow              chan statesmanager.GameState
	Done              chan bool
}

func NewWindowActor(w *Window, serverPID, messageServicePID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		w.sharedCfg.ServerPID = serverPID
		w.sharedCfg.MessageServicePID = messageServicePID
		return w
	}
}

func NewWindow() *Window {
	return &Window{
		Done: make(chan bool, 1),
		Flow: make(chan statesmanager.GameState, 1024),
	}
}

func (w *Window) Receive(ctx *actor.Context) {
	w.sharedCfg.Ctx = ctx
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Actor initialized")
	case actor.Started:
		utils.Logger.Info("Window actor started")
	case actor.Stopped:
		utils.Logger.Info("Actor stopped")
		close(w.inboxScene.MessageSection.MessageChan)
	case *proto.Ping:
		ctx.Respond(&proto.Pong{})
	case *proto.DeliverMessage:
		fmt.Println("ODEBRALEM,", msg.Message)
		w.inboxScene.MessageSection.MessageChan <- msg.Message
	default:
		utils.Logger.Warn("server got unknown errorMessage", reflect.TypeOf(msg).String())
	}
}
func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	w.stateManager = statesmanager.StateManager{
		Flow:         w.Flow,
		SceneStack:   make([]statesmanager.GameState, 0, 64),
		CurrentState: -1,
	}
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	w.width = WIDTH
	w.height = HEIGHT
	rl.InitWindow(int32(w.width), int32(w.height), "BCS application")
	rl.SetTargetFPS(60)

	go w.updateFlow()
	w.stateManager.Add(statesmanager.LoginState)

	/*w.stateManager.SceneStack = append(w.stateManager.SceneStack, statesmanager.LoginState)
	//init login state
	w.stateManager.CurrentState = statesmanager.LoginState
	w.loginScene.LoginSceneSetup(&w)
	*/
	w.running = true

}
func (w *Window) quit() {
	//maybe to quit some assets
	close(w.Flow)
}

func (w *Window) updateFlow() {
	for state := range w.Flow {
		if state == statesmanager.GoBackState {
			if len(w.stateManager.SceneStack) > 1 {
				w.stateManager.SceneStack = w.stateManager.SceneStack[:len(w.stateManager.SceneStack)-1]
				w.stateManager.CurrentState = w.stateManager.SceneStack[len(w.stateManager.SceneStack)-1]
				w.setupSceneForState(w.stateManager.CurrentState)
			}
			continue
		}

		w.stateManager.CurrentState = state
		w.stateManager.SceneStack = append(w.stateManager.SceneStack, state)
		w.setupSceneForState(state)

	}
}

func (w *Window) setupSceneForState(state statesmanager.GameState) {
	switch state {
	case statesmanager.LoginState:
		w.loginScene.LoginSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.HCMenuState:
		w.hcMenuScene.MenuHCSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.CreateUnitState:
		w.createUnitScene.CreateUnitSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.InfoUnitState:
		w.infoUnitScene.InfoUnitSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.CreateUserState:
		w.createUserScene.CreateUserSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.CreateDeviceState:
		w.createDeviceScene.CreateDeviceSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.InfoUserState:
		w.infoUserScene.InfoUserSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.CreateTaskState:
		w.createTaskScene.CreateTaskSceneSetup(&w.stateManager, &w.sharedCfg)
	case statesmanager.InboxState:
		w.inboxScene.SetupInboxScene(&w.stateManager, &w.sharedCfg)
	}
}
func (w *Window) update() {
	w.running = !rl.WindowShouldClose()
	switch w.stateManager.CurrentState {
	case statesmanager.LoginState:
		w.loginScene.UpdateLoginState()
	case statesmanager.HCMenuState:
		w.hcMenuScene.UpdateHCMenuState()
	case statesmanager.CreateUnitState:
		w.createUnitScene.UpdateCreateUnitState()
	case statesmanager.InfoUnitState:
		w.infoUnitScene.UpdateInfoUnitState()
	case statesmanager.CreateUserState:
		w.createUserScene.UpdateCreateUserState()
	case statesmanager.CreateDeviceState:
		w.createDeviceScene.UpdateCreateDeviceState()
	case statesmanager.InfoUserState:
		w.infoUserScene.UpdateInfoUserState()
	case statesmanager.CreateTaskState:
		w.createTaskScene.UpdateCreateTaskState()
	case statesmanager.InboxState:
		w.inboxScene.UpdateInboxState()
	}

}
func (w *Window) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	switch w.stateManager.CurrentState {
	case statesmanager.LoginState:
		w.loginScene.RenderLoginState()
	case statesmanager.HCMenuState:
		w.hcMenuScene.RenderHCMenuState()
	case statesmanager.CreateUnitState:
		w.createUnitScene.RenderCreateUnitState()
	case statesmanager.InfoUnitState:
		w.infoUnitScene.RenderInfoUnitState()
	case statesmanager.CreateUserState:
		w.createUserScene.RenderCreateUserState()
	case statesmanager.InfoUserState:
		w.infoUserScene.RenderInfoUserState()
	case statesmanager.CreateDeviceState:
		w.createDeviceScene.RenderCreateDeviceState()
	case statesmanager.CreateTaskState:
		w.createTaskScene.RenderCreateTaskState()
	case statesmanager.InboxState:
		w.inboxScene.RenderInboxState()
	default:
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
		w.update()
		w.render()
	}
	w.quit()
}
