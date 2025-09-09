package infouserstate

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/types/proto"
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

func (i *InfoUserScene) UpdateDescription() {

	currentUserIdx := i.userListSection.usersList.IdxActiveElement
	if currentUserIdx != -1 && currentUserIdx != i.userListSection.lastProcessedUserIdx {
		user := i.userListSection.users[i.userListSection.usersList.IdxActiveElement]
		i.userListSection.currSelectedUserID = user.Id
		//TODO in the v2 version we need to track more than
		// one unit ID
		if _, ok := i.unitListSection.userToUnitCache[user.Id]; ok {
			i.userListSection.isInUnit = true
		} else {
			i.userListSection.isInUnit = false
		}
		//
		i.descriptionSection.descriptionName = user.Personal.Name
		i.descriptionSection.descriptionSurname = user.Personal.Surname
		i.descriptionSection.descriptionLVL = strconv.Itoa(int(user.RuleLvl))
		i.userListSection.lastProcessedUserIdx = currentUserIdx
	}

}

func (i *InfoUserScene) AddToUnit() {
	if !i.userListSection.isInUnit { // shows add to unit
		//fil userUnits ( TODO in v2 for loop through many units)
		if i.actionSection.showAddModal {
			for _, unit := range i.unitListSection.units {
				i.addActionSection.unitsToAssignSlider.Strings = append(i.addActionSection.unitsToAssignSlider.Strings, unit.Id)
				/*
					cacheUnit := i.userToUnitCache[i.currSelectedUserID]
					if cacheUnit == unit.Id {
						continue
					}
					in v2 version. we dont need to show units that we are already enrolled in
				*/

			}
			//i.actionSection.showAddModal = true
		}

		if i.addActionSection.isConfirmAddButtonPressed {
			if i.addActionSection.unitsToAssignSlider.IdxActiveElement >= 0 {
				unit := i.unitListSection.units[i.addActionSection.unitsToAssignSlider.IdxActiveElement]

				res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.AssignUserToUnit{
					UserID: i.userListSection.currSelectedUserID,
					UnitID: unit.Id,
				}))
				if err != nil {
					//context error deadline
				}
				if _, ok := res.(*proto.SuccessOfAssign); ok {
					// TODO failure
					i.unitListSection.userToUnitCache[i.userListSection.currSelectedUserID] = unit.Id
					i.userListSection.isInUnit = true
				} else {
					//error
				}
			}

		}
	} else {
		rl.DrawRectangle(int32(i.actionSection.inUnitBackground.X),
			int32(i.actionSection.inUnitBackground.Y),
			int32(i.actionSection.inUnitBackground.Width),
			int32(i.actionSection.inUnitBackground.Height),
			rl.Gray)
		rl.DrawText(
			"User is \n in unit",
			int32(i.actionSection.inUnitBackground.X),
			int32(i.actionSection.inUnitBackground.Y),
			16,
			rl.White)

	}

}

func (i *InfoUserScene) RemoveFromUnit() {
	if i.userListSection.isInUnit { // shows remove  unit
		if i.actionSection.showRemoveModal {
			i.removeActionSection.usersUnitsSlider.Strings =
				append(
					i.removeActionSection.usersUnitsSlider.Strings,
					i.unitListSection.userToUnitCache[i.userListSection.currSelectedUserID])
			//	i.showRemoveModal = true
		}
		if i.removeActionSection.isConfirmRemoveButtonPressed {
			if i.removeActionSection.usersUnitsSlider.IdxActiveElement >= 0 {
				unit := i.unitListSection.units[i.removeActionSection.usersUnitsSlider.IdxActiveElement]

				res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.DeleteUserFromUnit{
					UserID: i.userListSection.currSelectedUserID,
					UnitID: unit.Id,
				}))
				if err != nil {
					//error context deadline exceeded
				}

				if _, ok := res.(*proto.SuccessOfDelete); ok {
					//TODO success

					//TODO in v2 map str->[]str and then we have to iterate through
					// this slice and delete exact unit
					delete(i.unitListSection.userToUnitCache, i.userListSection.currSelectedUserID)
					i.userListSection.isInUnit = false

				} else {
					//todo error
				}

			}

		}
	} else {
		rl.DrawRectangle(
			int32(i.actionSection.notInUnitBackground.X),
			int32(i.actionSection.notInUnitBackground.Y),
			int32(i.actionSection.notInUnitBackground.Width),
			int32(i.actionSection.notInUnitBackground.Height),
			rl.Gray)
		rl.DrawText(
			"User is not \n in unit",
			int32(i.actionSection.notInUnitBackground.X),
			int32(i.actionSection.notInUnitBackground.Y),
			16,
			rl.White)

	}

}

