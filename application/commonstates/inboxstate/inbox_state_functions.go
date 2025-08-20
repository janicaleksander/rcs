package inboxstate

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	component2 "github.com/janicaleksander/bcs/component"
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
	i.MessageSection.isSendButtonPressed = false
	i.conversationSection.isConversationSelected = false
	i.MessageSection.messages = i.MessageSection.messages[:0]

}

func (i *InboxScene) GetLoggedID() {
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.ServerPID, &proto.GetLoggedInUUID{
		Pid: &proto.PID{
			Address: i.cfg.Ctx.PID().Address,
			Id:      i.cfg.Ctx.PID().ID,
		},
	}))
	if err != nil {
		//TODO error ctx deadline exceeded
	}
	if v, ok := res.(*proto.LoggedInUUID); ok {
		i.tempUserID = v.Id
	} else {
		//TODO some general error and if its true we cant go further (maybe some error screen)
	}
}
func (i *InboxScene) GetUserConversation() {

	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID,
		&proto.GetUserConversation{
			Id: i.tempUserID},
	))
	if err != nil {
		//TODO error ctx deadline exceeded
	}

	if v, ok := res.(*proto.SuccessGetUserConversation); ok {
		i.conversationSection.usersConversations = v.ConvSummary
	} else {
		// maybe one message proto.Error with msg
	}

}

func (i *InboxScene) GetUserToNewConversation() {
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.GetUsersToNewConversation{Id: i.tempUserID}))
	if err != nil {
		//TODO error ctx deadline exceeded
	}

	if v, ok := res.(*proto.SuccessUsersToNewConversation); ok {
		i.modalSection.users = v.Users
	} else {
		// maybe one message proto.Error with msg
	}
}

func (i *InboxScene) AppendMessage(msg *proto.Message) {
	var xPosition float32
	if msg.SenderID == i.tempUserID {
		xPosition = i.MessageSection.messagePanelLayout.rightSide
	} else {
		xPosition = i.MessageSection.messagePanelLayout.leftSide
	}

	//TODO repair styling these boxes
	content := utils.WrapText(
		int32(i.MessageSection.messagePanelLayout.messageWidth),
		msg.Content,
		i.MessageSection.messagePanelLayout.messageFontSize,
	)
	height := float32(i.MessageSection.messagePanelLayout.messageFontSize) * float32(strings.Count(content, "\n")+1)
	i.MessageSection.messagePanelLayout.mu.Lock()
	currHeight := i.MessageSection.messagePanelLayout.currHeight
	i.MessageSection.messagePanelLayout.mu.Unlock()

	i.MessageSection.messages = append(i.MessageSection.messages, component2.Message{
		Bounds: rl.NewRectangle(
			xPosition,
			currHeight,
			i.MessageSection.messagePanelLayout.messageWidth,
			height),
		Content:   content,
		OriginalY: i.MessageSection.messagePanelLayout.currHeight,
	})

	i.MessageSection.messagePanelLayout.mu.Lock()
	i.MessageSection.messagePanelLayout.currHeight += 2*i.MessageSection.messagePanelLayout.padding + height
	i.MessageSection.messagePanelLayout.mu.Unlock()
	i.MessageSection.messagePanel.Content.Height = i.MessageSection.messagePanelLayout.currHeight
	i.MessageSection.messagePanel.Scroll.Y = -i.MessageSection.messagePanel.Content.Height
	// - or +?
}

func (i *InboxScene) addNewConversation() {
	i.modalSection.isErrorModal = false
	i.modalSection.textErrorModal = ""
	if i.modalSection.usersSlider.IdxActiveElement == -1 {
		i.modalSection.isErrorModal = true
		i.modalSection.textErrorModal = "Select user!"
	} else {
		selectedUSer := i.modalSection.users[i.modalSection.usersSlider.IdxActiveElement]
		res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.CreateConversation{
			Id:         uuid.New().String(),
			SenderID:   i.tempUserID,
			ReceiverID: selectedUSer.Id,
		}))
		if err != nil {
			//context deadline exceeded
			i.modalSection.isErrorModal = true
			i.modalSection.textErrorModal = "Error" + err.Error()

		}

		if _, ok := res.(*proto.SuccessOfCreateConversation); ok {
			i.refreshConversationsPanel()
		} else {
			i.modalSection.isErrorModal = true
			i.modalSection.textErrorModal = "Error"
			//error
		}
	}

}

// every refresh its possible to gave other sequence of tabs
func (i *InboxScene) refreshConversationsPanel() {
	i.conversationSection.activeConversation = -1
	i.conversationSection.conversationsTabs = i.conversationSection.conversationsTabs[:0]
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID,
		&proto.GetUserConversation{Id: i.tempUserID}))
	if err != nil {
		//TODO context deadline exceeded
	}

	if v, ok := res.(*proto.SuccessGetUserConversation); ok {
		fmt.Println(v.ConvSummary)
		i.conversationSection.usersConversations = v.ConvSummary
	} else {
		//TODO error
	}

	i.conversationSection.conversationPanelLayout.currHeight = i.conversationSection.conversationPanelLayout.startHeight
	for k, conversation := range i.conversationSection.usersConversations {
		i.conversationSection.conversationsTabs = append(i.conversationSection.conversationsTabs, component2.ConversationTab{
			ID:             int32(k),
			WithID:         conversation.WithID,
			ConversationID: conversation.ConversationId,
			Bounds: rl.NewRectangle(
				i.toolboxSection.toolboxArea.X,
				i.conversationSection.conversationPanelLayout.currHeight,
				i.toolboxSection.toolboxArea.Width,
				i.conversationSection.conversationPanelLayout.panelHeight,
			),
			OriginalY: i.conversationSection.conversationPanelLayout.currHeight,
			Nametag:   conversation.Nametag,
			EnterConversation: *component2.NewButton(component2.NewButtonConfig(), rl.NewRectangle(
				(3.0/4.0)*i.toolboxSection.toolboxArea.Width,
				i.conversationSection.conversationPanelLayout.currHeight,
				80,
				40), "ENTER", true),
		})
		i.conversationSection.conversationPanelLayout.currHeight += i.conversationSection.conversationPanelLayout.panelHeight
		i.conversationSection.conversationPanel.Content.Height = i.conversationSection.conversationPanelLayout.currHeight
		i.conversationSection.conversationPanel.Scroll.Y = i.conversationSection.conversationPanel.Content.Height

	}
}
