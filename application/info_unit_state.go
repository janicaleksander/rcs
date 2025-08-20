package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
)

type InfoUnitScene struct {
	backButton               component.Button
	unitsSection             UnitsSection
	usersSection             UsersSection
	descriptionButtonSection DescriptionButtonSection
}

type UnitsSection struct {
	units                []*proto.Unit
	unitsToUserCache     map[string][]*proto.User //preload users in each unit from units
	unitsSlider          ListSlider
	lastProcessedUnitIdx int32
}

type UsersSection struct {

	//last part with commander squad of selected unit (5,4,3 rank)
	usersSlider ListSlider
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
func (w *Window) infoUnitSceneSetup() {
	w.infoUnitScene.Reset()
	w.FetchUnits2()
	w.infoUnitScene.unitsSection.unitsSlider = ListSlider{
		strings: make([]string, 0),
		bounds: rl.NewRectangle(
			0,
			(1.0/8.0)*float32(w.height),
			(2.0/9.0)*float32(w.width),
			(7.0/8.0)*float32(w.height),
		),
		idxActiveElement: 0, // ?
		focus:            1,
		idxScroll:        0,
	}

	for _, v := range w.infoUnitScene.unitsSection.units {
		w.infoUnitScene.unitsSection.unitsSlider.strings = append(w.infoUnitScene.unitsSection.unitsSlider.strings, v.Id[:5]+"..."+v.Id[31:])
	}

	w.infoUnitScene.descriptionButtonSection.descriptionBounds = rl.NewRectangle(
		w.infoUnitScene.unitsSection.unitsSlider.bounds.Width,
		w.infoUnitScene.unitsSection.unitsSlider.bounds.Y,
		(4.0/9.0)*float32(w.width),
		w.infoUnitScene.unitsSection.unitsSlider.bounds.Height)

	w.infoUnitScene.usersSection.usersSlider = ListSlider{
		strings: make([]string, 0),
		bounds: rl.NewRectangle(
			w.infoUnitScene.unitsSection.unitsSlider.bounds.Width+w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width,
			w.infoUnitScene.unitsSection.unitsSlider.bounds.Y,
			(4.0/12.0)*float32(w.width),
			(2.0/3.0)*float32(w.height)),
		idxActiveElement: 0, // ?
		focus:            1,
		idxScroll:        0,
	}
	w.infoUnitScene.usersSection.userInfoBounds = rl.NewRectangle(
		w.infoUnitScene.unitsSection.unitsSlider.bounds.Width+w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width,
		(2.0/3.0)*float32(w.height),
		w.infoUnitScene.usersSection.usersSlider.bounds.Width,
		(1.0/3.0)*float32(w.height))

	//TODO do full refactor of this conditions and check where have to have it
	if len(w.infoUnitScene.unitsSection.units) > 0 {
		w.infoUnitScene.unitsSection.unitsSlider.idxActiveElement = 0
	} else {
		w.infoUnitScene.usersSection.usersSlider.idxActiveElement = -1
	}

	w.infoUnitScene.descriptionButtonSection.descriptionButton1 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.X,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Y,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width/2,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 1", false)

	w.infoUnitScene.descriptionButtonSection.descriptionButton2 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.X+w.infoUnitScene.descriptionButtonSection.descriptionButton1.Bounds.Width,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Y,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width/2,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 2", false)

	w.infoUnitScene.descriptionButtonSection.descriptionButton3 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.X,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Y+w.infoUnitScene.descriptionButtonSection.descriptionButton1.Bounds.Height,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width/2,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 3", false)

	w.infoUnitScene.descriptionButtonSection.descriptionButton4 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.X+w.infoUnitScene.descriptionButtonSection.descriptionButton1.Bounds.Width,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Y+w.infoUnitScene.descriptionButtonSection.descriptionButton1.Bounds.Height,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width/2,
		w.infoUnitScene.descriptionButtonSection.descriptionBounds.Height/2,
	), "Button info 4", false)

	w.infoUnitScene.usersSection.userInfoButton1 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUnitScene.usersSection.userInfoBounds.X,
		w.infoUnitScene.usersSection.userInfoBounds.Y,
		w.infoUnitScene.usersSection.userInfoBounds.Width,
		w.infoUnitScene.usersSection.userInfoBounds.Height/2,
	), "User info 1", false)

	w.infoUnitScene.usersSection.userInfoButton2 = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUnitScene.usersSection.userInfoBounds.X,
		w.infoUnitScene.usersSection.userInfoBounds.Y+w.infoUnitScene.usersSection.userInfoButton1.Bounds.Height,
		w.infoUnitScene.usersSection.userInfoBounds.Width,
		w.infoUnitScene.usersSection.userInfoBounds.Height/2,
	), "User info 2", false)

	w.infoUnitScene.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			float32(w.width-100),
			float32(w.height-50),
			100,
			50),
		"Go back",
		false,
	)
}