func (i *InfoUserScene) SendMessage() {
	if i.actionSection.showInboxModal {
		res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.IsOnline{
			Uuid: i.userListSection.currSelectedUserID,
		}))

		if err != nil {
			//error context deadline exceeded
		}

		if _, ok := res.(*proto.Online); ok {
			i.sendMessageSection.activeUserCircle.Color = rl.Green
		} else {
			i.sendMessageSection.activeUserCircle.Color = rl.Red
		}
	}
	if i.sendMessageSection.isSendMessageButtonPressed {
		message := i.sendMessageSection.inboxInput.GetText()

		res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetLoggedInUUID{
			Pid: &proto.PID{
				Address: i.cfg.Ctx.PID().Address,
				Id:      i.cfg.Ctx.PID().ID}}))

		if err != nil {
			//error context deadline exceeded
		}

		var sender string
		if v, ok := res.(*proto.LoggedInUUID); !ok {
			//todo error return
		} else {
			sender = v.Id
		}

		res, err = utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.FillConversationID{
			SenderID:   sender,
			ReceiverID: i.userListSection.currSelectedUserID,
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

		res, err = utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.SendMessage{
			Receiver: i.userListSection.currSelectedUserID,
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
		i.sendMessageSection.inboxInput.Clear()
	}

}

func (i *InfoUserScene) FetchUnits() {

	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetAllUnits{}))
	if err != nil {
		// context deadline exceeded
	}

	i.unitListSection.units = make([]*proto.Unit, 0, 64)

	if v, ok := res.(*proto.AllUnits); ok {
		for _, unit := range v.Units {
			i.unitListSection.units = append(i.unitListSection.units, unit)
		}
	} else {
		// TODO error
	}

}

func (i *InfoUserScene) FetchUsers() {
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetUserAboveLVL{Lvl: -1}))
	if err != nil {
		// context deadline exceeded
	}

	i.userListSection.users = make([]*proto.User, 0, 64)
	if v, ok := res.(*proto.UsersAboveLVL); ok {
		for _, user := range v.Users {
			i.userListSection.users = append(i.userListSection.users, user)
		}
	} else {
		// TODO error
	}

	//cache users information
	i.unitListSection.userToUnitCache = make(map[string]string, len(i.userListSection.users))
	var waitGroup sync.WaitGroup
	cacheChan := make(chan struct {
		userID string
		unitID string
	}, 1024)

	for _, user := range i.userListSection.users {
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup, userID string) {
			defer wg.Done()
			res, err = utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.IsUserInUnit{Id: userID}))
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
		i.unitListSection.userToUnitCache[v.userID] = v.unitID
	}
}

