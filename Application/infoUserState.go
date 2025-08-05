package Application

import (
	"fmt"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"strconv"
	"sync"
)

type InfoUserScene struct {
	// user list slider
	usersList       ListSlider
	users           []*Proto.User
	units           []*Proto.Unit
	userToUnitCache map[string]string // userID->unitID
	// description area
	descriptionBounds  rl.Rectangle
	descriptionName    string
	descriptionSurname string
	descriptionLVL     string
	currUserID         string
	isInUnit           bool

	// action button area
	actionButtonArea rl.Rectangle
	// add btn
	addButton                 Button
	isConfirmAddButtonPressed bool
	inUnitBackground          rl.Rectangle

	//TODO add errors field to two modals
	//modal after add btn
	showAddModal        bool
	unitsToAssignSlider ListSlider
	acceptAddButton     Button
	addModal            Modal

	// rmv btn
	removeButton                 Button
	isConfirmRemoveButtonPressed bool
	notInUnitBackground          rl.Rectangle

	//modal after rmv btn
	showRemoveModal    bool
	usersUnitsSlider   ListSlider
	acceptRemoveButton Button
	removeModal        Modal

	// inbox btn
	inboxButton          Button
	lastProcessedUserIdx int32
}

func (i *InfoUserScene) Reset() {
	i.lastProcessedUserIdx = -1
	i.descriptionName = ""
	i.descriptionSurname = ""
	i.descriptionLVL = ""
	i.currUserID = ""
}
func (w *Window) InfoUserSceneSetup() {
	w.infoUserScene.Reset()
	resp := w.ctx.Request(w.serverPID, &Proto.GetAllUnits{}, WaitTime)
	val, err := resp.Result()
	if err != nil {
		// TODO error
	}
	w.infoUserScene.units = make([]*Proto.Unit, 0, 64)
	if v, ok := val.(*Proto.AllUnits); ok {
		for _, unit := range v.Units {
			w.infoUserScene.units = append(w.infoUserScene.units, unit)
		}
	} else {
		// TODO error
	}

	//TODO get proper lvl
	resp = w.ctx.Request(w.serverPID, &Proto.GetUserAboveLVL{Lvl: -1}, WaitTime)
	val, err = resp.Result()
	if err != nil {
		// TODO error
	}
	w.infoUserScene.users = make([]*Proto.User, 0, 64)
	if v, ok := val.(*Proto.UsersAboveLVL); ok {
		for _, user := range v.Users {
			w.infoUserScene.users = append(w.infoUserScene.users, user)
		}
	} else {
		// TODO error
	}

	//cache users information
	w.infoUserScene.userToUnitCache = make(map[string]string, len(w.infoUserScene.users))
	var waitGroup sync.WaitGroup
	cacheChan := make(chan struct {
		userID string
		unitID string
	}, 1024)
	for _, user := range w.infoUserScene.users {
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup, userID string) {
			defer wg.Done()
			resp = w.ctx.Request(w.serverPID, &Proto.IsUserInUnit{Id: userID}, WaitTime)
			v, err := resp.Result()
			if err != nil {
				// TODO
				fmt.Println("ERROR")
			}
			if payload, ok := v.(*Proto.UserInUnit); ok {
				cacheChan <- struct {
					userID string
					unitID string
				}{userID: userID, unitID: payload.UnitID}
			}

		}(&waitGroup, user.Id)
	}
	go func() {
		waitGroup.Wait()
		close(cacheChan)
	}()
	for v := range cacheChan {
		w.infoUserScene.userToUnitCache[v.userID] = v.unitID
	}
	w.infoUserScene.usersList = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			0,
			0,
			(2.0/9.0)*float32(w.width),
			float32(w.height)),
		idxActiveElement: -1,
		focus:            0,
		idxScroll:        -1,
	}
	//TODO maybe check in all places to start -1
	for _, user := range w.infoUserScene.users {
		w.infoUserScene.usersList.strings = append(w.infoUserScene.usersList.strings, user.Personal.Name+"\n"+user.Personal.Surname)
	}
	w.infoUserScene.descriptionBounds = rl.NewRectangle(
		w.infoUserScene.usersList.bounds.Width,
		w.infoUserScene.usersList.bounds.Y,
		(7.0/9.0)*float32(w.width),
		(7.0/9.0)*float32(w.height),
	)
	w.infoUserScene.actionButtonArea = rl.NewRectangle(
		w.infoUserScene.descriptionBounds.X,
		w.infoUserScene.descriptionBounds.Y+w.infoUserScene.descriptionBounds.Height,
		w.infoUserScene.descriptionBounds.Width,
		(2.0/9.0)*float32(w.height))
	var padding float32 = 80
	//add to unit button
	w.infoUserScene.addButton = Button{
		bounds: rl.NewRectangle(
			w.infoUserScene.actionButtonArea.X+padding,
			w.infoUserScene.actionButtonArea.Y,
			100,
			80),
		text: "+",
	}
	w.infoUserScene.inUnitBackground = rl.NewRectangle(
		w.infoUserScene.addButton.bounds.X,
		w.infoUserScene.addButton.bounds.Y,
		w.infoUserScene.addButton.bounds.Width,
		w.infoUserScene.addButton.bounds.Height)
	//remove from unit
	w.infoUserScene.removeButton = Button{
		bounds: rl.NewRectangle(
			w.infoUserScene.actionButtonArea.X+padding+w.infoUserScene.addButton.bounds.Width,
			w.infoUserScene.actionButtonArea.Y,
			100,
			80),
		text: "-",
	}
	w.infoUserScene.notInUnitBackground = rl.NewRectangle(
		w.infoUserScene.removeButton.bounds.X,
		w.infoUserScene.removeButton.bounds.Y,
		w.infoUserScene.removeButton.bounds.Width,
		w.infoUserScene.removeButton.bounds.Height)

	if len(w.infoUserScene.users) > 0 {
		w.infoUserScene.usersList.idxActiveElement = 0
	} else {
		w.infoUserScene.usersList.idxActiveElement = -1

	}
	//TODO make one rule with ruleLVL when i can add what lvl and what lvl can do sth
	//e.g lvl5,lvl4 can only add lvl5; lvl4 can only add a lvl 3 2 1
	//and maybe here not include lvl 3 2 1(soldiers type)
	//or we cant add 5lvl to units cause their have access everywhere
	//POPUP after add button (sliders with units)

	w.infoUserScene.addModal = Modal{
		background: rl.NewRectangle(0, 0, float32(w.width), float32(w.height)),
		bgColor:    rl.Fade(rl.Gray, 0.3),
		core:       rl.NewRectangle(float32(w.width/2-150.0), float32(w.height/2-150.0), 300, 300),
	}
	w.infoUserScene.unitsToAssignSlider = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			w.infoUserScene.addModal.core.X+4,
			w.infoUserScene.addModal.core.Y+50,
			(3.9/4.0)*float32(w.infoUserScene.addModal.core.Width),
			(2.5/4.0)*float32(w.infoUserScene.addModal.core.Height)),
		idxActiveElement: -1,
		focus:            0,
		idxScroll:        0,
	}
	w.infoUserScene.acceptAddButton = Button{
		bounds: rl.NewRectangle(
			w.infoUserScene.unitsToAssignSlider.bounds.X,
			w.infoUserScene.unitsToAssignSlider.bounds.Y+200,
			(3.9/4.0)*float32(w.infoUserScene.addModal.core.Width),
			30),
		text: "Add to this unit",
	}

	//////
	w.infoUserScene.removeModal = Modal{
		background: rl.NewRectangle(0, 0, float32(w.width), float32(w.height)),
		bgColor:    rl.Fade(rl.Gray, 0.3),
		core:       rl.NewRectangle(float32(w.width/2-150.0), float32(w.height/2-150.0), 300, 300),
	}
	w.infoUserScene.usersUnitsSlider = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			w.infoUserScene.removeModal.core.X+4,
			w.infoUserScene.removeModal.core.Y+50,
			(3.9/4.0)*float32(w.infoUserScene.removeModal.core.Width),
			(2.5/4.0)*float32(w.infoUserScene.removeModal.core.Height)),
		idxActiveElement: -1,
		focus:            0,
		idxScroll:        0,
	}
	w.infoUserScene.acceptRemoveButton = Button{
		bounds: rl.NewRectangle(
			w.infoUserScene.usersUnitsSlider.bounds.X,
			w.infoUserScene.usersUnitsSlider.bounds.Y+200,
			(3.9/4.0)*float32(w.infoUserScene.removeModal.core.Width),
			30),
		text: "Remove from  this unit",
	}

}

