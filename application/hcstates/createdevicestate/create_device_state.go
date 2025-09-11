package createdevicestate

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

type CreateDeviceScene struct {
	cfg              *utils.SharedConfig
	state            *statesmanager.StateManager
	scheduler        utils.Scheduler
	backButton       component.Button
	newDeviceSection NewDeviceSection
	errorSection     ErrorSection
	infoSection      InfoSection
}

type NewDeviceSection struct {
	users        []*proto.User
	deviceTypes  []string
	nameInput    component.InputBox
	ownerSlider  component.ListSlider
	typeSlider   component.ListSlider
	acceptButton component.Button
}
type ErrorSection struct {
	errorPopup component.Popup
	error      string
}
type InfoSection struct {
	infoPopup component.Popup
	info      string
}

func (d *CreateDeviceScene) CreateDeviceSceneSetup(stateManager *statesmanager.StateManager, cfg *utils.SharedConfig) {
	d.cfg = cfg
	d.state = stateManager
	d.Reset()
	d.FetchUsers()
	d.FetchDeviceTypes()

	d.newDeviceSection.nameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-120),
			float32(rl.GetScreenHeight()/2-80),
			240, 40),
	)
	d.newDeviceSection.ownerSlider.IdxScroll = 0
	d.newDeviceSection.ownerSlider.IdxActiveElement = 0
	d.newDeviceSection.ownerSlider.Bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		float32(rl.GetScreenHeight()/2),
		240, 50,
	)

	d.newDeviceSection.typeSlider.IdxScroll = 0
	d.newDeviceSection.typeSlider.IdxActiveElement = 0
	d.newDeviceSection.typeSlider.Bounds = rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		float32(rl.GetScreenHeight()/2+80),
		240, 50,
	)

	d.newDeviceSection.acceptButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-120),
			float32(rl.GetScreenHeight()/2+160),
			240, 80), "Create", false)

	d.backButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(float32(rl.GetScreenWidth()/2-120),
			float32(rl.GetScreenHeight()/2+240), 80, 80), "Go back", false)

	d.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2.0-100.0),
		float32(rl.GetScreenHeight()/2.0+40.0),
		200,
		100), &d.errorSection.error)

	d.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(), rl.NewRectangle(
		float32(rl.GetScreenWidth()/2.0-100.0),
		float32(rl.GetScreenHeight()/2.0+40.0),
		200,
		100), &d.infoSection.info)
}
func (d *CreateDeviceScene) UpdateCreateDeviceState() {
	d.scheduler.Update(float64(rl.GetFrameTime()))
	d.newDeviceSection.nameInput.Update()

	if d.backButton.Update() {
		d.state.Add(statesmanager.GoBackState)
		return
	}
	if d.newDeviceSection.acceptButton.Update() {
		d.CreateDevice()
	}
}
func (d *CreateDeviceScene) RenderCreateDeviceState() {
	d.backButton.Render()
	d.newDeviceSection.acceptButton.Render()
	d.newDeviceSection.nameInput.Render()

	gui.ListViewEx(
		d.newDeviceSection.ownerSlider.Bounds,
		d.newDeviceSection.ownerSlider.Strings,
		&d.newDeviceSection.ownerSlider.IdxScroll,
		&d.newDeviceSection.ownerSlider.IdxActiveElement,
		d.newDeviceSection.ownerSlider.Focus)
	gui.ListViewEx(
		d.newDeviceSection.typeSlider.Bounds,
		d.newDeviceSection.typeSlider.Strings,
		&d.newDeviceSection.typeSlider.IdxScroll,
		&d.newDeviceSection.typeSlider.IdxActiveElement,
		d.newDeviceSection.typeSlider.Focus)
	d.errorSection.errorPopup.Render()
	d.infoSection.infoPopup.Render()
}

// need slider for device types
// need slider for users to assign ownership

//TODO think about logic to login to mobile app and pc app, distinguish this
