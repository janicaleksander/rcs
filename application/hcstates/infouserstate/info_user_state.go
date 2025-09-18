package infouserstate

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

type InfoUserScene struct {
	cfg                      *utils.SharedConfig
	stateManager             *statesmanager.StateManager
	backButton               component.Button
	unitListSection          UnitListSection
	userListSection          UserListSection
	descriptionSection       DescriptionSection
	actionSection            ActionSection
	addActionSection         AddActionSection
	removeActionSection      RemoveActionSection
	sendMessageSection       SendMessageSection
	trackUserLocationSection TrackUserLocationSection
	errorSection             ErrorSection
	infoSection              InfoSection
}

type UnitListSection struct {
	units           []*proto.Unit
	userToUnitCache map[string]string // userID->unitID
}
type UserListSection struct {
	users                []*proto.User
	userInformation      map[string]*proto.UserInformation
	usersList            component.ListSlider
	lastProcessedUserIdx int32
	isInUnit             bool
	currSelectedUserID   string
}

type DescriptionSection struct {
	descriptionBounds     rl.Rectangle
	descriptionID         string
	descriptionEmail      string
	descriptionName       string
	descriptionSurname    string
	descriptionLVL        string
	descriptionUnitPart   string
	descriptionDevicePart string
	descriptionTaskPart   string
}
type ActionSection struct {
	actionButtonArea    rl.Rectangle
	inUnitBackground    rl.Rectangle
	notInUnitBackground rl.Rectangle
	addButton           component.Button
	showAddModal        bool
	removeButton        component.Button
	showRemoveModal     bool
	inboxButton         component.Button
	showInboxModal      bool
	trackLocation       component.Button
	showLocationModal   bool
}

type AddActionSection struct {
	isConfirmAddButtonPressed bool
	unitsToAssignSlider       component.ListSlider
	acceptAddButton           component.Button
	addModal                  component.Modal
}
type RemoveActionSection struct {
	isConfirmRemoveButtonPressed bool
	acceptRemoveButton           component.Button
	removeModal                  component.Modal
}
type SendMessageSection struct {
	inboxModal                 component.Modal
	inboxInput                 component.InputBox
	sendMessage                component.Button
	activeUserCircle           component.Circle
	isSendMessageButtonPressed bool
}

type TrackUserLocationSection struct {
	LocationMap            LocationMap
	mapModal               component.Modal
	locationMapInformation component.LocationMapInformation
	userInfoTab            rl.Rectangle
	currentTaskTab         rl.Rectangle
}

type InfoSection struct {
	infoMessage string
	infoPopup   component.Popup
}
type ErrorSection struct {
	errorMessage string
	errorPopup   component.Popup
}