//TODO add some info where e.g. users are empty

func (w *Window) updateInfoUnitState() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.infoUnitScene.unitsSection.unitsSlider.bounds) {
			w.infoUnitScene.unitsSection.unitsSlider.focus = 1
			w.infoUnitScene.usersSection.usersSlider.focus = 0
		}
		if rl.CheckCollisionPointRec(mousePos, w.infoUnitScene.usersSection.usersSlider.bounds) {
			w.infoUnitScene.unitsSection.unitsSlider.focus = 0
			w.infoUnitScene.usersSection.usersSlider.focus = 1
		}

	}
	w.infoUnitScene.usersSection.userInfoButton1.Update() //info
	w.infoUnitScene.usersSection.userInfoButton2.Update() //send message
	w.infoUnitScene.descriptionButtonSection.descriptionButton1.Update()
	w.infoUnitScene.descriptionButtonSection.descriptionButton2.Update()
	w.infoUnitScene.descriptionButtonSection.descriptionButton3.Update()
	w.infoUnitScene.descriptionButtonSection.descriptionButton4.Update()

	if w.infoUnitScene.backButton.Update() {
		w.goSceneBack()
	}

	w.SelectUnit()
	w.SelectUser()
}

func (w *Window) renderInfoUnitState() {

	//unitsSlider
	gui.ListViewEx(
		w.infoUnitScene.unitsSection.unitsSlider.bounds,
		w.infoUnitScene.unitsSection.unitsSlider.strings,
		&w.infoUnitScene.unitsSection.unitsSlider.idxScroll,
		&w.infoUnitScene.unitsSection.unitsSlider.idxActiveElement,
		w.infoUnitScene.unitsSection.unitsSlider.focus,
	)
	//description box
	rl.DrawRectangle(int32(w.infoUnitScene.descriptionButtonSection.descriptionBounds.X),
		int32(w.infoUnitScene.descriptionButtonSection.descriptionBounds.Y),
		int32(w.infoUnitScene.descriptionButtonSection.descriptionBounds.Width),
		int32(w.infoUnitScene.descriptionButtonSection.descriptionBounds.Height),
		rl.Yellow)
	//users slider
	gui.ListViewEx(w.infoUnitScene.usersSection.usersSlider.bounds,
		w.infoUnitScene.usersSection.usersSlider.strings,
		&w.infoUnitScene.usersSection.usersSlider.idxScroll,
		&w.infoUnitScene.usersSection.usersSlider.idxActiveElement,
		w.infoUnitScene.usersSection.usersSlider.focus)

	//user info box
	rl.DrawRectangle(int32(w.infoUnitScene.usersSection.userInfoBounds.X),
		int32(w.infoUnitScene.usersSection.userInfoBounds.Y),
		int32(w.infoUnitScene.usersSection.userInfoBounds.Width),
		int32(w.infoUnitScene.usersSection.userInfoBounds.Height),
		rl.White)

	w.infoUnitScene.usersSection.userInfoButton1.Render() //info
	w.infoUnitScene.usersSection.userInfoButton2.Render() //send message
	w.infoUnitScene.descriptionButtonSection.descriptionButton1.Render()
	w.infoUnitScene.descriptionButtonSection.descriptionButton2.Render()
	w.infoUnitScene.descriptionButtonSection.descriptionButton3.Render()
	w.infoUnitScene.descriptionButtonSection.descriptionButton4.Render()
	w.infoUnitScene.backButton.Render()
}
