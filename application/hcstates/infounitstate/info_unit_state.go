package infounitstate

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

type InfoUnitScene struct {
	cfg                *utils.SharedConfig
	stateManager       *statesmanager.StateManager
	backButton         rl.Rectangle
	unitsSection       UnitsSection
	usersSection       UsersSection
	descriptionSection DescriptionSection
	errorSection       ErrorSection
}

type UnitsSection struct {
	units                []*proto.Unit
	unitsInformation     map[string]*proto.UnitInformation
	unitsSlider          component.ListSlider
	currUnitID           string
	lastProcessedUnitIdx int32
}

type UsersSection struct {

	//last part with commander squad of selected unit (5,4,3 rank)
	usersSlider component.ListSlider
	//section with info about selected user
	userInfoBounds       rl.Rectangle
	userInfoButton1      rl.Rectangle //info
	userInfoButton2      rl.Rectangle //send message
	currUserID           string
	lastProcessedUserIdx int32
}

type DescriptionSection struct {
	descriptionBounds rl.Rectangle
	descButtonDevices rl.Rectangle
	descButtonTasks   rl.Rectangle
	devicesModal      component.Modal
	devicesSlider     component.ScrollPanel
	devicesElements   map[string][]struct {
		bounds rl.Rectangle
		name   string
		desc   string
	}
	showDevicesModal bool
	tasksModal       component.Modal
	tasksSlider      component.ScrollPanel
	tasksElements    map[string][]struct {
		bounds rl.Rectangle
		name   string
		desc   string
	}
	showTasksModal bool
}
type ErrorSection struct {
	message  string
	errPopup component.Popup
}

// TODO if slice is empty show some info about this (Nothing is here maybe)+maybe err field if error
func (i *InfoUnitScene) InfoUnitSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	i.cfg = cfg
	i.stateManager = state
	i.Reset()
	i.FetchUnits()
	i.UnitsDescription()
	i.unitsSection.unitsSlider = component.ListSlider{
		Strings: make([]string, 0),
		Bounds: rl.NewRectangle(
			0,
			(1.0/8.0)*float32(rl.GetScreenHeight()),
			(2.0/9.0)*float32(rl.GetScreenWidth()),
			(7.0/8.0)*float32(rl.GetScreenHeight()),
		),
		IdxActiveElement: 0, // ?
		Focus:            1,
		IdxScroll:        0,
	}

	for _, v := range i.unitsSection.units {
		i.unitsSection.unitsSlider.Strings = append(i.unitsSection.unitsSlider.Strings, v.Id[:5]+"..."+v.Id[31:])
	}

	i.descriptionSection.descriptionBounds = rl.NewRectangle(
		i.unitsSection.unitsSlider.Bounds.Width,
		i.unitsSection.unitsSlider.Bounds.Y,
		(4.0/9.0)*float32(rl.GetScreenWidth()),
		i.unitsSection.unitsSlider.Bounds.Height)

	i.descriptionSection.descButtonDevices = rl.NewRectangle(
		i.descriptionSection.descriptionBounds.X,
		i.descriptionSection.descriptionBounds.Y,
		i.descriptionSection.descriptionBounds.Width/2,
		i.descriptionSection.descriptionBounds.Height/2,
	)
	i.descriptionSection.descButtonTasks = rl.NewRectangle(
		i.descriptionSection.descButtonDevices.X+i.descriptionSection.descButtonDevices.Width,
		i.descriptionSection.descButtonDevices.Y,
		i.descriptionSection.descriptionBounds.Width/2,
		i.descriptionSection.descriptionBounds.Height/2,
	)
	//i.descriptionButtonSection.descButtonTasks = rl.NewRectangle()

	i.usersSection.usersSlider = component.ListSlider{
		Strings: make([]string, 0),
		Bounds: rl.NewRectangle(
			i.unitsSection.unitsSlider.Bounds.Width+i.descriptionSection.descriptionBounds.Width,
			i.unitsSection.unitsSlider.Bounds.Y,
			(4.0/12.0)*float32(rl.GetScreenWidth()),
			(2.0/3.0)*float32(rl.GetScreenHeight())),
		IdxActiveElement: -1, // ?
		Focus:            1,
		IdxScroll:        0,
	}
	i.usersSection.userInfoBounds = rl.NewRectangle(
		i.unitsSection.unitsSlider.Bounds.Width+i.descriptionSection.descriptionBounds.Width,
		(2.0/3.0)*float32(rl.GetScreenHeight()),
		i.usersSection.usersSlider.Bounds.Width,
		(1.0/3.0)*float32(rl.GetScreenHeight()))

	i.usersSection.userInfoButton1 = rl.NewRectangle(
		i.usersSection.userInfoBounds.X,
		i.usersSection.userInfoBounds.Y,
		i.usersSection.userInfoBounds.Width,
		50,
	)

	i.usersSection.userInfoButton2 = rl.NewRectangle(
		i.usersSection.userInfoBounds.X,
		i.usersSection.userInfoBounds.Y+i.usersSection.userInfoButton1.Height,
		i.usersSection.userInfoBounds.Width,
		50,
	)

	modalWindowBg := rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
	modalWindowBgColor := rl.Fade(rl.Gray, 0.3)
	modalCore := rl.NewRectangle(
		(1.0/3.0)*float32(rl.GetScreenWidth())-180,
		(1.0/3.0)*float32(rl.GetScreenHeight())-150,
		800,
		500,
	)
	i.descriptionSection.devicesModal = component.Modal{
		Background: modalWindowBg,
		BgColor:    modalWindowBgColor,
		Core:       modalCore,
	}
	//owner + lastonline
	i.descriptionSection.devicesSlider = component.ScrollPanel{
		Bounds: rl.NewRectangle(
			i.descriptionSection.devicesModal.Core.X,
			i.descriptionSection.devicesModal.Core.Y+25,
			i.descriptionSection.devicesModal.Core.Width,
			i.descriptionSection.devicesModal.Core.Height-25,
		),
		Content: rl.NewRectangle(
			0,
			0,
			i.descriptionSection.devicesModal.Core.Width-30,
			i.descriptionSection.devicesModal.Core.Height*10,
		),
		Scroll: rl.Vector2{X: 0, Y: 0},
		View:   rl.Rectangle{},
	}

	i.descriptionSection.tasksModal = component.Modal{
		Background: modalWindowBg,
		BgColor:    modalWindowBgColor,
		Core:       modalCore,
	}
	//owner + lastonline
	i.descriptionSection.tasksSlider = component.ScrollPanel{
		Bounds: rl.NewRectangle(
			i.descriptionSection.devicesModal.Core.X,
			i.descriptionSection.devicesModal.Core.Y+25,
			i.descriptionSection.devicesModal.Core.Width,
			i.descriptionSection.devicesModal.Core.Height-25,
		),
		Content: rl.NewRectangle(
			0,
			0,
			i.descriptionSection.devicesModal.Core.Width-30,
			i.descriptionSection.devicesModal.Core.Height*10,
		),
		Scroll: rl.Vector2{X: 0, Y: 0},
		View:   rl.Rectangle{},
	}

	i.prepareDeviceSlider()
	i.prepareTaskSlider()
	i.backButton = rl.NewRectangle(
		float32(rl.GetScreenWidth()-100),
		float32(20),
		100,
		50)
}