func (i *InfoUserScene) InfoUserSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	i.cfg = cfg
	i.stateManager = state
	i.userListSection.users = make([]*proto.User, 0, 32)
	i.unitListSection.units = make([]*proto.Unit, 0, 32)
	i.userListSection.userInformation = make(map[string]*proto.UserInformation)

	i.Reset()

	i.userListSection.usersList = component.ListSlider{
		Strings: make([]string, 0, 64),
		Bounds: rl.NewRectangle(
			0,
			0,
			(2.0/9.0)*float32(rl.GetScreenWidth()),
			float32(rl.GetScreenHeight())),
		IdxActiveElement: -1, // ?
		Focus:            0,
		IdxScroll:        0,
	}
	//TODO maybe check in all places to start -1

	i.descriptionSection.descriptionBounds = rl.NewRectangle(
		i.userListSection.usersList.Bounds.Width,
		i.userListSection.usersList.Bounds.Y,
		(7.0/9.0)*float32(rl.GetScreenWidth()),
		(7.0/9.0)*float32(rl.GetScreenHeight()),
	)
	i.actionSection.actionButtonArea = rl.NewRectangle(
		i.descriptionSection.descriptionBounds.X,
		i.descriptionSection.descriptionBounds.Y+i.descriptionSection.descriptionBounds.Height,
		i.descriptionSection.descriptionBounds.Width,
		(2.0/9.0)*float32(rl.GetScreenHeight()))

	var padding float32 = 20
	var btnWidth float32 = 120
	var btnHeight float32 = 65

	startX := i.actionSection.actionButtonArea.X + padding*2
	startY := i.actionSection.actionButtonArea.Y + padding*2

	// add to unit (+)
	i.actionSection.addButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(startX, startY, btnWidth, btnHeight),
		"+", false)

	// remove from unit (-)
	i.actionSection.removeButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.actionSection.addButton.Bounds.X+i.actionSection.addButton.Bounds.Width+padding,
			startY, btnWidth, btnHeight),
		"-", false)

	i.actionSection.inUnitBackground = rl.NewRectangle(
		i.actionSection.addButton.Bounds.X,
		i.actionSection.addButton.Bounds.Y,
		i.actionSection.addButton.Bounds.Width,
		i.actionSection.addButton.Bounds.Height)

	i.actionSection.notInUnitBackground = rl.NewRectangle(
		i.actionSection.removeButton.Bounds.X,
		i.actionSection.removeButton.Bounds.Y,
		i.actionSection.removeButton.Bounds.Width,
		i.actionSection.removeButton.Bounds.Height)
	// send message
	i.actionSection.inboxButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.actionSection.removeButton.Bounds.X+i.actionSection.removeButton.Bounds.Width+padding,
			startY, btnWidth, btnHeight),
		"Send message!", false)

	// location
	i.actionSection.trackLocation = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.actionSection.inboxButton.Bounds.X+i.actionSection.inboxButton.Bounds.Width+padding,
			startY, btnWidth, btnHeight),
		"Location", false)

	i.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()-110),
			float32(rl.GetScreenHeight()-68),
			100,
			50,
		),
		"Go back",
		false)
	if len(i.userListSection.users) > 0 {
		i.userListSection.usersList.IdxActiveElement = 0
	} else {
		i.userListSection.usersList.IdxActiveElement = -1
	}
	i.addActionSection.addModal = component.Modal{
		Background: rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		BgColor:    rl.Fade(rl.Gray, 0.3),
		Core:       rl.NewRectangle(float32(rl.GetScreenWidth()/2-150.0), float32(rl.GetScreenHeight()/2-150.0), 300, 300),
	}

	i.addActionSection.unitsToAssignSlider = component.ListSlider{
		Strings: make([]string, 0, 64),
		Bounds: rl.NewRectangle(
			i.addActionSection.addModal.Core.X+4,
			i.addActionSection.addModal.Core.Y+50,
			(3.9/4.0)*float32(i.addActionSection.addModal.Core.Width),
			(1.5/4.0)*float32(i.addActionSection.addModal.Core.Height)),
		IdxActiveElement: -1, // ?
		Focus:            0,
		IdxScroll:        0,
	}

	i.addActionSection.acceptAddButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.addActionSection.unitsToAssignSlider.Bounds.X,
		i.addActionSection.unitsToAssignSlider.Bounds.Y+120,
		(3.9/4.0)*float32(i.addActionSection.addModal.Core.Width),
		30), "Add to this unit", false)

	i.removeActionSection.removeModal = component.Modal{
		Background: rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		BgColor:    rl.Fade(rl.Gray, 0.3),
		Core:       rl.NewRectangle(float32(rl.GetScreenWidth()/2-150.0), float32(rl.GetScreenHeight()/2-150.0), 300, 300),
	}

	i.removeActionSection.acceptRemoveButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.removeActionSection.removeModal.Core.X+4,
		i.removeActionSection.removeModal.Core.Y+50,
		(3.9/4.0)*float32(i.removeActionSection.removeModal.Core.Width),
		30), "Remove from unit", false)

	popupRect := rl.NewRectangle(
		i.removeActionSection.removeModal.Core.X+20,
		i.removeActionSection.removeModal.Core.Y+(3.0/4.0)*i.removeActionSection.removeModal.Core.Height,
		250,
		40,
	)
	i.errorSection.errorMessage = ""
	i.errorSection.errorPopup = *component.NewPopup(
		component.NewPopupConfig(component.WithBgColor(rl.Red), component.WithFontColor(rl.White)),
		popupRect,
		&i.errorSection.errorMessage,
	)
	i.infoSection.infoMessage = ""
	i.infoSection.infoPopup = *component.NewPopup(
		component.NewPopupConfig(component.WithBgColor(rl.Green), component.WithFontColor(rl.Black)),
		popupRect,
		&i.infoSection.infoMessage,
	)

	i.sendMessageSection.inboxModal = component.Modal{
		Background: rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		BgColor:    rl.Fade(rl.Gray, 0.3),
		Core:       rl.NewRectangle(float32(rl.GetScreenWidth()/2-150.0), float32(rl.GetScreenHeight()/2-150.0), 400, 300),
	}

	i.sendMessageSection.inboxInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			i.sendMessageSection.inboxModal.Core.X+10,
			i.sendMessageSection.inboxModal.Core.Y+100,
			300,
			40))

	i.sendMessageSection.sendMessage = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			i.sendMessageSection.inboxInput.Bounds.X+i.sendMessageSection.inboxInput.Bounds.Width+10,
			i.sendMessageSection.inboxInput.Bounds.Y,
			50,
			50),
		"Send!", false)

	i.sendMessageSection.activeUserCircle = component.Circle{
		X:      int32(i.sendMessageSection.inboxModal.Core.X + 10),
		Y:      int32(i.sendMessageSection.inboxModal.Core.Y),
		Radius: 10,
	}
	i.trackUserLocationSection.locationMapInformation = component.LocationMapInformation{
		MapCurrentTask:    make(map[string]*component.CurrentTaskTab),
		MapPinInformation: make(map[string]*component.PinInformation),
	}
	i.trackUserLocationSection.mapModal = component.Modal{
		Background: rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		BgColor:    rl.Fade(rl.Gray, 0.3),
		Core:       rl.NewRectangle(float32(rl.GetScreenWidth()/2-400.0), float32(rl.GetScreenHeight()/2-225.0), 800, 450),
	}
	i.trackUserLocationSection.LocationMap = LocationMap{
		width:  800,
		height: 450,
		tm:     NewTileManager(),
	}

	var boxHeight float32 = 120
	i.trackUserLocationSection.currentTaskTab = rl.NewRectangle(
		i.trackUserLocationSection.mapModal.Core.X,
		i.trackUserLocationSection.mapModal.Core.Y+i.trackUserLocationSection.mapModal.Core.Height+5,
		i.trackUserLocationSection.mapModal.Core.Width,
		boxHeight,
	)
	i.trackUserLocationSection.userInfoTab = rl.NewRectangle(
		i.trackUserLocationSection.mapModal.Core.X,
		i.trackUserLocationSection.mapModal.Core.Y-boxHeight-5,
		i.trackUserLocationSection.mapModal.Core.Width,
		boxHeight,
	)
	i.FetchUsers()
	i.FetchUnits()
	i.GetUserInformation()
	for _, user := range i.userListSection.users {
		i.userListSection.usersList.Strings = append(i.userListSection.usersList.Strings, user.Personal.Name+"\n"+user.Personal.Surname)
	}
	for _, unit := range i.unitListSection.units {
		i.addActionSection.unitsToAssignSlider.Strings = append(i.addActionSection.unitsToAssignSlider.Strings, unit.Id)
	}
	i.prepareMap()

}