func (w *Window) updateInfoUserState() {
	currentUserIdx := w.infoUserScene.usersList.idxActiveElement
	if currentUserIdx != -1 && currentUserIdx != w.infoUserScene.lastProcessedUserIdx {
		//
		user := w.infoUserScene.users[w.infoUserScene.usersList.idxActiveElement]
		w.infoUserScene.currUserID = user.Id
		//TODO in the v2 version we need to track more than
		// one unit ID
		if _, ok := w.infoUserScene.userToUnitCache[user.Id]; ok {
			w.infoUserScene.isInUnit = true
		} else {
			w.infoUserScene.isInUnit = false
		}
		//
		w.infoUserScene.descriptionName = user.Personal.Name
		w.infoUserScene.descriptionSurname = user.Personal.Surname
		w.infoUserScene.descriptionLVL = strconv.Itoa(int(user.RuleLvl))
		w.infoUserScene.lastProcessedUserIdx = currentUserIdx
	}
	//TODO in v2 version add ability to have more than one unit by commanders type
	//and here change layout when he has more than one unit modal shows up with all units
	//and we have to choose unit to perform chose action
	if !w.infoUserScene.isInUnit { // shows add to unit
		//fil userUnits ( TODO in v2 for loop through many units)
		if gui.Button(w.infoUserScene.addButton.bounds, w.infoUserScene.addButton.text) {
			for _, unit := range w.infoUserScene.units {
				w.infoUserScene.unitsToAssignSlider.strings = append(w.infoUserScene.unitsToAssignSlider.strings, unit.Id)
				/*
					cacheUnit := w.infoUserScene.userToUnitCache[w.infoUserScene.currUserID]
					if cacheUnit == unit.Id {
						continue
					}
					in v2 version. we dont need to show units that we are already enrolled in
				*/

			}
			w.infoUserScene.showAddModal = true

		}
		if w.infoUserScene.isConfirmAddButtonPressed {
			if w.infoUserScene.unitsToAssignSlider.idxActiveElement >= 0 {
				unit := w.infoUserScene.units[w.infoUserScene.unitsToAssignSlider.idxActiveElement]
				resp := w.ctx.Request(w.serverPID, &Proto.AssignUserToUnit{
					UserID: w.infoUserScene.currUserID,
					UnitID: unit.Id,
				}, WaitTime)
				val, err := resp.Result()
				if _, ok := val.(*Proto.FailureOfAssign); ok || err != nil {
					// TODO failure
				}
				if _, ok := val.(*Proto.SuccessOfAssign); ok {
					//TODO success
					w.infoUserScene.userToUnitCache[w.infoUserScene.currUserID] = unit.Id
					w.infoUserScene.isInUnit = true

				}
			}

		}
	} else {
		rl.DrawRectangle(int32(w.infoUserScene.inUnitBackground.X),
			int32(w.infoUserScene.inUnitBackground.Y),
			int32(w.infoUserScene.inUnitBackground.Width),
			int32(w.infoUserScene.inUnitBackground.Height),
			rl.Gray)
		rl.DrawText("User is \n in unit", int32(w.infoUserScene.inUnitBackground.X), int32(w.infoUserScene.inUnitBackground.Y), 16, rl.White)

	}
	if w.infoUserScene.isInUnit { // shows remove  unit
		if gui.Button(w.infoUserScene.removeButton.bounds, w.infoUserScene.removeButton.text) {
			w.infoUserScene.usersUnitsSlider.strings = append(w.infoUserScene.usersUnitsSlider.strings, w.infoUserScene.userToUnitCache[w.infoUserScene.currUserID])
			w.infoUserScene.showRemoveModal = true
		}
		if w.infoUserScene.isConfirmRemoveButtonPressed {
			if w.infoUserScene.usersUnitsSlider.idxActiveElement >= 0 {
				unit := w.infoUserScene.units[w.infoUserScene.usersUnitsSlider.idxActiveElement]
				resp := w.ctx.Request(w.serverPID, &Proto.DeleteUserFromUnit{
					UserID: w.infoUserScene.currUserID,
					UnitID: unit.Id,
				}, WaitTime)
				val, err := resp.Result()
				if _, ok := val.(*Proto.FailureOfDelete); ok || err != nil {
					// TODO failure
				}
				if _, ok := val.(*Proto.SuccessOfDelete); ok {
					//TODO success

					//TODO in v2 map str->[]str and then we have to iterate through
					// this slice and delete exact unit
					delete(w.infoUserScene.userToUnitCache, w.infoUserScene.currUserID)
					w.infoUserScene.isInUnit = false

				}
			}

		}
	} else {
		rl.DrawRectangle(int32(w.infoUserScene.notInUnitBackground.X),
			int32(w.infoUserScene.notInUnitBackground.Y),
			int32(w.infoUserScene.notInUnitBackground.Width),
			int32(w.infoUserScene.notInUnitBackground.Height),
			rl.Gray)
		rl.DrawText("User is not \n in unit", int32(w.infoUserScene.notInUnitBackground.X), int32(w.infoUserScene.notInUnitBackground.Y), 16, rl.White)

	}

}

