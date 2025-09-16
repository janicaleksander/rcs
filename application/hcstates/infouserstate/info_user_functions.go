package infouserstate

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/application/component"
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
				if _, ok := res.(*proto.AcceptAssignUserToUnit); ok {
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

				if _, ok := res.(*proto.AcceptDeleteUserFromUnit); ok {
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
			UserID: i.userListSection.currSelectedUserID,
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
		if v, ok := res.(*proto.FilledConversationID); ok {
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

		if _, ok := res.(*proto.AcceptSend); !ok {
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
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetUserAboveLVL{
		Lower: -1,
		Upper: 10,
	})) //TODO
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
			if v, ok := res.(*proto.UserIsInUnit); ok {
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

func (i *InfoUserScene) prepareMap() {
	cfg := struct {
		General struct {
			Place string `toml:"place"`
		}
	}{}
	var startPoint = utils.STARTPLACE
	var centerLat = 0.0
	var centerLon = 0.0
	_, err := toml.DecodeFile("configproduction/general.toml", &cfg)
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		startPoint = cfg.General.Place
	}
	switch startPoint {
	case "WRO":
		centerLat = 51.11080123267171
		centerLon = 17.018041879680265
	}

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

func (i *InfoUserScene) FetchPins() {
	res, err := utils.MakeRequest(utils.NewRequest(
		i.cfg.Ctx,
		i.cfg.ServerPID,
		&proto.GetPins{}))
	if err != nil {
		//TODO
	}
	if v, ok := res.(*proto.Pins); !ok {
		//error
		return
	} else {
		//map with pins information
		for _, p := range v.Pins {
			x, y := latLonToPixel(p.Location.Latitude, p.Location.Longitude, ZOOM)
			i.trackUserLocationSection.locationMapInformation.MapPinInformation[p.DeviceID] = &component.PinInformation{
				Position:       rl.Vector2{X: x, Y: y},
				DeviceID:       p.DeviceID,
				OwnerName:      p.OwnerName,
				OwnerSurname:   p.OwnerSurname,
				LastTimeOnline: p.LastOnline.AsTime(),
			}
		}
		//map with last task tab information
		var waitGroup sync.WaitGroup
		for _, p := range v.Pins {
			go func(wg *sync.WaitGroup) {
				wg.Add(1)
				defer wg.Done()
				res, err = utils.MakeRequest(
					utils.NewRequest(
						i.cfg.Ctx,
						i.cfg.ServerPID,
						&proto.GetCurrentTask{DeviceID: p.DeviceID}))
				if err != nil {
					//todo
				} else {
					if v, ok := res.(*proto.CurrentTask); ok {
						i.trackUserLocationSection.locationMapInformation.MapCurrentTask[p.DeviceID] = &component.CurrentTaskTab{
							OwnerID:        v.UserID,
							OwnerName:      p.OwnerName,
							OwnerSurname:   p.OwnerSurname,
							DeviceID:       p.DeviceID,
							LastTimeOnline: p.LastOnline.AsTime(),
							Task:           v.Task,
						}
					}
				}
			}(&waitGroup)

		}
		waitGroup.Wait()
	}

}

func (i *InfoUserScene) drawPins() {
	for _, p := range i.trackUserLocationSection.locationMapInformation.MapPinInformation {
		rl.DrawTexture(i.trackUserLocationSection.LocationMap.pinTexture, int32(p.Position.X), int32(p.Position.Y), rl.White)
	}

}
func (i *InfoUserScene) showPinInformationOnCollision(mousePos rl.Vector2) {
	for _, p := range i.trackUserLocationSection.locationMapInformation.MapPinInformation {
		if checkMousePinCollision(p.Position, mousePos) {
			drawInfoBox(p)
		}
	}

}
func (i *InfoUserScene) showTabInformationOnCollision(mousePos rl.Vector2) {
	for _, p := range i.trackUserLocationSection.locationMapInformation.MapPinInformation {
		if checkMousePinCollision(p.Position, mousePos) {
			if _, ok := i.trackUserLocationSection.locationMapInformation.MapCurrentTask[p.DeviceID]; ok {
				i.drawInfoTab(i.trackUserLocationSection.locationMapInformation.MapCurrentTask[p.DeviceID])
			}
		}
	}

}

func checkMousePinCollision(pinPos, mousePos rl.Vector2) bool {
	pinBox := rl.NewRectangle(
		pinPos.X,
		pinPos.Y,
		32,
		32,
	)
	return rl.CheckCollisionPointRec(mousePos, pinBox)

}
func (i *InfoUserScene) drawInfoTab(currentTaskTab *component.CurrentTaskTab) {
	//upper box - users info
	rl.DrawRectangle(
		int32(i.trackUserLocationSection.userInfoTab.X),
		int32(i.trackUserLocationSection.userInfoTab.Y),
		int32(i.trackUserLocationSection.userInfoTab.Width),
		int32(i.trackUserLocationSection.userInfoTab.Height),
		rl.NewColor(250, 250, 250, 255))
	height := int32(i.trackUserLocationSection.userInfoTab.Y)
	// Owner ID
	rl.DrawText(
		"Owner ID:",
		int32(i.trackUserLocationSection.userInfoTab.X),
		height,
		16,
		rl.Gray)
	rl.DrawText(
		currentTaskTab.OwnerID,
		int32(i.trackUserLocationSection.userInfoTab.X)+200,
		height,
		16,
		rl.Black)

	// Last Online
	rl.DrawText(
		"Last Online:",
		int32(i.trackUserLocationSection.userInfoTab.X),
		height+25,
		16,
		rl.Gray)
	rl.DrawText(
		currentTaskTab.LastTimeOnline.Format("2006.01.02 -*- 15:04"),
		int32(i.trackUserLocationSection.userInfoTab.X)+200,
		height+25,
		16,
		rl.Black)

	// Device ID
	rl.DrawText(
		"Device ID:",
		int32(i.trackUserLocationSection.userInfoTab.X),
		height+50,
		16,
		rl.Gray)
	rl.DrawText(
		currentTaskTab.DeviceID,
		int32(i.trackUserLocationSection.userInfoTab.X)+200,
		height+50,
		16,
		rl.Black)

	// Owner Name
	rl.DrawText(
		"Owner Name:",
		int32(i.trackUserLocationSection.userInfoTab.X),
		height+75,
		16,
		rl.Gray)
	rl.DrawText(
		currentTaskTab.OwnerName,
		int32(i.trackUserLocationSection.userInfoTab.X)+200,
		height+75,
		16,
		rl.Black)

	// Owner Surname
	rl.DrawText(
		"Owner Surname:",
		int32(i.trackUserLocationSection.userInfoTab.X),
		height+100,
		16,
		rl.Gray)
	rl.DrawText(
		currentTaskTab.OwnerSurname,
		int32(i.trackUserLocationSection.userInfoTab.X)+200,
		height+100,
		16,
		rl.Black)

	//lower box - current task info
	rl.DrawRectangle(
		int32(i.trackUserLocationSection.currentTaskTab.X),
		int32(i.trackUserLocationSection.currentTaskTab.Y),
		int32(i.trackUserLocationSection.currentTaskTab.Width),
		int32(i.trackUserLocationSection.currentTaskTab.Height),
		rl.NewColor(250, 250, 250, 255))

	//TODO repair what if I dont have any current task
	text := utils.WrapText(
		int32(i.trackUserLocationSection.currentTaskTab.Width),
		currentTaskTab.Task.Description,
		15)

	//name
	rl.DrawText(
		"TASK NAME: ",
		int32(i.trackUserLocationSection.currentTaskTab.X),
		int32(i.trackUserLocationSection.currentTaskTab.Y+5),
		15,
		rl.LightGray)
	rl.DrawText(
		currentTaskTab.Task.Name,
		int32(i.trackUserLocationSection.currentTaskTab.X)+100,
		int32(i.trackUserLocationSection.currentTaskTab.Y+5),
		15,
		rl.Black)

	//TODO max 500chars here
	//desc
	rl.DrawText(
		text,
		int32(i.trackUserLocationSection.currentTaskTab.X),
		int32(i.trackUserLocationSection.currentTaskTab.Y+35),
		15,
		rl.Black)
}

func drawInfoBox(pin *component.PinInformation) {
	rl.SetMouseCursor(rl.MouseCursorPointingHand)
	notificationBox := rl.NewRectangle(
		pin.Position.X-64,
		pin.Position.Y-64,
		350,
		64)
	text := pin.OwnerName + "\n" + pin.OwnerSurname
	textWidth := rl.MeasureText(text, 25)
	rl.DrawRectangle(int32(notificationBox.X), int32(notificationBox.Y), int32(notificationBox.Width), int32(notificationBox.Height), rl.White)
	x := int32(float32(notificationBox.X) + float32(notificationBox.Width)/2 - float32(textWidth)/2)

	rl.DrawText(
		text,
		x,
		int32(notificationBox.Y),
		25,
		rl.Black)
}

func (i *InfoUserScene) drawMap() rl.Vector2 {
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

	if !i.trackUserLocationSection.LocationMap.isPinLoaded {
		i.trackUserLocationSection.LocationMap.pinTexture = rl.LoadTexture("osm/output.png")
		i.trackUserLocationSection.LocationMap.isPinLoaded = true
	}
	i.drawPins()
	scale := float32(rl.GetRenderWidth()) / float32(rl.GetScreenWidth())
	mouse := rl.GetMousePosition()
	mouse.X *= scale
	mouse.Y *= scale
	mousePos := rl.GetScreenToWorld2D(mouse, i.trackUserLocationSection.LocationMap.camera)
	i.showPinInformationOnCollision(mousePos)
	rl.EndMode2D()
	return mousePos
}
