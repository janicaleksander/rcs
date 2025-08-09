package application

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/proto"
	"sync"
	"time"
)

type InfoUnitScene struct {

	//TODO refactor	other down list with other [], divide to two slices

	//list of all units
	units []*proto.Unit

	usersCache  map[string][]*proto.User //preload users in each unit from units
	unitsSlider ListSlider

	//middle with general info about selected
	descriptionBounds  rl.Rectangle
	descriptionButton1 Button
	descriptionButton2 Button
	descriptionButton3 Button
	descriptionButton4 Button
	//... some elements about unit

	//last part with commander squad of selected unit (5,4,3 rank)
	usersSlider ListSlider

	//section with info about selected user
	userInfoBounds  rl.Rectangle
	userInfoButton1 Button //info
	userInfoButton2 Button //send message

	lastProcessedUnitIdx int32
	lastProcessedUserIdx int32

	//TODO add error's field
	//TODO add go back button
}

func (s *InfoUnitScene) Reset() {
	s.lastProcessedUnitIdx = -1
	s.lastProcessedUserIdx = -1
}

// TODO if slice is empty show some info about this (Nothing is here maybe)+maybe err field if error
func (w *Window) infoUnitSceneSetup() {
	w.infoUnitScene.Reset()
	resp := w.ctx.Request(w.serverPID, &proto.GetAllUnits{}, time.Second*5)
	val, err := resp.Result()
	if err != nil {
		// todo
	}
	if v, ok := val.(*proto.AllUnits); ok {
		w.infoUnitScene.units = v.Units
	}
	w.infoUnitScene.usersCache = make(map[string][]*proto.User, 0)
	var wg sync.WaitGroup
	var cacheChan = make(chan struct {
		id   string
		data []*proto.User
	}, len(w.infoUnitScene.units))

	for _, v := range w.infoUnitScene.units {
		wg.Add(1)
		go func(wGroup *sync.WaitGroup, id string) {
			defer wGroup.Done()
			resp = w.ctx.Request(w.serverPID, &proto.GetAllUsersInUnit{Id: id}, time.Second*5)
			val, err := resp.Result()
			if err != nil {
				// todo
			}
			if v, ok := val.(*proto.AllUsersInUnit); ok {
				cacheChan <- struct {
					id   string
					data []*proto.User
				}{id: id, data: v.Users}
			}

		}(&wg, v.Id)
	}
	go func() {
		wg.Wait()
		close(cacheChan)
	}()
	for v := range cacheChan {
		w.infoUnitScene.usersCache[v.id] = v.data
	}

	w.infoUnitScene.unitsSlider = ListSlider{
		strings: make([]string, 0),
		bounds: rl.NewRectangle(
			0,
			(1.0/8.0)*float32(w.height),
			(2.0/9.0)*float32(w.width),
			(7.0/8.0)*float32(w.height),
		),
		idxActiveElement: 0,
		focus:            1,
		idxScroll:        0,
	}
	for _, v := range w.infoUnitScene.units {
		w.infoUnitScene.unitsSlider.strings = append(w.infoUnitScene.unitsSlider.strings, v.Id[:5]+"..."+v.Id[31:])
	}
	w.infoUnitScene.descriptionBounds = rl.NewRectangle(
		w.infoUnitScene.unitsSlider.bounds.Width,
		w.infoUnitScene.unitsSlider.bounds.Y,
		(4.0/9.0)*float32(w.width),
		w.infoUnitScene.unitsSlider.bounds.Height)

	w.infoUnitScene.usersSlider = ListSlider{
		strings: make([]string, 0),
		bounds: rl.NewRectangle(
			w.infoUnitScene.unitsSlider.bounds.Width+w.infoUnitScene.descriptionBounds.Width,
			w.infoUnitScene.unitsSlider.bounds.Y,
			(4.0/12.0)*float32(w.width),
			(2.0/3.0)*float32(w.height)),
		idxActiveElement: 0,
		focus:            1,
		idxScroll:        0,
	}
	w.infoUnitScene.userInfoBounds = rl.NewRectangle(
		w.infoUnitScene.unitsSlider.bounds.Width+w.infoUnitScene.descriptionBounds.Width,
		(2.0/3.0)*float32(w.height),
		w.infoUnitScene.usersSlider.bounds.Width,
		(1.0/3.0)*float32(w.height))

	//TODO do full refactor of this conditions and check where have to have it
	if len(w.infoUnitScene.units) > 0 {
		w.infoUnitScene.unitsSlider.idxActiveElement = 0
	} else {
		w.infoUnitScene.usersSlider.idxActiveElement = -1
	}

	w.infoUnitScene.descriptionButton1 = Button{
		bounds: rl.NewRectangle(
			w.infoUnitScene.descriptionBounds.X,
			w.infoUnitScene.descriptionBounds.Y,
			w.infoUnitScene.descriptionBounds.Width/2,
			w.infoUnitScene.descriptionBounds.Height/2,
		),
		text: "Button Info \n 1",
	}
	w.infoUnitScene.descriptionButton2 = Button{
		bounds: rl.NewRectangle(
			w.infoUnitScene.descriptionBounds.X+w.infoUnitScene.descriptionButton1.bounds.Width,
			w.infoUnitScene.descriptionBounds.Y,
			w.infoUnitScene.descriptionBounds.Width/2,
			w.infoUnitScene.descriptionBounds.Height/2,
		),
		text: "Button Info \n 2",
	}
	w.infoUnitScene.descriptionButton3 = Button{
		bounds: rl.NewRectangle(
			w.infoUnitScene.descriptionBounds.X,
			w.infoUnitScene.descriptionBounds.Y+w.infoUnitScene.descriptionButton1.bounds.Height,
			w.infoUnitScene.descriptionBounds.Width/2,
			w.infoUnitScene.descriptionBounds.Height/2,
		),
		text: "Button Info \n 3",
	}
	w.infoUnitScene.descriptionButton4 = Button{
		bounds: rl.NewRectangle(
			w.infoUnitScene.descriptionBounds.X+w.infoUnitScene.descriptionButton1.bounds.Width,
			w.infoUnitScene.descriptionBounds.Y+w.infoUnitScene.descriptionButton1.bounds.Height,
			w.infoUnitScene.descriptionBounds.Width/2,
			w.infoUnitScene.descriptionBounds.Height/2,
		),
		text: "Button Info \n 4",
	}

	w.infoUnitScene.userInfoButton1 = Button{
		bounds: rl.NewRectangle(
			w.infoUnitScene.userInfoBounds.X,
			w.infoUnitScene.userInfoBounds.Y,
			w.infoUnitScene.userInfoBounds.Width,
			w.infoUnitScene.userInfoBounds.Height/2,
		),
		text: "User info \n 1",
	}
	w.infoUnitScene.userInfoButton2 = Button{
		bounds: rl.NewRectangle(
			w.infoUnitScene.userInfoBounds.X,
			w.infoUnitScene.userInfoBounds.Y+w.infoUnitScene.userInfoButton1.bounds.Height,
			w.infoUnitScene.userInfoBounds.Width,
			w.infoUnitScene.userInfoBounds.Height/2,
		),
		text: "User info \n 2",
	}

}