func (w *Window) renderInfoUserState() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.infoUserScene.usersList.bounds) {
			w.infoUserScene.usersList.focus = 1
		}
		if rl.CheckCollisionPointRec(mousePos, w.infoUserScene.unitsToAssignSlider.bounds) {
			w.infoUserScene.usersList.focus = 0
			w.infoUserScene.unitsToAssignSlider.focus = 1
		}
		if w.infoUserScene.showAddModal || w.infoUserScene.showRemoveModal {
			w.infoUserScene.usersList.focus = 0
		}

	}
	gui.ListViewEx(w.infoUserScene.usersList.bounds, w.infoUserScene.usersList.strings, &w.infoUserScene.usersList.idxScroll, &w.infoUserScene.usersList.idxActiveElement, w.infoUserScene.usersList.focus)
	rl.DrawRectangle(int32(w.infoUserScene.descriptionBounds.X), int32(w.infoUserScene.descriptionBounds.Y), int32(w.infoUserScene.descriptionBounds.Width), int32(w.infoUserScene.descriptionBounds.Height), rl.White)
	rl.DrawText(
		w.infoUserScene.descriptionName+"\n"+
			w.infoUserScene.descriptionSurname+"\n"+
			w.infoUserScene.descriptionLVL+"\n", int32(w.infoUserScene.descriptionBounds.X), int32(w.infoUserScene.descriptionBounds.Y), 43, rl.Yellow)
	if w.infoUserScene.showAddModal {
		rl.DrawRectangle(
			int32(w.infoUserScene.addModal.background.X),
			int32(w.infoUserScene.addModal.background.Y),
			int32(w.infoUserScene.addModal.background.Width),
			int32(w.infoUserScene.addModal.background.Height),
			w.infoUserScene.addModal.bgColor)
		if gui.WindowBox(w.infoUserScene.addModal.core, "TITLE") {
			w.infoUserScene.showAddModal = false
			w.infoUserScene.unitsToAssignSlider.strings = w.infoUserScene.unitsToAssignSlider.strings[:0]
		}
		gui.ListViewEx(w.infoUserScene.unitsToAssignSlider.bounds,
			w.infoUserScene.unitsToAssignSlider.strings,
			&w.infoUserScene.unitsToAssignSlider.idxScroll,
			&w.infoUserScene.unitsToAssignSlider.idxActiveElement,
			w.infoUserScene.unitsToAssignSlider.focus)

		w.infoUserScene.isConfirmAddButtonPressed = gui.Button(w.infoUserScene.acceptAddButton.bounds, w.infoUserScene.acceptAddButton.text)
	}
	if w.infoUserScene.showRemoveModal {
		rl.DrawRectangle(
			int32(w.infoUserScene.removeModal.background.X),
			int32(w.infoUserScene.removeModal.background.Y),
			int32(w.infoUserScene.removeModal.background.Width),
			int32(w.infoUserScene.removeModal.background.Height),
			w.infoUserScene.removeModal.bgColor)
		if gui.WindowBox(w.infoUserScene.removeModal.core, "TITLE") {
			w.infoUserScene.showRemoveModal = false
			w.infoUserScene.usersUnitsSlider.strings = w.infoUserScene.usersUnitsSlider.strings[:0]
		}
		gui.ListViewEx(w.infoUserScene.usersUnitsSlider.bounds,
			w.infoUserScene.usersUnitsSlider.strings,
			&w.infoUserScene.usersUnitsSlider.idxScroll,
			&w.infoUserScene.usersUnitsSlider.idxActiveElement,
			w.infoUserScene.usersUnitsSlider.focus)

		w.infoUserScene.isConfirmRemoveButtonPressed = gui.Button(w.infoUserScene.acceptRemoveButton.bounds, w.infoUserScene.acceptRemoveButton.text)
	}
}
