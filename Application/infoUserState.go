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
	userToUnitCache map[string]bool
	// description area
	descriptionBounds  rl.Rectangle
	descriptionName    string
	descriptionSurname string
	descriptionLVL     string
	userIsInUnit       bool
	currUserID         string

	// action button area
	actionButtonArea rl.Rectangle
	// add btn
	addButton        Button
	inUnitBackground rl.Rectangle
	// rmv btn
	removeButton        Button
	notInUnitBackground rl.Rectangle
	// inbox btn
	inboxButton Button

	lastProcessedUserIdx int32
}

func (i *InfoUserScene) Reset() {
	i.lastProcessedUserIdx = -1
	i.descriptionName = ""
	i.descriptionSurname = ""
	i.descriptionLVL = ""
}
func (w *Window) InfoUserSceneSetup() {
	w.infoUserScene.Reset()
	//TODO get proper lvl
	resp := w.ctx.Request(w.serverPID, &Proto.GetUserAboveLVL{Lvl: -1}, WaitTime)
	val, err := resp.Result()
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
	w.infoUserScene.userToUnitCache = make(map[string]bool, len(w.infoUserScene.users))
	var waitGroup sync.WaitGroup
	cacheChan := make(chan struct {
		id       string
		isInUnit bool
	}, 1024)
	var isInUnit bool
	for _, user := range w.infoUserScene.users {
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup, id string) {
			defer wg.Done()
			resp = w.ctx.Request(w.serverPID, &Proto.IsUserInUnit{Id: id}, WaitTime)
			v, err := resp.Result()
			if err != nil {
				// TODO
				fmt.Println("ERROR")
			}
			if _, ok := v.(*Proto.UserInUnit); ok {
				isInUnit = true
			}
			if _, ok := v.(*Proto.UserNotInUnit); ok {
				isInUnit = false
			}
			cacheChan <- struct {
				id       string
				isInUnit bool
			}{id: id, isInUnit: isInUnit}
		}(&waitGroup, user.Id)
	}
	go func() {
		waitGroup.Wait()
		close(cacheChan)
	}()
	for v := range cacheChan {
		w.infoUserScene.userToUnitCache[v.id] = v.isInUnit
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
}

func (w *Window) updateInfoUserState() {
	currentUserIdx := w.infoUserScene.usersList.idxActiveElement
	if currentUserIdx != -1 && currentUserIdx != w.infoUserScene.lastProcessedUserIdx {
		//
		user := w.infoUserScene.users[w.infoUserScene.usersList.idxActiveElement]
		w.infoUserScene.currUserID = user.Id
		w.infoUserScene.userIsInUnit = w.infoUserScene.userToUnitCache[user.Id]
		//
		w.infoUserScene.descriptionName = user.Personal.Name
		w.infoUserScene.descriptionSurname = user.Personal.Surname
		w.infoUserScene.descriptionLVL = strconv.Itoa(int(user.RuleLvl))
		w.infoUserScene.lastProcessedUserIdx = currentUserIdx
	}

	if !w.infoUserScene.userIsInUnit {
		if gui.Button(w.infoUserScene.addButton.bounds, w.infoUserScene.addButton.text) {
			//button enable is user is not in unit
			//get vars
			// if yes  error
			// if no success

			//
			//
			w.infoUserScene.userToUnitCache[w.infoUserScene.currUserID] = true
			//db call

		}
	} else {
		rl.DrawRectangle(int32(w.infoUserScene.inUnitBackground.X),
			int32(w.infoUserScene.inUnitBackground.Y),
			int32(w.infoUserScene.inUnitBackground.Width),
			int32(w.infoUserScene.inUnitBackground.Height),
			rl.Gray)
		rl.DrawText("User is \n in unit", int32(w.infoUserScene.inUnitBackground.X), int32(w.infoUserScene.inUnitBackground.Y), 16, rl.White)

	}
	if w.infoUserScene.userIsInUnit {
		if gui.Button(w.infoUserScene.removeButton.bounds, w.infoUserScene.removeButton.text) {
			// button enable is user is in unit

			//
			//
			w.infoUserScene.userToUnitCache[w.infoUserScene.currUserID] = false
			w.infoUserScene.userIsInUnit = false
			//db call
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

	}
	gui.ListViewEx(w.infoUserScene.usersList.bounds, w.infoUserScene.usersList.strings, &w.infoUserScene.usersList.idxScroll, &w.infoUserScene.usersList.idxActiveElement, w.infoUserScene.usersList.focus)
	rl.DrawRectangle(int32(w.infoUserScene.descriptionBounds.X), int32(w.infoUserScene.descriptionBounds.Y), int32(w.infoUserScene.descriptionBounds.Width), int32(w.infoUserScene.descriptionBounds.Height), rl.White)
	rl.DrawText(
		w.infoUserScene.descriptionName+"\n"+
			w.infoUserScene.descriptionSurname+"\n"+
			w.infoUserScene.descriptionLVL+"\n", int32(w.infoUserScene.descriptionBounds.X), int32(w.infoUserScene.descriptionBounds.Y), 43, rl.Yellow)
}