func (i *InfoUserScene) UpdateInfoUserState() {
	modalAddOpen := i.actionSection.showAddModal
	modalRemoveOpen := i.actionSection.showRemoveModal
	modalSendOpen := i.actionSection.showInboxModal
	modalLocationOpen := i.actionSection.showLocationModal
	cond := !modalAddOpen && !modalRemoveOpen && !modalSendOpen && !modalLocationOpen
	i.sendMessageSection.inboxInput.SetActive(!cond)
	i.actionSection.addButton.SetActive(cond)
	i.actionSection.removeButton.SetActive(cond)
	i.actionSection.inboxButton.SetActive(cond)
	i.actionSection.trackLocation.SetActive(cond)
	//i.addActionSection.acceptAddButton.SetActive(cond)
	//i.removeActionSection.acceptRemoveButton.SetActive(cond)
	//i.sendMessageSection.sendMessage.SetActive(cond)
	i.backButton.SetActive(cond)
	if cond {
		i.userListSection.usersList.Focus = 1
	} else {
		i.userListSection.usersList.Focus = 0
	}

	i.sendMessageSection.inboxInput.Update()

	if i.actionSection.addButton.Update() {
		i.actionSection.showAddModal = true
	}
	if i.actionSection.removeButton.Update() {
		i.actionSection.showRemoveModal = true
	}
	if i.actionSection.inboxButton.Update() {
		i.actionSection.showInboxModal = true
	}
	if i.actionSection.trackLocation.Update() {
		i.FetchPins()
		i.actionSection.showLocationModal = true
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		i.actionSection.showLocationModal = false
	}
	i.addActionSection.isConfirmAddButtonPressed = i.addActionSection.acceptAddButton.Update()
	i.removeActionSection.isConfirmRemoveButtonPressed = i.removeActionSection.acceptRemoveButton.Update()
	i.sendMessageSection.isSendMessageButtonPressed = i.sendMessageSection.sendMessage.Update()
	if i.backButton.Update() {
		i.stateManager.Add(statesmanager.GoBackState)
		return
	}
	if i.actionSection.showLocationModal {
		i.updateMap()
	}
	i.UpdateDescription()
	//TODO in v2 version add ability to have more than one unit by commanders type
	//and here change layout when he has more than one unit modal shows up with all units
	//and we have to choose unit to perform chose action
	i.AddToUnit()
	i.RemoveFromUnit()
	i.SendMessage()

}

