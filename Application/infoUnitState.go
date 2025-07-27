package Application

import (
	"fmt"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"strconv"
	"sync"
	"time"
)

// TODO make components e.g slider component to not repeat this properties

type InfoUnitScene struct {
	//first left side slider with names of units 3/12

	//TODO refactor	other down list with other [], divide to two slices

	units []*Proto.Unit

	usersCache  map[string][]*Proto.User //preload users in each unit from units
	unitsSlider ListSlider

	//middle with general info about selected 5/12

	descriptionBounds rl.Rectangle
	//... some elements about unit

	//last part with commander squad of selected unit (5,4,3 rank)  4/12
	usersSlider ListSlider

	//section with info about selected user
	userInfoBounds  rl.Rectangle
	userInfoName    string
	userInfoSurname string
	userInfoLVL     string

	lastProcessedUnitIdx int32
	lastProcessedUserIdx int32
}

func (s *InfoUnitScene) Reset() {
	//TODO
	s.userInfoName = ""
	s.userInfoSurname = ""
	s.userInfoLVL = ""
	s.lastProcessedUnitIdx = -1
	s.lastProcessedUserIdx = -1
}

// TODO if slice is empty show some info about this (Nothing is here maybe)+maybe err field if error
func (w *Window) infoUnitSceneSetup() {
	w.infoUnitScene.Reset()
	resp := w.ctx.Request(w.serverPID, &Proto.GetAllUnits{}, time.Second*5)
	val, err := resp.Result()
	if err != nil {
		// todo
	}
	if v, ok := val.(*Proto.AllUnits); ok {
		w.infoUnitScene.units = v.Units
	}
	w.infoUnitScene.usersCache = make(map[string][]*Proto.User, 0)
	var wg sync.WaitGroup
	var cacheChan = make(chan struct {
		id   string
		data []*Proto.User
	}, len(w.infoUnitScene.units))

	for _, v := range w.infoUnitScene.units {
		wg.Add(1)
		go func(wGroup *sync.WaitGroup, id string) {
			defer wGroup.Done()
			resp = w.ctx.Request(w.serverPID, &Proto.GetAllUsersInUnit{Id: id}, time.Second*5)
			val, err := resp.Result()
			if err != nil {
				// todo
			}
			if v, ok := val.(*Proto.AllUsersInUnit); ok {
				cacheChan <- struct {
					id   string
					data []*Proto.User
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

	fmt.Println(w.infoUnitScene.usersCache)
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

	for _, v := range w.infoUnitScene.units {
		w.infoUnitScene.unitsSlider.strings = append(w.infoUnitScene.unitsSlider.strings, v.Id[:5]+"..."+v.Id[31:])
	}
	//TODO do full refactor of this conditions and check where have to have it
	if len(w.infoUnitScene.units) > 0 {
		w.infoUnitScene.unitsSlider.idxActiveElement = 0
	} else {
		w.infoUnitScene.usersSlider.idxActiveElement = -1
	}
}

//TODO add some info where e.g. users are empty

func (w *Window) updateInfoUnitState() {
	currentUnitIdx := w.infoUnitScene.unitsSlider.idxActiveElement // 0 0
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
		//w.infoUnitScene.lastProcessedUserIdx = -1
		w.infoUnitScene.lastProcessedUnitIdx = currentUnitIdx
	}

	currentUserIdx := w.infoUnitScene.usersSlider.idxActiveElement
	if currentUserIdx != -1 && currentUserIdx != w.infoUnitScene.lastProcessedUserIdx {
		selectedUnitIdx := w.infoUnitScene.unitsSlider.idxActiveElement

		if selectedUnitIdx != -1 {
			unit := w.infoUnitScene.units[selectedUnitIdx]
			users := w.infoUnitScene.usersCache[unit.Id]
			if int(currentUserIdx) < len(users) {
				user := users[currentUserIdx]
				w.infoUnitScene.userInfoName = user.Personal.Name
				w.infoUnitScene.userInfoSurname = user.Personal.Surname
				w.infoUnitScene.userInfoLVL = strconv.Itoa(int(user.RuleLvl))
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

	fontSize := int32(20)
	padding := float32(10)

	rl.DrawText(w.infoUnitScene.userInfoName,
		int32(w.infoUnitScene.userInfoBounds.X+padding),
		int32(w.infoUnitScene.userInfoBounds.Y+30),
		fontSize,
		rl.Black)

	rl.DrawText(w.infoUnitScene.userInfoSurname,
		int32(w.infoUnitScene.userInfoBounds.X+padding),
		int32(w.infoUnitScene.userInfoBounds.Y+30+30),
		fontSize,
		rl.Black)

	rl.DrawText(w.infoUnitScene.userInfoLVL,
		int32(w.infoUnitScene.userInfoBounds.X+padding),
		int32(w.infoUnitScene.userInfoBounds.Y+30+30+30),
		fontSize,
		rl.Black)

}