//TODO add some info where e.g. users are empty

func (w *Window) updateInfoUnitState() {
	currentUnitIdx := w.infoUnitScene.unitsSlider.idxActiveElement
	if currentUnitIdx != -1 && currentUnitIdx != w.infoUnitScene.lastProcessedUnitIdx {

		u := w.infoUnitScene.units[currentUnitIdx]
		users := w.infoUnitScene.usersCache[u.Id]
		w.infoUnitScene.usersSlider.strings = make([]string, 0, len(users))
		for _, v := range users {
			w.infoUnitScene.usersSlider.strings = append(w.infoUnitScene.usersSlider.strings, v.Personal.Name+" "+v.Personal.Surname)
		}

		if len(users) > 0 {
			w.infoUnitScene.usersSlider.idxActiveElement = 0
		} else {
			w.infoUnitScene.usersSlider.idxActiveElement = -1
		}

		w.infoUnitScene.Reset()
		w.infoUnitScene.lastProcessedUnitIdx = currentUnitIdx
	}

	currentUserIdx := w.infoUnitScene.usersSlider.idxActiveElement
	if currentUserIdx != -1 && currentUserIdx != w.infoUnitScene.lastProcessedUserIdx {
		selectedUnitIdx := w.infoUnitScene.unitsSlider.idxActiveElement

		if selectedUnitIdx != -1 {
			unit := w.infoUnitScene.units[selectedUnitIdx]
			users := w.infoUnitScene.usersCache[unit.Id]
			if int(currentUserIdx) < len(users) {
				//user := users[currentUserIdx]
			}
		}

		w.infoUnitScene.lastProcessedUserIdx = currentUserIdx
	}
}

func (w *Window) renderInfoUnitState() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.infoUnitScene.unitsSlider.bounds) {
			w.infoUnitScene.unitsSlider.focus = 1
			w.infoUnitScene.usersSlider.focus = 0
		}
		if rl.CheckCollisionPointRec(mousePos, w.infoUnitScene.usersSlider.bounds) {
			w.infoUnitScene.unitsSlider.focus = 0
			w.infoUnitScene.usersSlider.focus = 1
		}

	}
	//unitsSlider
	gui.ListViewEx(w.infoUnitScene.unitsSlider.bounds, w.infoUnitScene.unitsSlider.strings, &w.infoUnitScene.unitsSlider.idxScroll, &w.infoUnitScene.unitsSlider.idxActiveElement, w.infoUnitScene.unitsSlider.focus)
	//description box
	rl.DrawRectangle(int32(w.infoUnitScene.descriptionBounds.X),
		int32(w.infoUnitScene.descriptionBounds.Y),
		int32(w.infoUnitScene.descriptionBounds.Width),
		int32(w.infoUnitScene.descriptionBounds.Height),
		rl.Yellow)
	//users slider
	gui.ListViewEx(w.infoUnitScene.usersSlider.bounds, w.infoUnitScene.usersSlider.strings, &w.infoUnitScene.usersSlider.idxScroll, &w.infoUnitScene.usersSlider.idxActiveElement, w.infoUnitScene.usersSlider.focus)

	//user info box
	rl.DrawRectangle(int32(w.infoUnitScene.userInfoBounds.X),
		int32(w.infoUnitScene.userInfoBounds.Y),
		int32(w.infoUnitScene.userInfoBounds.Width),
		int32(w.infoUnitScene.userInfoBounds.Height),
		rl.White)

	gui.Button(w.infoUnitScene.descriptionButton1.bounds, w.infoUnitScene.descriptionButton1.text)
	gui.Button(w.infoUnitScene.descriptionButton2.bounds, w.infoUnitScene.descriptionButton2.text)
	gui.Button(w.infoUnitScene.descriptionButton3.bounds, w.infoUnitScene.descriptionButton3.text)
	gui.Button(w.infoUnitScene.descriptionButton4.bounds, w.infoUnitScene.descriptionButton4.text)

	gui.Button(w.infoUnitScene.userInfoButton1.bounds, w.infoUnitScene.userInfoButton1.text)
	gui.Button(w.infoUnitScene.userInfoButton2.bounds, w.infoUnitScene.userInfoButton2.text)

}
