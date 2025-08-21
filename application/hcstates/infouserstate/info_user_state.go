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
	cfg          *utils.SharedConfig
	stateManager *statesmanager.StateManager
	//scheduler TODO
	backButton          component.Button
	unitListSection     UnitListSection
	userListSection     UserListSection
	descriptionSection  DescriptionSection
	actionSection       ActionSection
	addActionSection    AddActionSection
	removeActionSection RemoveActionSection
	sendMessageSection  SendMessageSection
}

type UnitListSection struct {
	units           []*proto.Unit
	userToUnitCache map[string]string // userID->unitID
}
type UserListSection struct {
	users                []*proto.User
	usersList            component.ListSlider
	lastProcessedUserIdx int32
	currSelectedUserID   string
	isInUnit             bool
}

type DescriptionSection struct {
	descriptionBounds  rl.Rectangle
	descriptionName    string
	descriptionSurname string
	descriptionLVL     string
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
}

type AddActionSection struct {
	isConfirmAddButtonPressed bool
	unitsToAssignSlider       component.ListSlider
	acceptAddButton           component.Button
	addModal                  component.Modal
}
type RemoveActionSection struct {
	isConfirmRemoveButtonPressed bool
	usersUnitsSlider             component.ListSlider
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

func (i *InfoUserScene) InfoUserSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	i.cfg = cfg
	i.stateManager = state
	i.Reset()
	i.FetchUnits()
	//TODO get proper lvl
	i.FetchUsers()
	i.userListSection.usersList = component.ListSlider{
		Strings: make([]string, 0, 64),
		Bounds: rl.NewRectangle(
			0,
			0,
			(2.0/9.0)*float32(rl.GetScreenWidth()),
			float32(rl.GetScreenHeight())),
		IdxActiveElement: -1, // ?
		Focus:            0,
		IdxScroll:        -1,
	}
	//TODO maybe check in all places to start -1
	for _, user := range i.userListSection.users {
		i.userListSection.usersList.Strings = append(i.userListSection.usersList.Strings, user.Personal.Name+"\n"+user.Personal.Surname)
	}

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
	var padding float32 = 80
	//add to unit button
	i.actionSection.addButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.actionSection.actionButtonArea.X+padding,
		i.actionSection.actionButtonArea.Y,
		100,
		80), "+", false)

	i.actionSection.inUnitBackground = rl.NewRectangle(
		i.actionSection.addButton.Bounds.X,
		i.actionSection.addButton.Bounds.Y,
		i.actionSection.addButton.Bounds.Width,
		i.actionSection.addButton.Bounds.Height)

	//remove from unit
	i.actionSection.removeButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.actionSection.actionButtonArea.X+padding+i.actionSection.addButton.Bounds.Width,
		i.actionSection.actionButtonArea.Y,
		100,
		80), "-", false)

	i.actionSection.notInUnitBackground = rl.NewRectangle(
		i.actionSection.removeButton.Bounds.X,
		i.actionSection.removeButton.Bounds.Y,
		i.actionSection.removeButton.Bounds.Width,
		i.actionSection.removeButton.Bounds.Height)

	i.actionSection.inboxButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.actionSection.removeButton.Bounds.X+padding,
		i.actionSection.removeButton.Bounds.Y,
		i.actionSection.removeButton.Bounds.Width,
		i.actionSection.removeButton.Bounds.Height), "Send message!", false)

	if len(i.userListSection.users) > 0 {
		i.userListSection.usersList.IdxActiveElement = 0
	} else {
		i.userListSection.usersList.IdxActiveElement = -1
	}
	//TODO make one rule with ruleLVL when i can add what lvl and what lvl can do sth
	//e.g lvl5,lvl4 can only add lvl5; lvl4 can only add a lvl 3 2 1
	//and maybe here not include lvl 3 2 1(soldiers type)
	//or we cant add 5lvl to units cause their have access everywhere
	//POPUP after add button (sliders with units)

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
			(2.5/4.0)*float32(i.addActionSection.addModal.Core.Height)),
		IdxActiveElement: -1, // ?
		Focus:            0,
		IdxScroll:        0,
	}

	i.addActionSection.acceptAddButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.addActionSection.unitsToAssignSlider.Bounds.X,
		i.addActionSection.unitsToAssignSlider.Bounds.Y+200,
		(3.9/4.0)*float32(i.addActionSection.addModal.Core.Width),
		30), "Add to this unit", false)

	i.removeActionSection.removeModal = component.Modal{
		Background: rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		BgColor:    rl.Fade(rl.Gray, 0.3),
		Core:       rl.NewRectangle(float32(rl.GetScreenWidth()/2-150.0), float32(rl.GetScreenHeight()/2-150.0), 300, 300),
	}
	i.removeActionSection.usersUnitsSlider = component.ListSlider{
		Strings: make([]string, 0, 64),
		Bounds: rl.NewRectangle(
			i.removeActionSection.removeModal.Core.X+4,
			i.removeActionSection.removeModal.Core.Y+50,
			(3.9/4.0)*float32(i.removeActionSection.removeModal.Core.Width),
			(2.5/4.0)*float32(i.removeActionSection.removeModal.Core.Height)),
		IdxActiveElement: -1, // ?
		Focus:            0,
		IdxScroll:        0,
	}
	i.removeActionSection.acceptRemoveButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		i.removeActionSection.usersUnitsSlider.Bounds.X,
		i.removeActionSection.usersUnitsSlider.Bounds.Y+200,
		(3.9/4.0)*float32(i.removeActionSection.removeModal.Core.Width),
		30), "Remove from unit", false)

	i.sendMessageSection.inboxModal = component.Modal{
		Background: rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		BgColor:    rl.Fade(rl.Gray, 0.3),
		Core:       rl.NewRectangle(float32(rl.GetScreenWidth()/2-150.0), float32(rl.GetScreenHeight()/2-150.0), 400, 200),
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

	i.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			float32(rl.GetScreenWidth()-100),
			float32(rl.GetScreenHeight()-50),
			100,
			50,
		),
		"Go back",
		false)

}

