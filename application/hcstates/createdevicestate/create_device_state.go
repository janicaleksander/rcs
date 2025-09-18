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
	backButton       component.Button
	newDeviceSection NewDeviceSection
	errorSection     ErrorSection
	infoSection      InfoSection
}

type NewDeviceSection struct {
	users        []*proto.User
	deviceTypes  []int
	nameInput    component.InputBox
	ownerSlider  component.ListSlider
	typesToggle  component.ToggleGroup
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
	d.newDeviceSection.nameInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()/2-120),
			float32(rl.GetScreenHeight()/2)-100,
			240, 40),
	)
	d.newDeviceSection.ownerSlider.IdxScroll = 0
	d.newDeviceSection.ownerSlider.IdxActiveElement = 0
	d.newDeviceSection.ownerSlider.Bounds = rl.NewRectangle(
		d.newDeviceSection.nameInput.Bounds.X,
		d.newDeviceSection.nameInput.Bounds.Y+d.newDeviceSection.nameInput.Bounds.Height+10,
		240, 50,
	)

	d.newDeviceSection.typesToggle.Labels = []string{"Mobile (0)"}
	d.newDeviceSection.typesToggle.Bounds = []rl.Rectangle{rl.NewRectangle(
		float32(rl.GetScreenWidth()/2-120),
		d.newDeviceSection.ownerSlider.Bounds.Y+d.newDeviceSection.ownerSlider.Bounds.Height+20,
		100, 40)}

	d.newDeviceSection.acceptButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			d.newDeviceSection.nameInput.Bounds.X,
			d.newDeviceSection.typesToggle.Bounds[0].Y+50,
			240, 50), "CREATE", false)

	d.backButton = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()-110),
			float32(rl.GetScreenHeight()-50),
			100,
			40,
		), "GO BACK", false)

	d.errorSection.errorPopup = *component.NewPopup(component.NewPopupConfig(component.WithBgColor(utils.POPUPERRORBG)), rl.NewRectangle(
		d.newDeviceSection.acceptButton.Bounds.X,
		d.newDeviceSection.acceptButton.Bounds.Y+d.newDeviceSection.acceptButton.Bounds.Height+10,
		d.newDeviceSection.acceptButton.Bounds.Width,
		100), &d.errorSection.error)

	d.infoSection.infoPopup = *component.NewPopup(component.NewPopupConfig(component.WithBgColor(utils.POPUPINFOBG)), rl.NewRectangle(
		d.newDeviceSection.acceptButton.Bounds.X,
		d.newDeviceSection.acceptButton.Bounds.Y+d.newDeviceSection.acceptButton.Bounds.Height+10,
		d.newDeviceSection.acceptButton.Bounds.Width,
		60), &d.infoSection.info)

	d.FetchUsers()
	d.FetchDeviceTypes()
}
func (d *CreateDeviceScene) UpdateCreateDeviceState() {
	d.newDeviceSection.nameInput.Update()
	for i := range d.newDeviceSection.typesToggle.Labels {
		toggleState := d.newDeviceSection.typesToggle.Selected == i
		if gui.Toggle(
			d.newDeviceSection.typesToggle.Bounds[i],
			d.newDeviceSection.typesToggle.Labels[i],
			toggleState,
		) {
			d.newDeviceSection.typesToggle.Selected = i
		}
		if d.newDeviceSection.acceptButton.Update() {
			d.CreateDevice()
		}
	}
	if d.backButton.Update() {
		d.state.Add(statesmanager.GoBackState)
		return
	}

}
func (d *CreateDeviceScene) RenderCreateDeviceState() {
	rl.ClearBackground(utils.CREATEDEVICEBG)
	d.backButton.Render()
	d.newDeviceSection.acceptButton.Render()
	d.newDeviceSection.nameInput.Render()

	rl.DrawText(
		"CREATE DEVICE",
		int32(rl.GetScreenWidth()/2)-int32(rl.MeasureText("CREATE DEVICE", 50)/2),
		50,
		50,
		rl.DarkGray,
	)
	rl.DrawText(
		"NAME",
		int32(d.newDeviceSection.nameInput.Bounds.X)-180,
		int32(d.newDeviceSection.nameInput.Bounds.Y)+10,
		25,
		rl.DarkGray,
	)
	rl.DrawText(
		"OWNER",
		int32(d.newDeviceSection.ownerSlider.Bounds.X)-180,
		int32(d.newDeviceSection.ownerSlider.Bounds.Y)+10,
		25,
		rl.DarkGray,
	)
	rl.DrawText(
		"TYPE",
		int32(d.newDeviceSection.typesToggle.Bounds[0].X)-180,
		int32(d.newDeviceSection.typesToggle.Bounds[0].Y)+10,
		25,
		rl.DarkGray,
	)
	gui.ListViewEx(
		d.newDeviceSection.ownerSlider.Bounds,
		d.newDeviceSection.ownerSlider.Strings,
		&d.newDeviceSection.ownerSlider.IdxScroll,
		&d.newDeviceSection.ownerSlider.IdxActiveElement,
		d.newDeviceSection.ownerSlider.Focus)

	d.errorSection.errorPopup.Render()
	d.infoSection.infoPopup.Render()
}
