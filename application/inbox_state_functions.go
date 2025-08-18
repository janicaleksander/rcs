package application

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (i *InboxScene) Reset() {
	i.toolboxSection.isGoBackButtonPressed = false
	i.toolboxSection.showAddConversationModal = false // if we want to show window to add new conversation
	i.toolboxSection.isRefreshConversationPressed = false
	i.modalSection.isAcceptAddConversationPressed = false
	i.modalSection.isErrorModal = false
	i.modalSection.textErrorModal = ""
	i.messageSection.isSendButtonPressed = false
	i.conversationSection.isConversationSelected = false
	i.messageSection.messages = i.messageSection.messages[:0]

}

func (w *Window) GetLoggedID() {
	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.serverPID, &proto.GetLoggedInUUID{
		Pid: &proto.PID{
			Address: w.ctx.PID().Address,
			Id:      w.ctx.PID().ID,
		},
	}))
	if err != nil {
		//TODO error ctx deadline exceeded
	}
	if v, ok := res.(*proto.LoggedInUUID); ok {
		w.inboxScene.tempUserID = v.Id
	} else {
		//TODO some general error and if its true we cant go further (maybe some error screen)
	}
}
func (w *Window) GetUserConversation() {

	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID,
		&proto.GetUserConversation{
			Id: w.inboxScene.tempUserID},
	))
	if err != nil {
		//TODO error ctx deadline exceeded
	}

	if v, ok := res.(*proto.SuccessGetUserConversation); ok {
		w.inboxScene.conversationSection.usersConversations = v.ConvSummary
	} else {
		// maybe one message proto.Error with msg
	}

}

func (w *Window) GetUserToNewConversation() {
	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.GetUsersToNewConversation{Id: w.inboxScene.tempUserID}))
	if err != nil {
		//TODO error ctx deadline exceeded
	}

	if v, ok := res.(*proto.SuccessUsersToNewConversation); ok {
		w.inboxScene.modalSection.users = v.Users
	} else {
		// maybe one message proto.Error with msg
	}
}

func (w *Window) AppendMessage(msg *proto.Message) {
	var xPosition float32
	if msg.SenderID == w.inboxScene.tempUserID {
		xPosition = w.inboxScene.messageSection.messagePanelLayout.rightSide
	} else {
		xPosition = w.inboxScene.messageSection.messagePanelLayout.leftSide
	}

	//TODO repair styling these boxes
	content := utils.WrapText(
		int32(w.inboxScene.messageSection.messagePanelLayout.messageWidth),
		msg.Content,
		w.inboxScene.messageSection.messagePanelLayout.messageFontSize,
	)
	height := float32(w.inboxScene.messageSection.messagePanelLayout.messageFontSize) * float32(strings.Count(content, "\n")+1)
	w.inboxScene.messageSection.messagePanelLayout.mu.Lock()
	currHeight := w.inboxScene.messageSection.messagePanelLayout.currHeight
	w.inboxScene.messageSection.messagePanelLayout.mu.Unlock()

	w.inboxScene.messageSection.messages = append(w.inboxScene.messageSection.messages, Message{
		bounds: rl.NewRectangle(
			xPosition,
			currHeight,
			w.inboxScene.messageSection.messagePanelLayout.messageWidth,
			height),
		content:   content,
		originalY: w.inboxScene.messageSection.messagePanelLayout.currHeight,
	})

	w.inboxScene.messageSection.messagePanelLayout.mu.Lock()
	w.inboxScene.messageSection.messagePanelLayout.currHeight += 2*w.inboxScene.messageSection.messagePanelLayout.padding + height
	w.inboxScene.messageSection.messagePanelLayout.mu.Unlock()
	w.inboxScene.messageSection.messagePanel.content.Height = w.inboxScene.messageSection.messagePanelLayout.currHeight
	w.inboxScene.messageSection.messagePanel.scroll.Y = -w.inboxScene.messageSection.messagePanel.content.Height
	// - or +?
}

func (w *Window) addNewConversation() {
	w.inboxScene.modalSection.isErrorModal = false
	w.inboxScene.modalSection.textErrorModal = ""
	if w.inboxScene.modalSection.usersSlider.idxActiveElement == -1 {
		w.inboxScene.modalSection.isErrorModal = true
		w.inboxScene.modalSection.textErrorModal = "Select user!"
	} else {
		selectedUSer := w.inboxScene.modalSection.users[w.inboxScene.modalSection.usersSlider.idxActiveElement]
		res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.CreateConversation{
			Id:         uuid.New().String(),
			SenderID:   w.inboxScene.tempUserID,
			ReceiverID: selectedUSer.Id,
		}))
		if err != nil {
			//context deadline exceeded
			w.inboxScene.modalSection.isErrorModal = true
			w.inboxScene.modalSection.textErrorModal = "Error" + err.Error()

		}

		if _, ok := res.(*proto.SuccessOfCreateConversation); ok {
			w.refreshConversationsPanel()
		} else {
			w.inboxScene.modalSection.isErrorModal = true
			w.inboxScene.modalSection.textErrorModal = "Error"
			//error
		}
	}

}

// every refresh its possible to gave other sequence of tabs
func (w *Window) refreshConversationsPanel() {
	w.inboxScene.conversationSection.activeConversation = -1
	w.inboxScene.conversationSection.conversationsTabs = w.inboxScene.conversationSection.conversationsTabs[:0]
	res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID,
		&proto.GetUserConversation{Id: w.inboxScene.tempUserID}))
	if err != nil {
		//TODO context deadline exceeded
	}

	if v, ok := res.(*proto.SuccessGetUserConversation); ok {
		fmt.Println(v.ConvSummary)
		w.inboxScene.conversationSection.usersConversations = v.ConvSummary
	} else {
		//TODO error
	}

	w.inboxScene.conversationSection.conversationPanelLayout.currHeight = w.inboxScene.conversationSection.conversationPanelLayout.startHeight
	for i, conversation := range w.inboxScene.conversationSection.usersConversations {
		w.inboxScene.conversationSection.conversationsTabs = append(w.inboxScene.conversationSection.conversationsTabs, ConversationTab{
			ID:             int32(i),
			withID:         conversation.WithID,
			conversationID: conversation.ConversationId,
			bounds: rl.NewRectangle(
				w.inboxScene.toolboxSection.toolboxArea.X,
				w.inboxScene.conversationSection.conversationPanelLayout.currHeight,
				w.inboxScene.toolboxSection.toolboxArea.Width,
				w.inboxScene.conversationSection.conversationPanelLayout.panelHeight,
			),
			originalY: w.inboxScene.conversationSection.conversationPanelLayout.currHeight,
			nametag:   conversation.Nametag,
			enterConversation: *component.NewButton(component.NewButtonConfig(), rl.NewRectangle(
				(3.0/4.0)*w.inboxScene.toolboxSection.toolboxArea.Width,
				w.inboxScene.conversationSection.conversationPanelLayout.currHeight,
				80,
				40), "ENTER", true),
		})
		w.inboxScene.conversationSection.conversationPanelLayout.currHeight += w.inboxScene.conversationSection.conversationPanelLayout.panelHeight
		w.inboxScene.conversationSection.conversationPanel.content.Height = w.inboxScene.conversationSection.conversationPanelLayout.currHeight
		w.inboxScene.conversationSection.conversationPanel.scroll.Y = w.inboxScene.conversationSection.conversationPanel.content.Height

	}
}