func (i *InfoUserScene) UpdateInfoUserState() {
	//TODO Maybe add some control to slider etc their focus
	modalAddOpen := i.actionSection.showAddModal
	modalRemoveOpen := i.actionSection.showRemoveModal
	modalSendOpen := i.actionSection.showInboxModal
	cond := !modalAddOpen && !modalRemoveOpen && !modalSendOpen
	i.sendMessageSection.inboxInput.SetActive(cond)
	i.actionSection.addButton.SetActive(cond)
	i.actionSection.removeButton.SetActive(cond)
	i.actionSection.inboxButton.SetActive(cond)
	i.addActionSection.acceptAddButton.SetActive(cond)
	i.removeActionSection.acceptRemoveButton.SetActive(cond)
	i.sendMessageSection.sendMessage.SetActive(cond)
	i.backButton.SetActive(cond)

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

	i.addActionSection.isConfirmAddButtonPressed = i.addActionSection.acceptAddButton.Update()
	i.removeActionSection.isConfirmRemoveButtonPressed = i.removeActionSection.acceptRemoveButton.Update()
	i.sendMessageSection.isSendMessageButtonPressed = i.sendMessageSection.sendMessage.Update()
	if i.backButton.Update() {
		i.stateManager.Add(statesmanager.GoBackState)
		return
	}

	i.UpdateDescription()
	//TODO in v2 version add ability to have more than one unit by commanders type
	//and here change layout when he has more than one unit modal shows up with all units
	//and we have to choose unit to perform chose action
	i.AddToUnit()
	i.RemoveFromUnit()
	i.SendMessage()
}

//goback
// add to unit
// remove from unit
// send message

// modal remove
// modal add
// modal send

// slider one
// slider two
// input box
func (i *InfoUserScene) RenderInfoUserState() {
	i.backButton.Render()
	i.actionSection.addButton.Render()
	i.actionSection.removeButton.Render()
	i.actionSection.inboxButton.Render()

	gui.ListViewEx(
		i.userListSection.usersList.Bounds,
		i.userListSection.usersList.Strings,
		&i.userListSection.usersList.IdxScroll,
		&i.userListSection.usersList.IdxActiveElement,
		i.userListSection.usersList.Focus)
	rl.DrawRectangle(
		int32(i.descriptionSection.descriptionBounds.X),
		int32(i.descriptionSection.descriptionBounds.Y),
		int32(i.descriptionSection.descriptionBounds.Width),
		int32(i.descriptionSection.descriptionBounds.Height),
		rl.White)

	rl.DrawText(
		i.descriptionSection.descriptionName+"\n"+
			i.descriptionSection.descriptionSurname+"\n"+
			i.descriptionSection.descriptionLVL+"\n",
		int32(i.descriptionSection.descriptionBounds.X),
		int32(i.descriptionSection.descriptionBounds.Y),
		43, rl.Yellow)

	if i.actionSection.showAddModal {
		rl.DrawRectangle(
			int32(i.addActionSection.addModal.Background.X),
			int32(i.addActionSection.addModal.Background.Y),
			int32(i.addActionSection.addModal.Background.Width),
			int32(i.addActionSection.addModal.Background.Height),
			i.addActionSection.addModal.BgColor)
		if gui.WindowBox(i.addActionSection.addModal.Core, "TITLE") {
			i.actionSection.showAddModal = false
			i.addActionSection.unitsToAssignSlider.Strings = i.addActionSection.unitsToAssignSlider.Strings[:0]
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
		rl.DrawRectangle(
			int32(i.removeActionSection.removeModal.Background.X),
			int32(i.removeActionSection.removeModal.Background.Y),
			int32(i.removeActionSection.removeModal.Background.Width),
			int32(i.removeActionSection.removeModal.Background.Height),
			i.removeActionSection.removeModal.BgColor)
		if gui.WindowBox(i.removeActionSection.removeModal.Core, "TITLE") {
			i.actionSection.showRemoveModal = false
			i.removeActionSection.usersUnitsSlider.Strings = i.removeActionSection.usersUnitsSlider.Strings[:0]
		}

		gui.ListViewEx(i.removeActionSection.usersUnitsSlider.Bounds,
			i.removeActionSection.usersUnitsSlider.Strings,
			&i.removeActionSection.usersUnitsSlider.IdxScroll,
			&i.removeActionSection.usersUnitsSlider.IdxActiveElement,
			i.removeActionSection.usersUnitsSlider.Focus)
		i.removeActionSection.acceptRemoveButton.Render()

	}

	if i.actionSection.showInboxModal {
		if gui.WindowBox(i.sendMessageSection.inboxModal.Core, "TITLE") {
			i.actionSection.showInboxModal = false
		}
		rl.DrawCircle(
			i.sendMessageSection.activeUserCircle.X,
			i.sendMessageSection.activeUserCircle.Y,
			i.sendMessageSection.activeUserCircle.Radius,
			i.sendMessageSection.activeUserCircle.Color)

		i.sendMessageSection.sendMessage.Render()
		i.sendMessageSection.inboxInput.Render()

	}

}

//BIG TODO: remove a currently logged in user from e.g user info (cant send to myself message)
// and from other places
