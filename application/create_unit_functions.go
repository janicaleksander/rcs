package application

import (
	"time"

	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (w *Window) CreateUnit() {
	name := w.createUnitScene.newUnitSection.nameInput.GetText()
	user := w.createUnitScene.newUnitSection.usersDropdown.idxActiveElement
	if len(name) <= 0 || user <= 0 {
		w.createUnitScene.scheduler.After((3 * time.Second).Seconds(), func() {
			w.createUnitScene.errorSection.errorPopup.Hide()
		})
		w.createUnitScene.errorSection.errorPopup.Show()
		w.createUnitScene.errorSection.errorMessage = "Zero length error"
	} else {
		//user can be only in one unit in the same time -> error
		res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.CreateUnit{
			Name:         name,
			IsConfigured: false,
			UserID:       w.createUnitScene.newUnitSection.usersDropdown.strings[user],
		}))

		if err != nil {
			//context deadline exceeded
			w.createUnitScene.scheduler.After((3 * time.Second).Seconds(), func() {
				w.createUnitScene.errorSection.errorPopup.Hide()
			})
			w.createUnitScene.errorSection.errorPopup.Show()
			w.createUnitScene.errorSection.errorMessage = err.Error()
			return
		}
		if _, ok := res.(*proto.AcceptCreateUnit); ok {
			w.createUnitScene.scheduler.After((3 * time.Second).Seconds(), func() {
				w.createUnitScene.infoSection.infoPopup.Hide()
			})
			w.createUnitScene.infoSection.infoPopup.Show()
			w.createUnitScene.infoSection.infoMessage = "Success!"
		} else {
			w.createUnitScene.scheduler.After((3 * time.Second).Seconds(), func() {
				w.createUnitScene.errorSection.errorPopup.Hide()
			})
			w.createUnitScene.errorSection.errorPopup.Show()
			w.createUnitScene.errorSection.errorMessage = "Error!"

		}

	}

}
