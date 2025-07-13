package application

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"runtime"
)

type Window struct {
	state      state
	components components
	running    bool

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

type components struct{}
type state struct {
	loginScene bool
	//menuScene bool
	//...Scene bool
	//...Scene bool
}

func init() {
	runtime.LockOSThread()
}
func (w *Window) setup() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.SetTargetFPS(60)

	w.running = true

}
func (w *Window) quit() {
	//maybe to quit some assets

}

func (w *Window) drawScene() {}
func (w *Window) input() {

}

var cmsg string

func (w *Window) update() {
	w.running = !rl.WindowShouldClose()
	select {
	case msg := <-w.Test:
		cmsg = string(msg.Data)
		fmt.Println("Received message:", msg)
	default:
		fmt.Println("xd")
		cmsg = "nic nie ma"
	}

}
func (w *Window) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText(cmsg, 190, 200, 20, rl.Red)

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
