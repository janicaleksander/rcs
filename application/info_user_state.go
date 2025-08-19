package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
)

type InfoUserScene struct {
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
	usersList            ListSlider
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
	unitsToAssignSlider       ListSlider
	acceptAddButton           component.Button
	addModal                  Modal
}
type RemoveActionSection struct {
	isConfirmRemoveButtonPressed bool
	usersUnitsSlider             ListSlider
	acceptRemoveButton           component.Button
	removeModal                  Modal
}
type SendMessageSection struct {
	inboxModal                 Modal
	inboxInput                 component.InputBox
	sendMessage                component.Button
	activeUserCircle           Circle
	isSendMessageButtonPressed bool
}

func (w *Window) InfoUserSceneSetup() {
	w.infoUserScene.Reset()
	w.FetchUnits()
	//TODO get proper lvl
	w.FetchUsers()
	w.infoUserScene.userListSection.usersList = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			0,
			0,
			(2.0/9.0)*float32(w.width),
			float32(w.height)),
		idxActiveElement: -1, // ?
		focus:            0,
		idxScroll:        -1,
	}
	//TODO maybe check in all places to start -1
	for _, user := range w.infoUserScene.userListSection.users {
		w.infoUserScene.userListSection.usersList.strings = append(w.infoUserScene.userListSection.usersList.strings, user.Personal.Name+"\n"+user.Personal.Surname)
	}

	w.infoUserScene.descriptionSection.descriptionBounds = rl.NewRectangle(
		w.infoUserScene.userListSection.usersList.bounds.Width,
		w.infoUserScene.userListSection.usersList.bounds.Y,
		(7.0/9.0)*float32(w.width),
		(7.0/9.0)*float32(w.height),
	)

	w.infoUserScene.actionSection.actionButtonArea = rl.NewRectangle(
		w.infoUserScene.descriptionSection.descriptionBounds.X,
		w.infoUserScene.descriptionSection.descriptionBounds.Y+w.infoUserScene.descriptionSection.descriptionBounds.Height,
		w.infoUserScene.descriptionSection.descriptionBounds.Width,
		(2.0/9.0)*float32(w.height))
	var padding float32 = 80
	//add to unit button
	w.infoUserScene.actionSection.addButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUserScene.actionSection.actionButtonArea.X+padding,
		w.infoUserScene.actionSection.actionButtonArea.Y,
		100,
		80), "+", false)

	w.infoUserScene.actionSection.inUnitBackground = rl.NewRectangle(
		w.infoUserScene.actionSection.addButton.Bounds.X,
		w.infoUserScene.actionSection.addButton.Bounds.Y,
		w.infoUserScene.actionSection.addButton.Bounds.Width,
		w.infoUserScene.actionSection.addButton.Bounds.Height)

	//remove from unit
	w.infoUserScene.actionSection.removeButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUserScene.actionSection.actionButtonArea.X+padding+w.infoUserScene.actionSection.addButton.Bounds.Width,
		w.infoUserScene.actionSection.actionButtonArea.Y,
		100,
		80), "-", false)

	w.infoUserScene.actionSection.notInUnitBackground = rl.NewRectangle(
		w.infoUserScene.actionSection.removeButton.Bounds.X,
		w.infoUserScene.actionSection.removeButton.Bounds.Y,
		w.infoUserScene.actionSection.removeButton.Bounds.Width,
		w.infoUserScene.actionSection.removeButton.Bounds.Height)

	w.infoUserScene.actionSection.inboxButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUserScene.actionSection.removeButton.Bounds.X+padding,
		w.infoUserScene.actionSection.removeButton.Bounds.Y,
		w.infoUserScene.actionSection.removeButton.Bounds.Width,
		w.infoUserScene.actionSection.removeButton.Bounds.Height), "Send message!", false)

	if len(w.infoUserScene.userListSection.users) > 0 {
		w.infoUserScene.userListSection.usersList.idxActiveElement = 0
	} else {
		w.infoUserScene.userListSection.usersList.idxActiveElement = -1
	}
	//TODO make one rule with ruleLVL when i can add what lvl and what lvl can do sth
	//e.g lvl5,lvl4 can only add lvl5; lvl4 can only add a lvl 3 2 1
	//and maybe here not include lvl 3 2 1(soldiers type)
	//or we cant add 5lvl to units cause their have access everywhere
	//POPUP after add button (sliders with units)

	w.infoUserScene.addActionSection.addModal = Modal{
		background: rl.NewRectangle(0, 0, float32(w.width), float32(w.height)),
		bgColor:    rl.Fade(rl.Gray, 0.3),
		core:       rl.NewRectangle(float32(w.width/2-150.0), float32(w.height/2-150.0), 300, 300),
	}

	w.infoUserScene.addActionSection.unitsToAssignSlider = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			w.infoUserScene.addActionSection.addModal.core.X+4,
			w.infoUserScene.addActionSection.addModal.core.Y+50,
			(3.9/4.0)*float32(w.infoUserScene.addActionSection.addModal.core.Width),
			(2.5/4.0)*float32(w.infoUserScene.addActionSection.addModal.core.Height)),
		idxActiveElement: -1, // ?
		focus:            0,
		idxScroll:        0,
	}

	w.infoUserScene.addActionSection.acceptAddButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUserScene.addActionSection.unitsToAssignSlider.bounds.X,
		w.infoUserScene.addActionSection.unitsToAssignSlider.bounds.Y+200,
		(3.9/4.0)*float32(w.infoUserScene.addActionSection.addModal.core.Width),
		30), "Add to this unit", false)

	w.infoUserScene.removeActionSection.removeModal = Modal{
		background: rl.NewRectangle(0, 0, float32(w.width), float32(w.height)),
		bgColor:    rl.Fade(rl.Gray, 0.3),
		core:       rl.NewRectangle(float32(w.width/2-150.0), float32(w.height/2-150.0), 300, 300),
	}
	w.infoUserScene.removeActionSection.usersUnitsSlider = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			w.infoUserScene.removeActionSection.removeModal.core.X+4,
			w.infoUserScene.removeActionSection.removeModal.core.Y+50,
			(3.9/4.0)*float32(w.infoUserScene.removeActionSection.removeModal.core.Width),
			(2.5/4.0)*float32(w.infoUserScene.removeActionSection.removeModal.core.Height)),
		idxActiveElement: -1, // ?
		focus:            0,
		idxScroll:        0,
	}
	w.infoUserScene.removeActionSection.acceptRemoveButton = *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
		w.infoUserScene.removeActionSection.usersUnitsSlider.bounds.X,
		w.infoUserScene.removeActionSection.usersUnitsSlider.bounds.Y+200,
		(3.9/4.0)*float32(w.infoUserScene.removeActionSection.removeModal.core.Width),
		30), "Remove from unit", false)

	w.infoUserScene.sendMessageSection.inboxModal = Modal{
		background: rl.NewRectangle(0, 0, float32(w.width), float32(w.height)),
		bgColor:    rl.Fade(rl.Gray, 0.3),
		core:       rl.NewRectangle(float32(w.width/2-150.0), float32(w.height/2-150.0), 400, 200),
	}

	w.infoUserScene.sendMessageSection.inboxInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			w.infoUserScene.sendMessageSection.inboxModal.core.X+10,
			w.infoUserScene.sendMessageSection.inboxModal.core.Y+100,
			300,
			40))

	w.infoUserScene.sendMessageSection.sendMessage = *component.NewButton(component.NewButtonConfig(),
		rl.NewRectangle(
			w.infoUserScene.sendMessageSection.inboxInput.Bounds.X+w.infoUserScene.sendMessageSection.inboxInput.Bounds.Width+10,
			w.infoUserScene.sendMessageSection.inboxInput.Bounds.Y,
			50,
			50),
		"Send!", false)

	w.infoUserScene.sendMessageSection.activeUserCircle = Circle{
		x:      int32(w.infoUserScene.sendMessageSection.inboxModal.core.X + 10),
		y:      int32(w.infoUserScene.sendMessageSection.inboxModal.core.Y),
		radius: 10,
	}

	w.infoUserScene.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			float32(w.width-100),
			float32(w.height-50),
			100,
			50,
		),
		"Go back",
		false)

}

