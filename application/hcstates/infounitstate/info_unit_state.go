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
	cfg                      *utils.SharedConfig
	stateManager             *statesmanager.StateManager
	backButton               component.Button
	unitsSection             UnitsSection
	usersSection             UsersSection
	descriptionButtonSection DescriptionButtonSection
}

type UnitsSection struct {
	units                []*proto.Unit
	unitsToUserCache     map[string][]*proto.User //preload users in each unit from units
	unitsSlider          component.ListSlider
	lastProcessedUnitIdx int32
}

type UsersSection struct {

	//last part with commander squad of selected unit (5,4,3 rank)
	usersSlider component.ListSlider
	//section with info about selected user
	userInfoBounds       rl.Rectangle
	userInfoButton1      component.Button //info
	userInfoButton2      component.Button //send message
	lastProcessedUserIdx int32
}

type DescriptionButtonSection struct {

	//middle with general info about selected
	descriptionBounds  rl.Rectangle
	descriptionButton1 component.Button
	descriptionButton2 component.Button
	descriptionButton3 component.Button
	descriptionButton4 component.Button
	//... some elements about unit

}

// TODO if slice is empty show some info about this (Nothing is here maybe)+maybe err field if error
func (i *InfoUnitScene) InfoUnitSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	i.cfg = cfg
	i.stateManager = state
	i.Reset()
	i.FetchUnits()
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

	i.descriptionButtonSection.descriptionBounds = rl.NewRectangle(
		i.unitsSection.unitsSlider.Bounds.Width,
		i.unitsSection.unitsSlider.Bounds.Y,
		(4.0/9.0)*float32(rl.GetScreenWidth()),
		i.unitsSection.unitsSlider.Bounds.Height)

	i.usersSection.usersSlider = component.ListSlider{
		Strings: make([]string, 0),
		Bounds: rl.NewRectangle(
			i.unitsSection.unitsSlider.Bounds.Width+i.descriptionButtonSection.descriptionBounds.Width,
			i.unitsSection.unitsSlider.Bounds.Y,
			(4.0/12.0)*float32(rl.GetScreenWidth()),
			(2.0/3.0)*float32(rl.GetScreenHeight())),
		IdxActiveElement: 0, // ?
		Focus:            1,
		IdxScroll:        0,
	}
	i.usersSection.userInfoBounds = rl.NewRectangle(
		i.unitsSection.unitsSlider.Bounds.Width+i.descriptionButtonSection.descriptionBounds.Width,
		(2.0/3.0)*float32(rl.GetScreenHeight()),
		i.usersSection.usersSlider.Bounds.Width,
		(1.0/3.0)*float32(rl.GetScreenHeight()))

	//TODO do full refactor of this conditions and check where have to have it
	if len(i.unitsSection.units) > 0 {
		i.unitsSection.unitsSlider.IdxActiveElement = 0
	} else {
		i.usersSection.usersSlider.IdxActiveElement = -1
	}

	i.descriptionButtonSection.descriptionButton1 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.descriptionButtonSection.descriptionBounds.X,
		i.descriptionButtonSection.descriptionBounds.Y,
		i.descriptionButtonSection.descriptionBounds.Width/2,
		i.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 1", false)

	i.descriptionButtonSection.descriptionButton2 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.descriptionButtonSection.descriptionBounds.X+i.descriptionButtonSection.descriptionButton1.Bounds.Width,
		i.descriptionButtonSection.descriptionBounds.Y,
		i.descriptionButtonSection.descriptionBounds.Width/2,
		i.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 2", false)

	i.descriptionButtonSection.descriptionButton3 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.descriptionButtonSection.descriptionBounds.X,
		i.descriptionButtonSection.descriptionBounds.Y+i.descriptionButtonSection.descriptionButton1.Bounds.Height,
		i.descriptionButtonSection.descriptionBounds.Width/2,
		i.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 3", false)

	i.descriptionButtonSection.descriptionButton4 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.descriptionButtonSection.descriptionBounds.X+i.descriptionButtonSection.descriptionButton1.Bounds.Width,
		i.descriptionButtonSection.descriptionBounds.Y+i.descriptionButtonSection.descriptionButton1.Bounds.Height,
		i.descriptionButtonSection.descriptionBounds.Width/2,
		i.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 4", false)

	i.usersSection.userInfoButton1 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.usersSection.userInfoBounds.X,
		i.usersSection.userInfoBounds.Y,
		i.usersSection.userInfoBounds.Width,
		i.usersSection.userInfoBounds.Height/2,
	), "User info 1", false)

	i.usersSection.userInfoButton2 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.usersSection.userInfoBounds.X,
		i.usersSection.userInfoBounds.Y+i.usersSection.userInfoButton1.Bounds.Height,
		i.usersSection.userInfoBounds.Width,
		i.usersSection.userInfoBounds.Height/2,
	), "User info 2", false)

	i.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()-100),
			float32(rl.GetScreenHeight()-50),
			100,
			50),
		"Go back",
		false,
	)
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
	i.usersSection.userInfoButton1.Update() //info
	i.usersSection.userInfoButton2.Update() //send message
	i.descriptionButtonSection.descriptionButton1.Update()
	i.descriptionButtonSection.descriptionButton2.Update()
	i.descriptionButtonSection.descriptionButton3.Update()
	i.descriptionButtonSection.descriptionButton4.Update()

	if i.backButton.Update() {
		i.stateManager.Add(statesmanager.GoBackState)
		return
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
	rl.DrawRectangle(int32(i.descriptionButtonSection.descriptionBounds.X),
		int32(i.descriptionButtonSection.descriptionBounds.Y),
		int32(i.descriptionButtonSection.descriptionBounds.Width),
		int32(i.descriptionButtonSection.descriptionBounds.Height),
		rl.Yellow)
	//users slider
	gui.ListViewEx(i.usersSection.usersSlider.Bounds,
		i.usersSection.usersSlider.Strings,
		&i.usersSection.usersSlider.IdxScroll,
		&i.usersSection.usersSlider.IdxActiveElement,
		i.usersSection.usersSlider.Focus)

	//user info box
	rl.DrawRectangle(int32(i.usersSection.userInfoBounds.X),
		int32(i.usersSection.userInfoBounds.Y),
		int32(i.usersSection.userInfoBounds.Width),
		int32(i.usersSection.userInfoBounds.Height),
		rl.White)

	i.usersSection.userInfoButton1.Render() //info
	i.usersSection.userInfoButton2.Render() //send message
	i.descriptionButtonSection.descriptionButton1.Render()
	i.descriptionButtonSection.descriptionButton2.Render()
	i.descriptionButtonSection.descriptionButton3.Render()
	i.descriptionButtonSection.descriptionButton4.Render()
	i.backButton.Render()
}