//TODO add some info where e.g. users are empty

func (i *InfoUnitScene) UpdateInfoUnitState() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, i.unitsSection.unitsSlider.Bounds) {
			i.unitsSection.unitsSlider.Focus = 1
			i.usersSection.usersSlider.Focus = 0
		}
		if rl.CheckCollisionPointRec(mousePos, i.usersSection.usersSlider.Bounds) {
			i.unitsSection.unitsSlider.Focus = 0
			i.usersSection.usersSlider.Focus = 1
		}

	}

	i.SelectUnit()
	i.SelectUser()
}

func (i *InfoUnitScene) RenderInfoUnitState() {

	//unitsSlider
	gui.ListViewEx(
		i.unitsSection.unitsSlider.Bounds,
		i.unitsSection.unitsSlider.Strings,
		&i.unitsSection.unitsSlider.IdxScroll,
		&i.unitsSection.unitsSlider.IdxActiveElement,
		i.unitsSection.unitsSlider.Focus,
	)
	//description box
	rl.DrawRectangleLinesEx(rl.NewRectangle(i.descriptionSection.descriptionBounds.X,
		i.descriptionSection.descriptionBounds.Y,
		i.descriptionSection.descriptionBounds.Width,
		i.descriptionSection.descriptionBounds.Height),
		1.5,
		rl.Black)
	//users slider
	gui.ListViewEx(i.usersSection.usersSlider.Bounds,
		i.usersSection.usersSlider.Strings,
		&i.usersSection.usersSlider.IdxScroll,
		&i.usersSection.usersSlider.IdxActiveElement,
		i.usersSection.usersSlider.Focus)
	if gui.Button(rl.NewRectangle(i.backButton.X, i.backButton.Y, i.backButton.Width, i.backButton.Height), "GO BACK") {
		i.stateManager.Add(statesmanager.GoBackState)
		return
	}
	if gui.Button(i.descriptionSection.descButtonDevices, "DEVICES") {
		i.descriptionSection.showDevicesModal = true
	}
	if gui.Button(i.descriptionSection.descButtonTasks, "TASKS") {
		i.descriptionSection.showTasksModal = true

	}
	if gui.Button(i.usersSection.userInfoButton1, "INFO 1") {

	}
	if gui.Button(i.usersSection.userInfoButton2, "INFO 2") {

	}
	if i.descriptionSection.showDevicesModal {
		rl.DrawRectangle(
			int32(i.descriptionSection.devicesModal.Background.X),
			int32(i.descriptionSection.devicesModal.Background.Y),
			int32(i.descriptionSection.devicesModal.Background.Width),
			int32(i.descriptionSection.devicesModal.Background.Height),
			i.descriptionSection.devicesModal.BgColor,
		)

		if gui.WindowBox(i.descriptionSection.devicesModal.Core, "See devices") {
			i.descriptionSection.showDevicesModal = false
		}

		gui.ScrollPanel(
			i.descriptionSection.devicesSlider.Bounds,
			"Devices",
			i.descriptionSection.devicesSlider.Content,
			&i.descriptionSection.devicesSlider.Scroll,
			&i.descriptionSection.devicesSlider.View,
		)

		view := i.descriptionSection.devicesSlider.View
		rl.BeginScissorMode(int32(view.X), int32(view.Y), int32(view.Width), int32(view.Height))
		//TODO curr unitID
		currMap := i.descriptionSection.devicesElements[i.unitsSection.currUnitID]
		padding := int32(10)
		rectHeight := int32(80)
		textSpacing := int32(5)
		descFontSize := int32(14)
		nameFontSize := int32(20)
		sliderX := int32(i.descriptionSection.devicesSlider.Bounds.X)
		sliderY := int32(i.descriptionSection.devicesSlider.Bounds.Y)

		for _, elem := range currMap {
			y := int32(25 + int(elem.bounds.Y) + int(i.descriptionSection.devicesSlider.Scroll.Y))

			rl.DrawRectangleRounded(
				rl.Rectangle{
					X:      float32(sliderX),
					Y:      float32(sliderY + y),
					Width:  float32(i.descriptionSection.devicesSlider.Bounds.Width),
					Height: float32(rectHeight),
				},
				0.2,
				10,
				rl.LightGray,
			)

			rl.DrawRectangleLinesEx(
				rl.Rectangle{
					X:      float32(sliderX),
					Y:      float32(sliderY + y),
					Width:  float32(i.descriptionSection.devicesSlider.Bounds.Width),
					Height: float32(rectHeight),
				},
				2, rl.Gray,
			)

			rl.DrawText(elem.name, sliderX+padding, sliderY+y+padding, nameFontSize, rl.Black)
			rl.DrawText(elem.desc, sliderX+padding, sliderY+y+padding+nameFontSize+textSpacing, descFontSize, rl.DarkGray)
		}

		rl.EndScissorMode()
	}
	//TODO below
	if i.descriptionSection.showTasksModal {
		rl.DrawRectangle(
			int32(i.descriptionSection.devicesModal.Background.X),
			int32(i.descriptionSection.devicesModal.Background.Y),
			int32(i.descriptionSection.devicesModal.Background.Width),
			int32(i.descriptionSection.devicesModal.Background.Height),
			i.descriptionSection.devicesModal.BgColor,
		)

		if gui.WindowBox(i.descriptionSection.tasksModal.Core, "See tasks") {
			i.descriptionSection.showTasksModal = false
		}

		gui.ScrollPanel(
			i.descriptionSection.tasksSlider.Bounds,
			"Tasks",
			i.descriptionSection.tasksSlider.Content,
			&i.descriptionSection.tasksSlider.Scroll,
			&i.descriptionSection.tasksSlider.View,
		)

		view := i.descriptionSection.tasksSlider.View
		rl.BeginScissorMode(int32(view.X), int32(view.Y), int32(view.Width), int32(view.Height))
		currMap := i.descriptionSection.tasksElements[i.unitsSection.currUnitID]
		padding := int32(10)
		rectHeight := int32(80)
		textSpacing := int32(5)
		descFontSize := int32(14)
		nameFontSize := int32(20)
		sliderX := int32(i.descriptionSection.tasksSlider.Bounds.X)
		sliderY := int32(i.descriptionSection.tasksSlider.Bounds.Y)

		for _, elem := range currMap {
			y := int32(25 + int(elem.bounds.Y) + int(i.descriptionSection.tasksSlider.Scroll.Y))

			rl.DrawRectangleRounded(
				rl.Rectangle{
					X:      float32(sliderX),
					Y:      float32(sliderY + y),
					Width:  float32(i.descriptionSection.tasksSlider.Bounds.Width),
					Height: float32(rectHeight),
				},
				0.2,
				10,
				rl.LightGray,
			)

			rl.DrawRectangleLinesEx(
				rl.Rectangle{
					X:      float32(sliderX),
					Y:      float32(sliderY + y),
					Width:  float32(i.descriptionSection.tasksSlider.Bounds.Width),
					Height: float32(rectHeight),
				},
				2, rl.Gray,
			)

			rl.DrawText(elem.name, sliderX+padding, sliderY+y+padding, nameFontSize, rl.Black)
			rl.DrawText(elem.desc, sliderX+padding, sliderY+y+padding+nameFontSize+textSpacing, descFontSize, rl.DarkGray)
		}

		rl.EndScissorMode()
	}

}
