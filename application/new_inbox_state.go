package application

import (
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

type messagePanelLayout struct {
	middle          float32
	padding         float32
	currHeight      float32
	messageWidth    float32
	messageFontSize int32
	leftSide        float32
	rightSide       float32
}
type InboxScene struct {
	messageChan      chan *proto.Message
	toolboxArea      rl.Rectangle
	conversationArea rl.Rectangle

	messagePanel       ScrollPanel
	messagePanelLayout messagePanelLayout
	usersConversations []*proto.ConversationSummary
	conversationsTabs  []ConversationTab
	messages           []Message
}

func (i *InboxScene) Reset() {

}

// TODO maybe use redis for fast cache to e.g user UUID
// TODO in refactor change all names to pattern verb+scene
func (w *Window) setupInboxScene() {

	//get users conversations
	res := w.ctx.Request(w.serverPID, &proto.GetLoggedInUUID{
		Pid: &proto.PID{
			Address: w.ctx.PID().Address,
			Id:      w.ctx.PID().ID}}, utils.WaitTime)
	resp, err := res.Result()
	if err != nil {
		//TODO STH
	}
	v, ok := resp.(*proto.LoggedInUUID)
	if !ok {
		//TODO
	}
	sender := v.Id
	res = w.ctx.Request(w.messageServicePID, &proto.GetUserConversation{Id: sender}, utils.WaitTime)
	resp, err = res.Result()
	_, ok = resp.(*proto.FailureGetUserConversation)
	if err != nil || !ok {
		//TODO error
	}
	if conversations, ok := resp.(*proto.SuccessGetUserConversation); ok {
		w.inboxScene.usersConversations = conversations.ConvSummary
	}

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

	w.inboxScene.messagePanel.bounds = rl.NewRectangle(
		w.inboxScene.toolboxArea.Width,
		w.inboxScene.toolboxArea.Y,
		(3.0/5.0)*float32(w.width),
		float32(w.height))
	w.inboxScene.messagePanel.content = rl.NewRectangle(
		w.inboxScene.toolboxArea.Width+5,
		w.inboxScene.toolboxArea.Y+5,
		(3.0/5.0)*float32(w.width)-15,
		float32(w.height)*10)

	w.inboxScene.messagePanel.view = rl.Rectangle{}
	w.inboxScene.messagePanel.scroll = rl.Vector2{}

	//todo check data type and change to float32 if its possible
	w.inboxScene.messagePanelLayout.middle = w.inboxScene.messagePanel.bounds.X + (w.inboxScene.messagePanel.bounds.Width)/2.0
	w.inboxScene.messagePanelLayout.padding = 10
	w.inboxScene.messagePanelLayout.currHeight = w.inboxScene.messagePanel.bounds.Y + w.inboxScene.messagePanelLayout.padding
	w.inboxScene.messagePanelLayout.messageWidth = 150
	w.inboxScene.messagePanelLayout.messageFontSize = 20
	w.inboxScene.messagePanelLayout.leftSide = w.inboxScene.messagePanel.bounds.X + w.inboxScene.messagePanelLayout.padding
	w.inboxScene.messagePanelLayout.rightSide = w.inboxScene.messagePanelLayout.middle + w.inboxScene.messagePanelLayout.padding

	var y float32 = w.inboxScene.toolboxArea.Height
	var height float32 = 40
	for _, conversation := range w.inboxScene.usersConversations {
		w.inboxScene.conversationsTabs = append(w.inboxScene.conversationsTabs, ConversationTab{
			withID:         conversation.WithID,
			conversationID: conversation.ConversationId,
			bounds:         rl.NewRectangle(w.inboxScene.toolboxArea.X, y, w.inboxScene.toolboxArea.Width, height),
			nametag:        conversation.Nametag,
			enterConversation: Button{
				bounds: rl.NewRectangle(
					(3.0/4.0)*w.inboxScene.toolboxArea.Width,
					y,
					80,
					40),
				text: "ENTER",
			},
		})

		y += height
	}
	w.inboxScene.messageChan = make(chan *proto.Message, 1024)
	go func() {
		//boxWidth := 200
		//spacing := 20
		//y := w.inboxScene.messagesArea.Y
		//	x := w.inboxScene.messagesArea.X
		/*	for msg := range w.inboxScene.messageChan {
				content := wrapText(int32(boxWidth), msg.Content, 15)
				height := (30)*strings.Count(content, "\n") + 1
				w.inboxScene.messages = append(w.inboxScene.messages, Message{
					bounds:  rl.NewRectangle(x, y, float32(boxWidth), float32(height)),
					content: content,
				})
				y += float32(spacing) + float32(height)

			}

		*/
	}()
}
func (w *Window) updateInboxState() {
	for i, tab := range w.inboxScene.conversationsTabs {
		if tab.isClicked {
			res := w.ctx.Request(w.serverPID, &proto.GetLoggedInUUID{
				Pid: &proto.PID{
					Address: w.ctx.PID().Address,
					Id:      w.ctx.PID().ID}}, utils.WaitTime)
			resp, err := res.Result()
			if err != nil {
				//TODO STH
			}
			v, ok := resp.(*proto.LoggedInUUID)
			if !ok {
				//TODO
			}
			sender := v.Id
			//open conversation
			w.ctx.Send(w.messageServicePID, &proto.UpdatePresence{
				Id: sender,
				Presence: &proto.PresenceType{
					Type: &proto.PresenceType_Inbox{
						Inbox: &proto.Inbox{
							WithID: tab.withID}}},
			})
			resp2 := w.ctx.Request(w.messageServicePID, &proto.OpenAndLoadConversation{
				UserID:         sender,
				ReceiverID:     tab.withID,
				ConversationID: tab.conversationID}, utils.WaitTime)
			res2, err2 := resp2.Result()
			//todo name var proper e.g res<what>
			if err2 != nil {
				panic(err2)
			}
			w.inboxScene.messages = w.inboxScene.messages[:0]
			if conversation, ok := res2.(*proto.SuccessOpenAndLoadConversation); ok {
				var xPosition float32
				var yPosition float32 = w.inboxScene.messagePanel.bounds.Y + w.inboxScene.messagePanelLayout.padding
				for _, msg := range conversation.Messages {
					if msg.SenderID == sender {
						xPosition = w.inboxScene.messagePanelLayout.rightSide
					} else {
						xPosition = w.inboxScene.messagePanelLayout.leftSide
					}

					//TODO repair styling these boxes
					content := wrapText(
						int32(w.inboxScene.messagePanelLayout.messageWidth),
						msg.Content,
						w.inboxScene.messagePanelLayout.messageFontSize,
					)
					height := float32(w.inboxScene.messagePanelLayout.messageFontSize) * float32(strings.Count(content, "\n")+1)
					w.inboxScene.messages = append(w.inboxScene.messages, Message{
						bounds: rl.NewRectangle(
							xPosition,
							yPosition,
							w.inboxScene.messagePanelLayout.messageWidth,
							height),
						content:   content,
						originalY: yPosition,
					})
					yPosition += 2*w.inboxScene.messagePanelLayout.padding + height
				}
			}
			//MSSVC -> ConversationManager where is spinning up new actor

			w.inboxScene.conversationsTabs[i].isClicked = false // to not load every time
		}

	}
}
func wrapText(maxWidth int32, input string, fontSize int32) string {
	var output strings.Builder
	var line strings.Builder
	for _, char := range input {
		line.WriteString(string(char))
		width := rl.MeasureText(line.String(), fontSize)
		if width >= maxWidth {
			output.WriteString("\n")
			line.Reset()
		}
		output.WriteString(string(char))
	}

	return output.String()
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

	gui.ScrollPanel(
		w.inboxScene.messagePanel.bounds,
		"MESSAGES",
		w.inboxScene.messagePanel.content,
		&w.inboxScene.messagePanel.scroll,
		&w.inboxScene.messagePanel.view,
	)
	rl.BeginScissorMode(
		int32(w.inboxScene.messagePanel.view.X),
		int32(w.inboxScene.messagePanel.view.Y),
		int32(w.inboxScene.messagePanel.view.Width),
		int32(w.inboxScene.messagePanel.view.Height),
	)

	for i := range w.inboxScene.messages {
		movingY := w.inboxScene.messages[i].originalY + +w.inboxScene.messagePanel.scroll.Y
		rl.DrawRectangle(
			int32(w.inboxScene.messages[i].bounds.X),
			int32(movingY),
			int32(w.inboxScene.messages[i].bounds.Width),
			int32(w.inboxScene.messages[i].bounds.Height),
			rl.SkyBlue)
		rl.DrawText(
			w.inboxScene.messages[i].content,
			int32(w.inboxScene.messages[i].bounds.X),
			int32(movingY),
			15,
			rl.White)

	}

	rl.EndScissorMode()

	for i := range w.inboxScene.conversationsTabs {
		rl.DrawRectangle(
			int32(w.inboxScene.conversationsTabs[i].bounds.X),
			int32(w.inboxScene.conversationsTabs[i].bounds.Y),
			int32(w.inboxScene.conversationsTabs[i].bounds.Width),
			int32(w.inboxScene.conversationsTabs[i].bounds.Height),
			rl.Red)
		rl.DrawText(
			w.inboxScene.conversationsTabs[i].nametag,
			int32(w.inboxScene.conversationsTabs[i].bounds.X),
			int32(w.inboxScene.conversationsTabs[i].bounds.Y),
			25,
			rl.Black)
		w.inboxScene.conversationsTabs[i].isClicked = gui.Button(w.inboxScene.conversationsTabs[i].enterConversation.bounds, w.inboxScene.conversationsTabs[i].enterConversation.text)
	}

}
