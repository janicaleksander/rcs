package application

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *InfoUserScene) Reset() {
	i.userListSection.lastProcessedUserIdx = -1
	i.descriptionSection.descriptionName = ""
	i.descriptionSection.descriptionSurname = ""
	i.descriptionSection.descriptionLVL = ""
	i.userListSection.currSelectedUserID = ""
	i.userListSection.isInUnit = false
}

func (w *Window) UpdateDescription() {

	currentUserIdx := w.infoUserScene.userListSection.usersList.idxActiveElement
	if currentUserIdx != -1 && currentUserIdx != w.infoUserScene.userListSection.lastProcessedUserIdx {
		user := w.infoUserScene.userListSection.users[w.infoUserScene.userListSection.usersList.idxActiveElement]
		w.infoUserScene.userListSection.currSelectedUserID = user.Id
		//TODO in the v2 version we need to track more than
		// one unit ID
		if _, ok := w.infoUserScene.unitListSection.userToUnitCache[user.Id]; ok {
			w.infoUserScene.userListSection.isInUnit = true
		} else {
			w.infoUserScene.userListSection.isInUnit = false
		}
		//
		w.infoUserScene.descriptionSection.descriptionName = user.Personal.Name
		w.infoUserScene.descriptionSection.descriptionSurname = user.Personal.Surname
		w.infoUserScene.descriptionSection.descriptionLVL = strconv.Itoa(int(user.RuleLvl))
		w.infoUserScene.userListSection.lastProcessedUserIdx = currentUserIdx
	}

}

func (w *Window) AddToUnit() {
	if !w.infoUserScene.userListSection.isInUnit { // shows add to unit
		//fil userUnits ( TODO in v2 for loop through many units)
		if w.infoUserScene.actionSection.showAddModal {
			for _, unit := range w.infoUserScene.unitListSection.units {
				w.infoUserScene.addActionSection.unitsToAssignSlider.strings = append(w.infoUserScene.addActionSection.unitsToAssignSlider.strings, unit.Id)
				/*
					cacheUnit := w.infoUserScene.userToUnitCache[w.infoUserScene.currSelectedUserID]
					if cacheUnit == unit.Id {
						continue
					}
					in v2 version. we dont need to show units that we are already enrolled in
				*/

			}
			//w.infoUserScene.actionSection.showAddModal = true
		}

		if w.infoUserScene.addActionSection.isConfirmAddButtonPressed {
			if w.infoUserScene.addActionSection.unitsToAssignSlider.idxActiveElement >= 0 {
				unit := w.infoUserScene.unitListSection.units[w.infoUserScene.addActionSection.unitsToAssignSlider.idxActiveElement]

				res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.AssignUserToUnit{
					UserID: w.infoUserScene.userListSection.currSelectedUserID,
					UnitID: unit.Id,
				}))
				if err != nil {
					//context error deadline
				}
				if _, ok := res.(*proto.SuccessOfAssign); ok {
					// TODO failure
					w.infoUserScene.unitListSection.userToUnitCache[w.infoUserScene.userListSection.currSelectedUserID] = unit.Id
					w.infoUserScene.userListSection.isInUnit = true
				} else {
					//error
				}
			}

		}
	} else {
		rl.DrawRectangle(int32(w.infoUserScene.actionSection.inUnitBackground.X),
			int32(w.infoUserScene.actionSection.inUnitBackground.Y),
			int32(w.infoUserScene.actionSection.inUnitBackground.Width),
			int32(w.infoUserScene.actionSection.inUnitBackground.Height),
			rl.Gray)
		rl.DrawText(
			"User is \n in unit",
			int32(w.infoUserScene.actionSection.inUnitBackground.X),
			int32(w.infoUserScene.actionSection.inUnitBackground.Y),
			16,
			rl.White)

	}

}

func (w *Window) RemoveFromUnit() {
	if w.infoUserScene.userListSection.isInUnit { // shows remove  unit
		if w.infoUserScene.actionSection.showRemoveModal {
			w.infoUserScene.removeActionSection.usersUnitsSlider.strings =
				append(
					w.infoUserScene.removeActionSection.usersUnitsSlider.strings,
					w.infoUserScene.unitListSection.userToUnitCache[w.infoUserScene.userListSection.currSelectedUserID])
			//	w.infoUserScene.showRemoveModal = true
		}
		if w.infoUserScene.removeActionSection.isConfirmRemoveButtonPressed {
			if w.infoUserScene.removeActionSection.usersUnitsSlider.idxActiveElement >= 0 {
				unit := w.infoUserScene.unitListSection.units[w.infoUserScene.removeActionSection.usersUnitsSlider.idxActiveElement]

				res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.DeleteUserFromUnit{
					UserID: w.infoUserScene.userListSection.currSelectedUserID,
					UnitID: unit.Id,
				}))
				if err != nil {
					//error context deadline exceeded
				}

				if _, ok := res.(*proto.SuccessOfDelete); ok {
					//TODO success

					//TODO in v2 map str->[]str and then we have to iterate through
					// this slice and delete exact unit
					delete(w.infoUserScene.unitListSection.userToUnitCache, w.infoUserScene.userListSection.currSelectedUserID)
					w.infoUserScene.userListSection.isInUnit = false

				} else {
					//todo error
				}

			}

		}
	} else {
		rl.DrawRectangle(
			int32(w.infoUserScene.actionSection.notInUnitBackground.X),
			int32(w.infoUserScene.actionSection.notInUnitBackground.Y),
			int32(w.infoUserScene.actionSection.notInUnitBackground.Width),
			int32(w.infoUserScene.actionSection.notInUnitBackground.Height),
			rl.Gray)
		rl.DrawText(
			"User is not \n in unit",
			int32(w.infoUserScene.actionSection.notInUnitBackground.X),
			int32(w.infoUserScene.actionSection.notInUnitBackground.Y),
			16,
			rl.White)

	}

}