// TODO add start point
// this is going to fetch from some kind of cfg when we have to set a default city (now is only WRO)
func (i *InfoUserScene) prepareMap() {
	centerLat := 51.11080123267171
	centerLon := 17.018041879680265
	mapX, mapY := latLonToPixel(centerLat, centerLon, ZOOM)
	i.trackUserLocationSection.LocationMap.camera = rl.Camera2D{
		Offset: rl.Vector2{
			X: i.trackUserLocationSection.mapModal.Core.Width/2 + (i.trackUserLocationSection.LocationMap.width)/2,
			Y: i.trackUserLocationSection.mapModal.Core.Height/2 + (i.trackUserLocationSection.LocationMap.height)/2,
		},
		Target: rl.Vector2{
			X: mapX,
			Y: mapY,
		},
		Rotation: 0,
		Zoom:     1,
	}
	i.trackUserLocationSection.LocationMap.tm.preloadNearbyTiles(
		mapX,
		mapY,
	)

}
func (i *InfoUserScene) updateMap() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		i.trackUserLocationSection.LocationMap.isDraggingCamera = true
	}
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		i.trackUserLocationSection.LocationMap.isDraggingCamera = false
	}
	if i.trackUserLocationSection.LocationMap.isDraggingCamera {
		delta := rl.GetMouseDelta()
		i.trackUserLocationSection.LocationMap.camera.Target.X -= delta.X
		i.trackUserLocationSection.LocationMap.camera.Target.Y -= delta.Y
	}

	select {
	case tile := <-i.trackUserLocationSection.LocationMap.tm.tileQueue:
		tile.loadTextureNow()
	default:
	}

	i.trackUserLocationSection.LocationMap.tm.setVisibleTiles(
		i.trackUserLocationSection.LocationMap.camera.Target.X,
		i.trackUserLocationSection.LocationMap.camera.Target.Y,
		int(i.trackUserLocationSection.LocationMap.width),
		int(i.trackUserLocationSection.LocationMap.height))
	i.trackUserLocationSection.LocationMap.tm.preloadNearbyTiles(
		i.trackUserLocationSection.LocationMap.camera.Target.X,
		i.trackUserLocationSection.LocationMap.camera.Target.Y,
	)
	i.trackUserLocationSection.LocationMap.tm.cleanupDistantTiles(
		i.trackUserLocationSection.LocationMap.camera.Target.X,
		i.trackUserLocationSection.LocationMap.camera.Target.Y,
	)
	//TODO
	/*
		tm.mu.Lock()
		defer tm.mu.Unlock()
		for _, tile := range tm.tiles {
			tile.unload()
		}
	*/

}
func (i *InfoUserScene) drawMap() {
	rl.BeginMode2D(i.trackUserLocationSection.LocationMap.camera)
	tiles := i.trackUserLocationSection.LocationMap.tm.getLoadedTiles()
	for _, tile := range tiles {
		if tile.isReady() {
			rl.DrawTexture(tile.getTexture(),
				int32(tile.x*TILESIZE),
				int32(tile.y*TILESIZE),
				rl.White)
		}
	}

	//drawPins()
	if !i.trackUserLocationSection.LocationMap.isPinLoaded {
		i.trackUserLocationSection.LocationMap.pinTexture = rl.LoadTexture("osm/output.png")
		i.trackUserLocationSection.LocationMap.isPinLoaded = true
	}

	scale := float32(rl.GetRenderWidth()) / float32(rl.GetScreenWidth())
	mouse := rl.GetMousePosition()
	mouse.X *= scale
	mouse.Y *= scale
	//	mouseWorldPos := rl.GetScreenToWorld2D(mouse, i.trackUserLocationSection.LocationMap.camera)
	rl.EndMode2D()
}

func isOnPin(mousePos, pinPos rl.Vector2) bool {
	pinBox := rl.NewRectangle(
		pinPos.X,
		pinPos.Y,
		64,
		64,
	)

	isColliding := rl.CheckCollisionPointRec(mousePos, pinBox)
	if isColliding {
		rl.SetMouseCursor(rl.MouseCursorPointingHand)
		notificationBox := rl.NewRectangle(
			pinPos.X-64,
			pinPos.Y-64,
			200,
			64)

		rl.DrawRectangle(
			int32(notificationBox.X),
			int32(notificationBox.Y),
			int32(notificationBox.Width),
			int32(notificationBox.Height),
			rl.Blue)

		rl.DrawText(
			"Pin Location",
			int32(notificationBox.X+5),
			int32(notificationBox.Y+20),
			20,
			rl.White)
	} else {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}
	return isColliding
}