func (w *Window) updateInfoUserState() {
	//TODO Maybe add some control to slider etc their focus
	modalAddOpen := w.infoUserScene.actionSection.showAddModal
	modalRemoveOpen := w.infoUserScene.actionSection.showRemoveModal
	modalSendOpen := w.infoUserScene.actionSection.showInboxModal
	cond := !modalAddOpen && !modalRemoveOpen && !modalSendOpen
	w.infoUserScene.sendMessageSection.inboxInput.SetActive(cond)
	w.infoUserScene.actionSection.addButton.SetActive(cond)
	w.infoUserScene.actionSection.removeButton.SetActive(cond)
	w.infoUserScene.actionSection.inboxButton.SetActive(cond)
	w.infoUserScene.addActionSection.acceptAddButton.SetActive(cond)
	w.infoUserScene.removeActionSection.acceptRemoveButton.SetActive(cond)
	w.infoUserScene.sendMessageSection.sendMessage.SetActive(cond)
	w.infoUserScene.backButton.SetActive(cond)

	w.infoUserScene.sendMessageSection.inboxInput.Update()

	if w.infoUserScene.actionSection.addButton.Update() {
		w.infoUserScene.actionSection.showAddModal = true
	}
	if w.infoUserScene.actionSection.removeButton.Update() {
		w.infoUserScene.actionSection.showRemoveModal = true
	}
	if w.infoUserScene.actionSection.inboxButton.Update() {
		w.infoUserScene.actionSection.showInboxModal = true
	}

	w.infoUserScene.addActionSection.isConfirmAddButtonPressed = w.infoUserScene.addActionSection.acceptAddButton.Update()
	w.infoUserScene.removeActionSection.isConfirmRemoveButtonPressed = w.infoUserScene.removeActionSection.acceptRemoveButton.Update()
	w.infoUserScene.sendMessageSection.isSendMessageButtonPressed = w.infoUserScene.sendMessageSection.sendMessage.Update()
	if w.infoUserScene.backButton.Update() {
		w.goSceneBack()
		return
	}

	w.UpdateDescription()
	//TODO in v2 version add ability to have more than one unit by commanders type
	//and here change layout when he has more than one unit modal shows up with all units
	//and we have to choose unit to perform chose action
	w.AddToUnit()
	w.RemoveFromUnit()
	w.SendMessage()
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
func (w *Window) renderInfoUserState() {
	w.infoUserScene.backButton.Render()
	w.infoUserScene.actionSection.addButton.Render()
	w.infoUserScene.actionSection.removeButton.Render()
	w.infoUserScene.actionSection.inboxButton.Render()

	gui.ListViewEx(
		w.infoUserScene.userListSection.usersList.bounds,
		w.infoUserScene.userListSection.usersList.strings,
		&w.infoUserScene.userListSection.usersList.idxScroll,
		&w.infoUserScene.userListSection.usersList.idxActiveElement,
		w.infoUserScene.userListSection.usersList.focus)
	rl.DrawRectangle(
		int32(w.infoUserScene.descriptionSection.descriptionBounds.X),
		int32(w.infoUserScene.descriptionSection.descriptionBounds.Y),
		int32(w.infoUserScene.descriptionSection.descriptionBounds.Width),
		int32(w.infoUserScene.descriptionSection.descriptionBounds.Height),
		rl.White)

	rl.DrawText(
		w.infoUserScene.descriptionSection.descriptionName+"\n"+
			w.infoUserScene.descriptionSection.descriptionSurname+"\n"+
			w.infoUserScene.descriptionSection.descriptionLVL+"\n",
		int32(w.infoUserScene.descriptionSection.descriptionBounds.X),
		int32(w.infoUserScene.descriptionSection.descriptionBounds.Y),
		43, rl.Yellow)

	if w.infoUserScene.actionSection.showAddModal {
		rl.DrawRectangle(
			int32(w.infoUserScene.addActionSection.addModal.background.X),
			int32(w.infoUserScene.addActionSection.addModal.background.Y),
			int32(w.infoUserScene.addActionSection.addModal.background.Width),
			int32(w.infoUserScene.addActionSection.addModal.background.Height),
			w.infoUserScene.addActionSection.addModal.bgColor)
		if gui.WindowBox(w.infoUserScene.addActionSection.addModal.core, "TITLE") {
			w.infoUserScene.actionSection.showAddModal = false
			w.infoUserScene.addActionSection.unitsToAssignSlider.strings = w.infoUserScene.addActionSection.unitsToAssignSlider.strings[:0]
		}
		gui.ListViewEx(
			w.infoUserScene.addActionSection.unitsToAssignSlider.bounds,
			w.infoUserScene.addActionSection.unitsToAssignSlider.strings,
			&w.infoUserScene.addActionSection.unitsToAssignSlider.idxScroll,
			&w.infoUserScene.addActionSection.unitsToAssignSlider.idxActiveElement,
			w.infoUserScene.addActionSection.unitsToAssignSlider.focus)
		w.infoUserScene.addActionSection.acceptAddButton.Render()

	}
	if w.infoUserScene.actionSection.showRemoveModal {
		rl.DrawRectangle(
			int32(w.infoUserScene.removeActionSection.removeModal.background.X),
			int32(w.infoUserScene.removeActionSection.removeModal.background.Y),
			int32(w.infoUserScene.removeActionSection.removeModal.background.Width),
			int32(w.infoUserScene.removeActionSection.removeModal.background.Height),
			w.infoUserScene.removeActionSection.removeModal.bgColor)
		if gui.WindowBox(w.infoUserScene.removeActionSection.removeModal.core, "TITLE") {
			w.infoUserScene.actionSection.showRemoveModal = false
			w.infoUserScene.removeActionSection.usersUnitsSlider.strings = w.infoUserScene.removeActionSection.usersUnitsSlider.strings[:0]
		}

		gui.ListViewEx(w.infoUserScene.removeActionSection.usersUnitsSlider.bounds,
			w.infoUserScene.removeActionSection.usersUnitsSlider.strings,
			&w.infoUserScene.removeActionSection.usersUnitsSlider.idxScroll,
			&w.infoUserScene.removeActionSection.usersUnitsSlider.idxActiveElement,
			w.infoUserScene.removeActionSection.usersUnitsSlider.focus)
		w.infoUserScene.removeActionSection.acceptRemoveButton.Render()

	}

	if w.infoUserScene.actionSection.showInboxModal {
		if gui.WindowBox(w.infoUserScene.sendMessageSection.inboxModal.core, "TITLE") {
			w.infoUserScene.actionSection.showInboxModal = false
		}
		rl.DrawCircle(
			w.infoUserScene.sendMessageSection.activeUserCircle.x,
			w.infoUserScene.sendMessageSection.activeUserCircle.y,
			w.infoUserScene.sendMessageSection.activeUserCircle.radius,
			w.infoUserScene.sendMessageSection.activeUserCircle.color)

		w.infoUserScene.sendMessageSection.sendMessage.Render()
		w.infoUserScene.sendMessageSection.inboxInput.Render()

	}

}

//BIG TODO: remove a currently logged in user from e.g user info (cant send to myself message)
// and from other places
