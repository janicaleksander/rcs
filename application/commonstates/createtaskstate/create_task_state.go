package createtaskstate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

type CreateTaskScene struct {
	cfg          *utils.SharedConfig
	stateManager *statesmanager.StateManager
	errorSection ErrorSection
	infoSection  InfoSection
}

type ErrorSection struct {
	isSetupError  bool
	isCreateError bool
	errorMessage  string
	errorPopup    component.Popup
}
type InfoSection struct {
	isInfoMessage bool
	infoMessage   string
	infoPopup     component.Popup
}

type NewTaskSection struct {
}

func (c *CreateTaskScene) CreateTaskSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	c.cfg = cfg
	c.stateManager = state
}

func (c *CreateTaskScene) UpdateCreateTaskState() {
}
func (c *CreateTaskScene) RenderCreateTaskState() {
	rl.ClearBackground(utils.CREATEUNITBG)
	rl.DrawText("CREATE TASK", int32(rl.GetScreenWidth()/2)-rl.MeasureText("CREATE TASK", 45)/2, 50, 45, rl.DarkGray)

}
