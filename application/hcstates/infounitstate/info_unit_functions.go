package infounitstate

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (i *InfoUnitScene) Reset() {
	i.unitsSection.lastProcessedUnitIdx = -1
	i.usersSection.lastProcessedUserIdx = -1
}

// if we have -1 on unit -> block this big button
// if we have -1 on users -> block small buttons
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
}

func (i *InfoUnitScene) SelectUnit() {
	currentUnitIdx := i.unitsSection.unitsSlider.IdxActiveElement
	if currentUnitIdx != -1 && currentUnitIdx != i.unitsSection.lastProcessedUnitIdx {
		selectedUnit := i.unitsSection.units[currentUnitIdx]
		if _, ok := i.unitsSection.unitsInformation[selectedUnit.Id]; ok {
			i.unitsSection.currUnitID = selectedUnit.Id

			users := i.unitsSection.unitsInformation[selectedUnit.Id].Users
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
			if _, ok := i.unitsSection.unitsInformation[unit.Id]; ok {
				users := i.unitsSection.unitsInformation[unit.Id].Users
				if int(currentUserIdx) < len(users) {
					//user := users[currentUserIdx] ???
				}
			}

		}
		i.usersSection.lastProcessedUserIdx = currentUserIdx
	}
}

func (i *InfoUnitScene) UnitsDescription() {
	i.unitsSection.unitsInformation = make(map[string]*proto.UnitInformation)
	for _, u := range i.unitsSection.units {
		res, err := utils.MakeRequest(
			utils.NewRequest(
				i.cfg.Ctx,
				i.cfg.ServerPID,
				&proto.GetUnitInformation{
					UnitID: u.Id},
			))
		if err != nil {
			//TODO err section
		}
		if v, ok := res.(*proto.UnitInformation); ok {
			i.unitsSection.unitsInformation[v.UnitID] = v
		} else {
			//TODO err
		}
	}
}
func (i *InfoUnitScene) prepareDeviceSlider() {
	i.descriptionSection.devicesElements = make(map[string][]struct {
		bounds rl.Rectangle
		name   string
		desc   string
	})
	for id, unitInf := range i.unitsSection.unitsInformation {
		elements := make([]struct {
			bounds rl.Rectangle
			name   string
			desc   string
		}, 0, len(unitInf.Devices))

		for idx, d := range unitInf.Devices {
			element := struct {
				bounds rl.Rectangle
				name   string
				desc   string
			}{}

			element.bounds = rl.Rectangle{
				X:      10,
				Y:      float32(idx * 60), // wysokość elementu 60px
				Width:  i.descriptionSection.devicesSlider.View.Width - 20,
				Height: 50,
			}
			element.name = d.Name
			element.desc = "Owner: " + string(d.Owner) + "\n" + d.LastTimeOnline.AsTime().Format(time.DateTime)
			elements = append(elements, element)
		}

		i.descriptionSection.devicesElements[id] = elements
	}
}
func (i *InfoUnitScene) prepareTaskSlider() {
	var width int32 = 700
	var fontSize int32 = 14
	i.descriptionSection.tasksElements = make(map[string][]struct {
		bounds rl.Rectangle
		name   string
		desc   string
	})
	for id, unitInf := range i.unitsSection.unitsInformation {
		elements := make([]struct {
			bounds rl.Rectangle
			name   string
			desc   string
		}, 0, len(unitInf.Devices))

		for idx, t := range unitInf.Tasks {
			element := struct {
				bounds rl.Rectangle
				name   string
				desc   string
			}{}

			element.bounds = rl.Rectangle{
				X:      10,
				Y:      float32(idx * 60),
				Width:  i.descriptionSection.devicesSlider.View.Width - 20,
				Height: 50,
			}
			element.name = t.Name
			element.desc = utils.WrapText(width, t.Description, fontSize)
			elements = append(elements, element)
		}

		i.descriptionSection.tasksElements[id] = elements
	}
}