func (i *InfoUserScene) RenderInfoUserState() {
	rl.ClearBackground(rl.White)

	upperBox := rl.NewRectangle(
		i.descriptionSection.descriptionBounds.X+2,
		i.descriptionSection.descriptionBounds.Y+5,
		i.descriptionSection.descriptionBounds.Width-5,
		i.descriptionSection.descriptionBounds.Height,
	)
	rl.DrawRectangle(int32(upperBox.X), int32(upperBox.Y), int32(upperBox.Width), int32(upperBox.Height), utils.USERDESCBG)
	rl.DrawRectangleLinesEx(upperBox, 2, rl.NewColor(203, 212, 205, 255))

	infoItems := []struct {
		label string
		value string
	}{
		{"USER ID:", i.descriptionSection.descriptionID},
		{"USER EMAIL:", i.descriptionSection.descriptionEmail},
		{"USER NAME:", i.descriptionSection.descriptionName},
		{"USER SURNAME:", i.descriptionSection.descriptionSurname},
		{"USER RULE LEVEL:", i.descriptionSection.descriptionLVL},
		{"DEVICE:", i.descriptionSection.descriptionDevicePart},
		{"TASK:", i.descriptionSection.descriptionTaskPart},
		{"UNIT:", i.descriptionSection.descriptionUnitPart},
	}

	itemHeight := float32(45)
	padding := float32(10)
	labelWidth := float32(210)
	startY := upperBox.Y + 80

	for idx, item := range infoItems {
		y := startY + float32(idx)*(itemHeight+padding)

		labelRect := rl.NewRectangle(
			upperBox.X+padding,
			y,
			labelWidth,
			itemHeight,
		)
		rl.DrawRectangle(int32(labelRect.X), int32(labelRect.Y), int32(labelRect.Width), int32(labelRect.Height), rl.NewColor(220, 220, 220, 255))
		rl.DrawRectangleLinesEx(labelRect, 2, rl.DarkGray)
		rl.DrawText(item.label, int32(labelRect.X+8), int32(labelRect.Y+10), 20, rl.Black)

		valueRect := rl.NewRectangle(
			labelRect.X+labelRect.Width+padding,
			y,
			upperBox.Width-3*padding-labelWidth,
			itemHeight,
		)
		if item.label == "UNIT:" {
			valueRect.Height += 20
		}
		rl.DrawRectangle(int32(valueRect.X), int32(valueRect.Y), int32(valueRect.Width), int32(valueRect.Height), rl.NewColor(245, 245, 245, 255))
		rl.DrawRectangleLinesEx(valueRect, 2, rl.Gray)
		rl.DrawText(item.value, int32(valueRect.X+8), int32(valueRect.Y+10), 20, rl.Black)
	}

	lowerBox := rl.NewRectangle(
		i.descriptionSection.descriptionBounds.X+2,
		i.descriptionSection.descriptionBounds.Y+i.descriptionSection.descriptionBounds.Height+10,
		i.descriptionSection.descriptionBounds.Width-5,
		float32(rl.GetScreenHeight())-i.descriptionSection.descriptionBounds.Height-20,
	)
	rl.DrawRectangle(int32(lowerBox.X), int32(lowerBox.Y), int32(lowerBox.Width), int32(lowerBox.Height), utils.USERBUTTONSBG)
	rl.DrawRectangleLinesEx(lowerBox, 2, rl.NewColor(203, 212, 205, 255))

	i.backButton.Render()

	i.actionSection.inboxButton.Render()
	i.actionSection.trackLocation.Render()
	if i.userListSection.usersList.Focus == 0 {
		gui.Disable()
		gui.ListViewEx(
			i.userListSection.usersList.Bounds,
			i.userListSection.usersList.Strings,
			&i.userListSection.usersList.IdxScroll,
			&i.userListSection.usersList.IdxActiveElement,
			i.userListSection.usersList.Focus,
		)
		gui.Enable()
	} else {
		gui.ListViewEx(
			i.userListSection.usersList.Bounds,
			i.userListSection.usersList.Strings,
			&i.userListSection.usersList.IdxScroll,
			&i.userListSection.usersList.IdxActiveElement,
			i.userListSection.usersList.Focus,
		)
	}

	if i.actionSection.showAddModal {
		rl.DrawRectangle(int32(i.addActionSection.addModal.Background.X),
			int32(i.addActionSection.addModal.Background.Y),
			int32(i.addActionSection.addModal.Background.Width),
			int32(i.addActionSection.addModal.Background.Height),
			i.addActionSection.addModal.BgColor)
		if gui.WindowBox(i.addActionSection.addModal.Core, "Add user to the unit") {
			i.actionSection.showAddModal = false
			i.addActionSection.unitsToAssignSlider.Strings = i.addActionSection.unitsToAssignSlider.Strings[:0]
			i.errorSection.errorPopup.Hide()
			i.infoSection.infoPopup.Hide()
		}
		gui.ListViewEx(
			i.addActionSection.unitsToAssignSlider.Bounds,
			i.addActionSection.unitsToAssignSlider.Strings,
			&i.addActionSection.unitsToAssignSlider.IdxScroll,
			&i.addActionSection.unitsToAssignSlider.IdxActiveElement,
			i.addActionSection.unitsToAssignSlider.Focus)
		i.addActionSection.acceptAddButton.Render()
	}

	if i.actionSection.showRemoveModal {
		rl.DrawRectangle(int32(i.removeActionSection.removeModal.Background.X),
			int32(i.removeActionSection.removeModal.Background.Y),
			int32(i.removeActionSection.removeModal.Background.Width),
			int32(i.removeActionSection.removeModal.Background.Height),
			i.removeActionSection.removeModal.BgColor)
		if gui.WindowBox(i.removeActionSection.removeModal.Core, "Remove user from unit") {
			i.actionSection.showRemoveModal = false
			i.errorSection.errorPopup.Hide()
			i.infoSection.infoPopup.Hide()
		}
		i.removeActionSection.acceptRemoveButton.Render()
	}

	if i.actionSection.showInboxModal {
		if gui.WindowBox(i.sendMessageSection.inboxModal.Core, "Send message") {
			i.actionSection.showInboxModal = false
		}
		rl.DrawCircle(
			i.sendMessageSection.activeUserCircle.X+int32(i.sendMessageSection.activeUserCircle.Radius),
			i.sendMessageSection.activeUserCircle.Y+4*int32(i.sendMessageSection.activeUserCircle.Radius),
			i.sendMessageSection.activeUserCircle.Radius,
			i.sendMessageSection.activeUserCircle.Color)
		i.sendMessageSection.sendMessage.Render()
		i.sendMessageSection.inboxInput.Render()
	}

	if i.actionSection.showLocationModal {
		rl.DrawRectangle(int32(i.trackUserLocationSection.mapModal.Background.X),
			int32(i.trackUserLocationSection.mapModal.Background.Y),
			int32(i.trackUserLocationSection.mapModal.Background.Width),
			int32(i.trackUserLocationSection.mapModal.Background.Height),
			i.trackUserLocationSection.mapModal.BgColor)
		rl.BeginScissorMode(int32(i.trackUserLocationSection.mapModal.Core.X),
			int32(i.trackUserLocationSection.mapModal.Core.Y),
			int32(i.trackUserLocationSection.mapModal.Core.Width),
			int32(i.trackUserLocationSection.mapModal.Core.Height))
		mousePos := i.drawMap()
		rl.EndScissorMode()
		i.showTabInformationOnCollision(mousePos)
		rl.DrawRectangleLines(
			int32(i.trackUserLocationSection.mapModal.Core.X),
			int32(i.trackUserLocationSection.mapModal.Core.Y),
			int32(i.trackUserLocationSection.mapModal.Core.Width),
			int32(i.trackUserLocationSection.mapModal.Core.Height),
			rl.Black)
	}
	if !i.userListSection.isInUnit {
		i.actionSection.addButton.Render()
		rl.DrawRectangle(
			int32(i.actionSection.notInUnitBackground.X),
			int32(i.actionSection.notInUnitBackground.Y),
			int32(i.actionSection.notInUnitBackground.Width),
			int32(i.actionSection.notInUnitBackground.Height),
			rl.Gray)
		rl.DrawText(
			"User is not \n in unit",
			int32(i.actionSection.notInUnitBackground.X),
			int32(i.actionSection.notInUnitBackground.Y),
			16,
			rl.White)
	} else {
		i.actionSection.removeButton.Render()
		rl.DrawRectangle(int32(i.actionSection.inUnitBackground.X),
			int32(i.actionSection.inUnitBackground.Y),
			int32(i.actionSection.inUnitBackground.Width),
			int32(i.actionSection.inUnitBackground.Height),
			rl.Gray)
		rl.DrawText(
			"User is \n in unit",
			int32(i.actionSection.inUnitBackground.X),
			int32(i.actionSection.inUnitBackground.Y),
			16,
			rl.White)

	}
	i.errorSection.errorPopup.Render()
	i.infoSection.infoPopup.Render()

}

//BIG TODO: remove a currently logged in user from e.g user info (cant send to myself message)
// and from other places
/*

 */
