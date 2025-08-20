package application

import (
	"sync"

	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (s *InfoUnitScene) Reset() {
	s.unitsSection.lastProcessedUnitIdx = -1
	s.usersSection.lastProcessedUserIdx = -1
}

func (w *Window) FetchUnits2() {

	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetAllUnits{}))
	if err != nil {
		// todo
	}

	if v, ok := res.(*proto.AllUnits); ok {
		w.infoUnitScene.unitsSection.units = v.Units
	} else {
		//todo error
	}

	w.infoUnitScene.unitsSection.unitsToUserCache = make(map[string][]*proto.User, len(w.infoUnitScene.unitsSection.units))
	var wg sync.WaitGroup
	var cacheChan = make(chan struct {
		id   string
		data []*proto.User
	}, len(w.infoUnitScene.unitsSection.units))

	for _, v := range w.infoUnitScene.unitsSection.units {
		wg.Add(1)
		go func(wGroup *sync.WaitGroup, id string) {
			defer wGroup.Done()

			res, err = utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetAllUsersInUnit{Id: id}))
			if err != nil {
				//todo error ctx deadline exceeded
			}

			if v, ok := res.(*proto.AllUsersInUnit); ok {
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
		w.infoUnitScene.unitsSection.unitsToUserCache[v.id] = v.data
	}

}

func (w *Window) SelectUnit() {
	currentUnitIdx := w.infoUnitScene.unitsSection.unitsSlider.idxActiveElement
	if currentUnitIdx != -1 && currentUnitIdx != w.infoUnitScene.unitsSection.lastProcessedUnitIdx {
		u := w.infoUnitScene.unitsSection.units[currentUnitIdx]
		users := w.infoUnitScene.unitsSection.unitsToUserCache[u.Id]
		w.infoUnitScene.usersSection.usersSlider.strings = make([]string, 0, len(users))
		for _, v := range users {
			w.infoUnitScene.usersSection.usersSlider.strings = append(w.infoUnitScene.usersSection.usersSlider.strings,
				v.Personal.Name+" "+v.Personal.Surname)
		}

		if len(users) > 0 {
			w.infoUnitScene.usersSection.usersSlider.idxActiveElement = 0
		} else {
			w.infoUnitScene.usersSection.usersSlider.idxActiveElement = -1
		}

		w.infoUnitScene.unitsSection.lastProcessedUnitIdx = currentUnitIdx
	}

}

func (w *Window) SelectUser() {
	currentUserIdx := w.infoUnitScene.usersSection.usersSlider.idxActiveElement
	if currentUserIdx != -1 && currentUserIdx != w.infoUnitScene.usersSection.lastProcessedUserIdx {
		selectedUnitIdx := w.infoUnitScene.unitsSection.unitsSlider.idxActiveElement

		if selectedUnitIdx != -1 {
			unit := w.infoUnitScene.unitsSection.units[selectedUnitIdx]
			users := w.infoUnitScene.unitsSection.unitsToUserCache[unit.Id]
			if int(currentUserIdx) < len(users) {
				//user := users[currentUserIdx] ???
			}
		}
		w.infoUnitScene.usersSection.lastProcessedUserIdx = currentUserIdx
	}
}
