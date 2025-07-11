package application

import rl "github.com/gen2brain/raylib-go/raylib"

type Window struct {
	state      state
	components components
	running    bool

	Done chan bool
}

func NewWindow() *Window {
	return &Window{
		Done: make(chan bool),
	}
}

type components struct{}
type state struct {
	loginScene bool
	//menuScene bool
	//...Scene bool
	//...Scene bool
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
func (w *Window) update() {
	w.running = !rl.WindowShouldClose()
}
func (w *Window) render() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

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
