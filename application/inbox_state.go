package application

import (
	"fmt"
	"strings"
	"sync"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type messagePanelLayout struct {
	middle          float32
	padding         float32
	currHeight      float32
	messageWidth    float32
	messageFontSize int32
	leftSide        float32
	rightSide       float32
	mu              sync.RWMutex
}
type InboxScene struct {
	messageChan        chan *proto.Message
	toolboxArea        rl.Rectangle
	conversationArea   rl.Rectangle
	tempWithID         string
	tempWithConvID     string
	messagePanel       ScrollPanel
	messagePanelLayout messagePanelLayout

	textInput           component.InputBox
	sendButton          Button
	isSendButtonPressed bool
	usersConversations  []*proto.ConversationSummary
	conversationsTabs   []ConversationTab
	messages            []Message
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

	w.inboxScene.textInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			w.inboxScene.messagePanel.bounds.X,
			w.inboxScene.messagePanel.bounds.Height-30, //height
			w.inboxScene.messagePanel.bounds.Width-70,  //for button
			30),
		false)

	w.inboxScene.sendButton = Button{
		bounds: rl.NewRectangle(
			w.inboxScene.messagePanel.bounds.X+w.inboxScene.textInput.Bounds.Width,
			w.inboxScene.messagePanel.bounds.Height-30,
			65,
			30),
		text: "SEND",
	}
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
		for msg := range w.inboxScene.messageChan {
			var xPosition float32
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
			w.inboxScene.messagePanelLayout.mu.Lock()
			currHeight := w.inboxScene.messagePanelLayout.currHeight
			w.inboxScene.messagePanelLayout.mu.Unlock()

			w.inboxScene.messages = append(w.inboxScene.messages, Message{
				bounds: rl.NewRectangle(
					xPosition,
					currHeight,
					w.inboxScene.messagePanelLayout.messageWidth,
					height),
				content:   content,
				originalY: w.inboxScene.messagePanelLayout.currHeight,
			})
			w.inboxScene.messagePanelLayout.mu.Lock()
			w.inboxScene.messagePanelLayout.currHeight += 2*w.inboxScene.messagePanelLayout.padding + height
			w.inboxScene.messagePanelLayout.mu.Unlock()

		}

	}()
}

// repair reset e.g. currHeight on every click at other conv
func (w *Window) updateInboxState() {
	w.inboxScene.textInput.Update()
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
			w.inboxScene.tempWithID = tab.withID
			w.inboxScene.tempWithConvID = tab.conversationID
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
					w.inboxScene.messagePanelLayout.mu.Lock()
					currHeight := w.inboxScene.messagePanelLayout.currHeight
					w.inboxScene.messagePanelLayout.mu.Unlock()
					w.inboxScene.messages = append(w.inboxScene.messages, Message{
						bounds: rl.NewRectangle(
							xPosition,
							currHeight,
							w.inboxScene.messagePanelLayout.messageWidth,
							height),
						content:   content,
						originalY: w.inboxScene.messagePanelLayout.currHeight,
					})
					w.inboxScene.messagePanelLayout.mu.Lock()
					w.inboxScene.messagePanelLayout.currHeight += 2*w.inboxScene.messagePanelLayout.padding + height
					w.inboxScene.messagePanelLayout.mu.Unlock()
				}
			}
			//MSSVC -> ConversationManager where is spinning up new actor

			w.inboxScene.conversationsTabs[i].isClicked = false // to not load every time
		}

	}
	if w.inboxScene.isSendButtonPressed {
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
		resp2 := w.ctx.Request(w.messageServicePID, &proto.SendMessage{
			Receiver: w.inboxScene.tempWithID,
			Message: &proto.Message{
				Id:             uuid.New().String(),
				SenderID:       sender,
				ConversationID: w.inboxScene.tempWithConvID,
				Content:        w.inboxScene.textInput.GetText(),
				SentAt:         timestamppb.Now(),
			},
		}, utils.WaitTime)
		res2, err := resp2.Result()
		if err != nil {
			panic(err)
		}
		if message, ok := res2.(*proto.DeliverMessage); ok {
			var xPosition float32
			if message.Message.SenderID == message.Receiver {
				xPosition = w.inboxScene.messagePanelLayout.rightSide
			} else {
				xPosition = w.inboxScene.messagePanelLayout.leftSide
			}

			//TODO repair styling these boxes
			content := wrapText(
				int32(w.inboxScene.messagePanelLayout.messageWidth),
				message.Message.Content,
				w.inboxScene.messagePanelLayout.messageFontSize,
			)
			height := float32(w.inboxScene.messagePanelLayout.messageFontSize) * float32(strings.Count(content, "\n")+1)
			w.inboxScene.messagePanelLayout.mu.Lock()
			currHeight := w.inboxScene.messagePanelLayout.currHeight
			w.inboxScene.messagePanelLayout.mu.Unlock()
			w.inboxScene.messages = append(w.inboxScene.messages, Message{
				bounds: rl.NewRectangle(
					xPosition,
					currHeight,
					w.inboxScene.messagePanelLayout.messageWidth,
					height),
				content:   content,
				originalY: w.inboxScene.messagePanelLayout.currHeight,
			})
			w.inboxScene.messagePanelLayout.mu.Lock()
			w.inboxScene.messagePanelLayout.currHeight += 2*w.inboxScene.messagePanelLayout.padding + height
			w.inboxScene.messagePanelLayout.mu.Unlock()
		}
	}

}

//TODO add starting positon to last message and increase on the fly the content area

// TODO Repair
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
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, w.inboxScene.textInput.Bounds) {
			w.inboxScene.textInput.Active()
		} else {
			w.inboxScene.textInput.Deactivate()
		}

	}
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
		fmt.Println(len(w.inboxScene.messages))
		movingY := w.inboxScene.messages[i].originalY + w.inboxScene.messagePanel.scroll.Y
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
	w.inboxScene.textInput.Render()
	w.inboxScene.isSendButtonPressed = gui.Button(w.inboxScene.sendButton.bounds, w.inboxScene.sendButton.text)
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

//TODO when i see e.g. modal disable all buttons
//e.g. in info user when i am sending message i would click other buttons from other layers
