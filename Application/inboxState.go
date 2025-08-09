package Application

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Utils"
)

type InboxScene struct {
	toolboxArea      rl.Rectangle
	conversationArea rl.Rectangle
	messagesArea     rl.Rectangle

	usersConversation []*Proto.ConversationSummary
}

func (i *InboxScene) Reset() {

}

// TODO maybe use redis for fast cache to e.g user UUID
// TODO in refactor change all names to pattern verb+scene
func (w *Window) setupInboxScene() {
	//get users conversations
	res := w.ctx.Request(w.serverPID, &Proto.GetLoggedInUUID{
		Pid: &Proto.PID{
			Address: w.ctx.PID().Address,
			Id:      w.ctx.PID().ID}}, Utils.WaitTime)
	resp, err := res.Result()
	if err != nil {
		//TODO STH
	}
	v, ok := resp.(*Proto.LoggedInUUID)
	if !ok {
		//TODO
	}
	fmt.Println("XD", v.Id)
	sender := v.Id

	//TODO in the future make own databae conected to msssvc and send everything through this server e.g this and store message
	res = w.ctx.Request(w.serverPID, &Proto.GetUserConversation{Id: sender}, Utils.WaitTime)
	resp, err = res.Result()
	_, ok = resp.(*Proto.FailureGetUserConversation)
	if err != nil || !ok {
		//TODO error
		fmt.Print(err)
	}
	if conversations, ok := resp.(*Proto.SuccessGetUserConversation); ok {
		w.inboxScene.usersConversation = conversations.ConvSummary
	}

	fmt.Println(w.inboxScene.usersConversation)
	w.inboxScene.toolboxArea = rl.NewRectangle(
		0,
		0,
		(2.0/5.0)*float32(w.width),
		(1.0/8.0)*float32(w.height))
	w.inboxScene.conversationArea = rl.NewRectangle(
		w.inboxScene.toolboxArea.X,
		w.inboxScene.toolboxArea.Height,
		w.inboxScene.toolboxArea.Width,
		(7.0/8.0)*float32(w.height))
	w.inboxScene.messagesArea = rl.NewRectangle(
		w.inboxScene.toolboxArea.Width,
		w.inboxScene.toolboxArea.Y,
		(3.0/5.0)*float32(w.width),
		float32(w.height))
}
func (w *Window) updateInboxState() {

}
func (w *Window) renderInboxState() {

	rl.DrawRectangle(
		int32(w.inboxScene.toolboxArea.X),
		int32(w.inboxScene.toolboxArea.Y),
		int32(w.inboxScene.toolboxArea.Width),
		int32(w.inboxScene.toolboxArea.Height),
		rl.Gray)

	rl.DrawRectangle(
		int32(w.inboxScene.conversationArea.X),
		int32(w.inboxScene.conversationArea.Y),
		int32(w.inboxScene.conversationArea.Width),
		int32(w.inboxScene.conversationArea.Height),
		rl.White)

	rl.DrawRectangle(
		int32(w.inboxScene.messagesArea.X),
		int32(w.inboxScene.messagesArea.Y),
		int32(w.inboxScene.messagesArea.Width),
		int32(w.inboxScene.messagesArea.Height),
		rl.LightGray)

}