func (w *Window) SendMessage() {
	if w.infoUserScene.actionSection.showInboxModal {
		res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.IsOnline{
			Uuid: w.infoUserScene.userListSection.currSelectedUserID,
		}))

		if err != nil {
			//error context deadline exceeded
		}

		if _, ok := res.(*proto.Online); ok {
			w.infoUserScene.sendMessageSection.activeUserCircle.color = rl.Green
		} else {
			w.infoUserScene.sendMessageSection.activeUserCircle.color = rl.Red
		}
	}
	if w.infoUserScene.sendMessageSection.isSendMessageButtonPressed {
		message := w.infoUserScene.sendMessageSection.inboxInput.GetText()

		res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetLoggedInUUID{
			Pid: &proto.PID{
				Address: w.ctx.PID().Address,
				Id:      w.ctx.PID().ID}}))

		if err != nil {
			//error context deadline exceeded
		}

		var sender string
		if v, ok := res.(*proto.LoggedInUUID); !ok {
			//todo error return
		} else {
			sender = v.Id
		}

		res, err = utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.FillConversationID{
			SenderID:   sender,
			ReceiverID: w.infoUserScene.userListSection.currSelectedUserID,
		}))

		//TOOD finish the err handling sth like messenger type of send error some maybe red circle idk
		if err != nil {
			//ctx error
		}
		var cnvID string
		if v, ok := res.(*proto.SuccessOfFillConversationID); ok {
			cnvID = v.Id
		} else {
			//todo
			panic("ERROR CNV ID")
		}
		n := time.Now()

		res, err = utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.SendMessage{
			Receiver: w.infoUserScene.userListSection.currSelectedUserID,
			Message: &proto.Message{
				Id:             uuid.New().String(),
				SenderID:       sender,
				ConversationID: cnvID,
				Content:        message,
				SentAt:         timestamppb.Now(),
			}}))

		//TOOD finish the err handling sth like messenger type of send error some maybe red circle idk
		if err != nil {
			//todo error
			panic(err.Error())
		}

		if _, ok := res.(*proto.SuccessSend); !ok {
			//todo error
		}

		fmt.Println("CZAS SENDINGu", time.Since(n))
		w.infoUserScene.sendMessageSection.inboxInput.Clear()
	}

}

func (w *Window) FetchUnits() {

	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetAllUnits{}))
	if err != nil {
		// context deadline exceeded
	}

	w.infoUserScene.unitListSection.units = make([]*proto.Unit, 0, 64)

	if v, ok := res.(*proto.AllUnits); ok {
		for _, unit := range v.Units {
			w.infoUserScene.unitListSection.units = append(w.infoUserScene.unitListSection.units, unit)
		}
	} else {
		// TODO error
	}

}

func (w *Window) FetchUsers() {
	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetUserAboveLVL{Lvl: -1}))
	if err != nil {
		// context deadline exceeded
	}

	w.infoUserScene.userListSection.users = make([]*proto.User, 0, 64)
	if v, ok := res.(*proto.UsersAboveLVL); ok {
		for _, user := range v.Users {
			w.infoUserScene.userListSection.users = append(w.infoUserScene.userListSection.users, user)
		}
	} else {
		// TODO error
	}

	//cache users information
	w.infoUserScene.unitListSection.userToUnitCache = make(map[string]string, len(w.infoUserScene.userListSection.users))
	var waitGroup sync.WaitGroup
	cacheChan := make(chan struct {
		userID string
		unitID string
	}, 1024)

	for _, user := range w.infoUserScene.userListSection.users {
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup, userID string) {
			defer wg.Done()
			res, err = utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.IsUserInUnit{Id: userID}))
			if err != nil {
				//context deadline exceeded
			}
			if v, ok := res.(*proto.UserInUnit); ok {
				cacheChan <- struct {
					userID string
					unitID string
				}{userID: userID, unitID: v.UnitID}
			}

		}(&waitGroup, user.Id)
	}
	go func() {
		waitGroup.Wait()
		close(cacheChan)
	}()
	for v := range cacheChan {
		w.infoUserScene.unitListSection.userToUnitCache[v.userID] = v.unitID
	}
}
