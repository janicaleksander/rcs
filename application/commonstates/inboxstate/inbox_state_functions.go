package inboxstate

import (
	"fmt"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *InboxScene) Reset() {
	i.toolboxSection.isGoBackButtonPressed = false
	i.toolboxSection.showAddConversationModal = false
	i.toolboxSection.isRefreshConversationPressed = false
	i.modalSection.isAcceptAddConversationPressed = false
	i.errorSection.message = ""
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
		i.errorSection.message = err.Error()
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}
	if v, ok := res.(*proto.LoggedInUUID); ok {
		i.tempUserID = v.Id
	} else {
		v, _ := res.(*proto.Error)
		i.errorSection.message = v.Content
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
	}
} //TODO Maybe setup error and its block next actions
func (i *InboxScene) GetUserConversation() {
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID,
		&proto.GetUserConversations{
			Id: i.tempUserID},
	))
	if err != nil {
		i.errorSection.message = err.Error()
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

	if v, ok := res.(*proto.UserConversations); ok {
		i.conversationSection.usersConversations = v.ConvSummary
	} else {
		v, _ := res.(*proto.Error)
		i.errorSection.message = v.Content
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
	}

}

func (i *InboxScene) GetUserToNewConversation() {
	res, err := utils.MakeRequest(utils.NewRequest(
		i.cfg.Ctx,
		i.cfg.MessageServicePID,
		&proto.GetUsersToNewConversation{
			UserID: i.tempUserID}))
	if err != nil {
		i.errorSection.message = err.Error()
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

	if v, ok := res.(*proto.UsersToNewConversation); ok {
		i.modalSection.users = v.Users
	} else {
		v, _ := res.(*proto.Error)
		i.errorSection.message = v.Content
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
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

	i.MessageSection.messages = append(i.MessageSection.messages, component.Message{
		Bounds: rl.NewRectangle(
			xPosition,
			currHeight,
			i.MessageSection.messagePanelLayout.messageWidth,
			height),
		Content:   content,
		SentAt:    msg.SentAt.AsTime().Format("02 Jan 2006 15:04"),
		OriginalY: i.MessageSection.messagePanelLayout.currHeight,
	})

	i.MessageSection.messagePanelLayout.mu.Lock()
	i.MessageSection.messagePanelLayout.currHeight += 2*i.MessageSection.messagePanelLayout.padding + height
	i.MessageSection.messagePanelLayout.mu.Unlock()
	i.MessageSection.messagePanel.Content.Height = i.MessageSection.messagePanelLayout.currHeight
	i.MessageSection.messagePanel.Scroll.Y = -i.MessageSection.messagePanel.Content.Height
}

func (i *InboxScene) addNewConversation() {
	if i.modalSection.usersSlider.IdxActiveElement == -1 {
		i.errorSection.message = "Select user!"
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
	} else {
		selectedUSer := i.modalSection.users[i.modalSection.usersSlider.IdxActiveElement]
		res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.CreateConversation{
			Id:         uuid.New().String(),
			SenderID:   i.tempUserID,
			ReceiverID: selectedUSer.Id,
		}))
		if err != nil {
			//context deadline exceeded
			i.errorSection.message = err.Error()
			i.errorSection.errorPopup.ShowFor(time.Second * 3)
			return
		}
		if _, ok := res.(*proto.AcceptCreateConversation); ok {
			i.refreshConversationsPanel()
		} else {
			v, _ := res.(*proto.Error)
			i.errorSection.message = v.Content
			i.errorSection.errorPopup.ShowFor(time.Second * 3)
		}
	}

}

// every refresh its possible to gave other sequence of tabs
func (i *InboxScene) refreshConversationsPanel() {
	i.conversationSection.activeConversation = -1
	i.conversationSection.conversationsTabs = i.conversationSection.conversationsTabs[:0]
	res, err := utils.MakeRequest(utils.NewRequest(
		i.cfg.Ctx,
		i.cfg.MessageServicePID,
		&proto.GetUserConversations{Id: i.tempUserID}))
	if err != nil {
		i.errorSection.message = err.Error()
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}
	if v, ok := res.(*proto.UserConversations); ok {
		fmt.Println(v.ConvSummary)
		i.conversationSection.usersConversations = v.ConvSummary
	} else {
		v, _ := res.(*proto.Error)
		i.errorSection.message = v.Content
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

	i.conversationSection.conversationPanelLayout.currHeight = i.conversationSection.conversationPanelLayout.startHeight
	for k, conversation := range i.conversationSection.usersConversations {
		i.conversationSection.conversationsTabs = append(i.conversationSection.conversationsTabs, component.ConversationTab{
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
			EnterConversation: *component.NewButton(
				component.NewButtonConfig(),
				rl.NewRectangle(
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

func (i *InboxScene) LoadMessages(tab *component.ConversationTab) {
	//open conversation
	i.cfg.Ctx.Send(i.cfg.MessageServicePID, &proto.UpdatePresence{
		Id: i.tempUserID,
		Presence: &proto.PresenceType{
			Type: &proto.PresenceType_Inbox{
				Inbox: &proto.Inbox{
					WithID: tab.WithID}}},
	})
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.OpenAndLoadConversation{
		UserID:         i.tempUserID,
		ReceiverID:     tab.WithID,
		ConversationID: tab.ConversationID,
	},
	))
	if err != nil {
		i.errorSection.message = err.Error()
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}
	i.MessageSection.messages = i.MessageSection.messages[:0]
	i.MessageSection.messagePanelLayout.currHeight = i.MessageSection.messagePanel.Bounds.Y + 3*i.MessageSection.messagePanelLayout.padding
	if v, ok := res.(*proto.LoadedConversation); ok {
		for _, msg := range v.Messages {
			i.AppendMessage(msg)
		}
	} else {
		v, _ := res.(*proto.Error)
		i.errorSection.message = v.Content
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
	}
}

func (i *InboxScene) SendMessage() {
	msg := &proto.Message{
		Id:             uuid.New().String(),
		SenderID:       i.tempUserID,
		ConversationID: i.conversationSection.activeConversationID,
		Content:        i.MessageSection.textInput.GetText(),
		SentAt:         timestamppb.Now(),
	}
	res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.SendMessage{
		Receiver: i.conversationSection.activeWithID,
		Message:  msg,
	}))
	if err != nil {
		i.errorSection.message = err.Error()
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
		return
	}

	if _, ok := res.(*proto.AcceptSend); ok {
		i.AppendMessage(msg)
	} else {
		v, _ := res.(*proto.Error)
		i.errorSection.message = v.Content
		i.errorSection.errorPopup.ShowFor(time.Second * 3)
	}

}
