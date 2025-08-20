package infounitstate

import (
	"sync"

	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (i *InfoUnitScene) Reset() {
	i.unitsSection.lastProcessedUnitIdx = -1
	i.usersSection.lastProcessedUserIdx = -1
}

func (i *InfoUnitScene) FetchUnits() {

	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetAllUnits{}))
	if err != nil {
		// todo
	}

	if v, ok := res.(*proto.AllUnits); ok {
		i.unitsSection.units = v.Units
	} else {
		//todo error
	}

	i.unitsSection.unitsToUserCache = make(map[string][]*proto.User, len(i.unitsSection.units))
	var wg sync.WaitGroup
	var cacheChan = make(chan struct {
		id   string
		data []*proto.User
	}, len(i.unitsSection.units))

	for _, v := range i.unitsSection.units {
		wg.Add(1)
		go func(wGroup *sync.WaitGroup, id string) {
			defer wGroup.Done()

			res, err = utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetAllUsersInUnit{Id: id}))
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
		i.unitsSection.unitsToUserCache[v.id] = v.data
	}

}

func (i *InfoUnitScene) SelectUnit() {
	currentUnitIdx := i.unitsSection.unitsSlider.IdxActiveElement
	if currentUnitIdx != -1 && currentUnitIdx != i.unitsSection.lastProcessedUnitIdx {
		u := i.unitsSection.units[currentUnitIdx]
		users := i.unitsSection.unitsToUserCache[u.Id]
		i.usersSection.usersSlider.Strings = make([]string, 0, len(users))
		for _, v := range users {
			i.usersSection.usersSlider.Strings = append(i.usersSection.usersSlider.Strings,
				v.Personal.Name+" "+v.Personal.Surname)
		}

		if len(users) > 0 {
			i.usersSection.usersSlider.IdxActiveElement = 0
		} else {
			i.usersSection.usersSlider.IdxActiveElement = -1
		}

		i.unitsSection.lastProcessedUnitIdx = currentUnitIdx
	}

}

func (i *InfoUnitScene) SelectUser() {
	currentUserIdx := i.usersSection.usersSlider.IdxActiveElement
	if currentUserIdx != -1 && currentUserIdx != i.usersSection.lastProcessedUserIdx {
		selectedUnitIdx := i.unitsSection.unitsSlider.IdxActiveElement

		if selectedUnitIdx != -1 {
			unit := i.unitsSection.units[selectedUnitIdx]
			users := i.unitsSection.unitsToUserCache[unit.Id]
			if int(currentUserIdx) < len(users) {
				//user := users[currentUserIdx] ???
			}
		}
		i.usersSection.lastProcessedUserIdx = currentUserIdx
	}
}
