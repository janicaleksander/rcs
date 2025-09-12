package inboxstate

import (
	"sync"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ToolboxSection struct {
	toolboxArea                rl.Rectangle     // area with navigation buttons
	backButton                 component.Button // go back button
	addConversationButton      component.Button // button to add new conversation
	refreshConversationsButton component.Button // button to refresh users conversation

	isGoBackButtonPressed        bool
	showAddConversationModal     bool // if we want to show window to add new conversation
	isRefreshConversationPressed bool
}

type ModalSection struct { // Modal add conversation
	addConversationModal           component.Modal
	usersSlider                    component.ListSlider // user slider inside modals
	users                          []*proto.User        // users from DB without logged in user
	acceptAddConversationButton    component.Button     // button inside modal to confirm
	isAcceptAddConversationPressed bool
	//error inside modal
	errorBoxModal  rl.Rectangle
	isErrorModal   bool
	textErrorModal string
}

type MessageSection struct {
	textInput           component.InputBox
	sendButton          component.Button
	isSendButtonPressed bool

	MessageChan        chan *proto.Message   //chan to transport messages from window actor
	messages           []component.Message   // messages for current conversation
	messagePanel       component.ScrollPanel //slider with all messages
	messagePanelLayout messagePanelLayout    // messagePanel configuration
	nameOnTheWindow    string
}
type ConversationSection struct {
	usersConversations      []*proto.ConversationSummary
	conversationPanelLayout conversationPanelLayout
	conversationPanel       component.ScrollPanel
	conversationsTabs       []component.ConversationTab
	isConversationSelected  bool
	activeConversation      int32
	activeConversationID    string
	activeWithID            string
	nameOnTheWindow         string
}

type InboxScene struct {
	cfg *utils.SharedConfig
	//scheduler
	stateManager        *statesmanager.StateManager
	tempUserID          string // id of current logged in a user
	toolboxSection      ToolboxSection
	modalSection        ModalSection
	MessageSection      MessageSection
	conversationSection ConversationSection
}

type conversationPanelLayout struct {
	currHeight  float32
	panelHeight float32
	startHeight float32
}

type messagePanelLayout struct {
	padding         float32
	middle          float32
	currHeight      float32
	messageWidth    float32
	messageFontSize int32
	leftSide        float32
	rightSide       float32
	mu              sync.RWMutex
}

// TODO maybe use redis for fast cache to e.g user UUID
// TODO in refactor change all names to pattern verb+scene
func (i *InboxScene) SetupInboxScene(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	i.cfg = cfg
	i.stateManager = state
	i.Reset()
	i.GetLoggedID()
	i.GetUserConversation()
	i.GetUserToNewConversation()

	i.toolboxSection.toolboxArea = rl.NewRectangle(
		0,
		0,
		(2.0/5.0)*float32(rl.GetScreenWidth()),
		(1.0/8.0)*float32(rl.GetScreenHeight()))

	//go back button
	var toolboxButtonPadding float32 = 10
	var toolBoxButtonWidth float32 = 150
	var toolBoxButtonHeight float32 = 50

	i.toolboxSection.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.toolboxSection.toolboxArea.X+toolboxButtonPadding,
			i.toolboxSection.toolboxArea.Y+toolboxButtonPadding,
			toolBoxButtonWidth,
			toolBoxButtonHeight),
		"GO BACK",
		false)

	i.toolboxSection.addConversationButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.toolboxSection.backButton.Bounds.X+i.toolboxSection.backButton.Bounds.Width+toolboxButtonPadding,
			i.toolboxSection.toolboxArea.Y+toolboxButtonPadding,
			toolBoxButtonWidth,
			toolBoxButtonHeight),
		"Add \n conversation",
		false)

	i.modalSection.addConversationModal = component.Modal{
		Background: rl.NewRectangle(
			0,
			0,
			float32(rl.GetScreenWidth()),
			float32(rl.GetScreenHeight())),
		BgColor: rl.Fade(rl.LightGray, 0.2),
		Core:    rl.NewRectangle(float32(rl.GetScreenWidth()/2-150.0), float32(rl.GetScreenHeight()/2-150.0), 300, 280),
	}
	var addConversationModalPadding float32 = 20
	i.modalSection.usersSlider = component.ListSlider{
		Strings: make([]string, 0, 64),
		Bounds: rl.NewRectangle(
			i.modalSection.addConversationModal.Core.X+i.modalSection.addConversationModal.Core.Width/2-(i.modalSection.addConversationModal.Core.Width/2)/2,
			i.modalSection.addConversationModal.Core.Y+2*addConversationModalPadding,
			i.modalSection.addConversationModal.Core.Width/2,
			80),
		IdxActiveElement: 0,
		Focus:            0,
		IdxScroll:        0,
	}
	//fill above slice with users
	for _, user := range i.modalSection.users {
		i.modalSection.usersSlider.Strings = append(i.modalSection.usersSlider.Strings,
			user.Email+"\n"+user.Personal.Name+user.Personal.Name)
	}

	i.modalSection.acceptAddConversationButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.modalSection.usersSlider.Bounds.X,
			i.modalSection.usersSlider.Bounds.Y+5*addConversationModalPadding,
			i.modalSection.usersSlider.Bounds.Width,
			25),
		"Add conversation", false)

	i.modalSection.errorBoxModal = rl.NewRectangle(
		i.modalSection.usersSlider.Bounds.X,
		i.modalSection.usersSlider.Bounds.Y+8*addConversationModalPadding,
		i.modalSection.usersSlider.Bounds.Width,
		75)

	i.toolboxSection.refreshConversationsButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.toolboxSection.addConversationButton.Bounds.X+i.toolboxSection.addConversationButton.Bounds.Width+toolboxButtonPadding,
			i.toolboxSection.toolboxArea.Y+toolboxButtonPadding,
			toolBoxButtonWidth,
			toolBoxButtonHeight),
		"Fetch \n conversations", false)

	//conversationPanel slider
	i.conversationSection.conversationPanel.Bounds = rl.NewRectangle(
		i.toolboxSection.toolboxArea.X,
		i.toolboxSection.toolboxArea.Height,
		i.toolboxSection.toolboxArea.Width,
		(7.0/8.0)*float32(rl.GetScreenHeight()))

	i.conversationSection.conversationPanel.Content = rl.NewRectangle(
		i.toolboxSection.toolboxArea.X,
		i.toolboxSection.toolboxArea.Height,
		i.toolboxSection.toolboxArea.Width-15,
		(7.0/8.0)*float32(rl.GetScreenHeight()))

	i.conversationSection.nameOnTheWindow = "CONVERSATIONS"
	i.conversationSection.conversationPanel.View = rl.Rectangle{}
	i.conversationSection.conversationPanel.Scroll = rl.Vector2{}

	//messagePanel slider
	i.MessageSection.messagePanel.Bounds = rl.NewRectangle(
		i.toolboxSection.toolboxArea.Width,
		i.toolboxSection.toolboxArea.Y,
		(3.0/5.0)*float32(rl.GetScreenWidth()),
		float32(rl.GetScreenHeight()))
	i.MessageSection.messagePanel.Content = rl.NewRectangle(
		i.toolboxSection.toolboxArea.Width,
		i.toolboxSection.toolboxArea.Y,
		(3.0/5.0)*float32(rl.GetScreenWidth())-15,
		float32(rl.GetScreenHeight())*10)

	i.MessageSection.nameOnTheWindow = "MESSAGES"
	i.MessageSection.messagePanel.View = rl.Rectangle{}
	i.MessageSection.messagePanel.Scroll = rl.Vector2{}

	i.MessageSection.messagePanelLayout.padding = 20.0
	i.MessageSection.messagePanelLayout.middle = i.MessageSection.messagePanel.Bounds.X + (i.MessageSection.messagePanel.Bounds.Width)/2.0
	i.MessageSection.messagePanelLayout.currHeight = i.MessageSection.messagePanel.Bounds.Y + 3*i.MessageSection.messagePanelLayout.padding
	i.MessageSection.messagePanelLayout.messageWidth = 150
	i.MessageSection.messagePanelLayout.messageFontSize = 20
	i.MessageSection.messagePanelLayout.leftSide = i.MessageSection.messagePanel.Bounds.X + i.MessageSection.messagePanelLayout.padding
	i.MessageSection.messagePanelLayout.rightSide = i.MessageSection.messagePanelLayout.middle + i.MessageSection.messagePanelLayout.padding

	i.MessageSection.textInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			i.MessageSection.messagePanel.Bounds.X,
			i.MessageSection.messagePanel.Bounds.Height-30, //height
			i.MessageSection.messagePanel.Bounds.Width-70,  //for button// TODO vars
			30))

	i.MessageSection.sendButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			i.MessageSection.messagePanel.Bounds.X+i.MessageSection.textInput.Bounds.Width,
			i.MessageSection.messagePanel.Bounds.Height-30,
			65,
			30),
		"SEND",
		false)
	i.conversationSection.conversationPanelLayout = conversationPanelLayout{
		currHeight:  80,
		panelHeight: 80,
		startHeight: 120,
	}
	i.refreshConversationsPanel()
	i.MessageSection.MessageChan = make(chan *proto.Message, 1024)
	go func() {
		for msg := range i.MessageSection.MessageChan {
			i.AppendMessage(msg)
		}
	}()

}

func (i *InboxScene) UpdateInboxState() {

	modalOpen := i.toolboxSection.showAddConversationModal
	i.toolboxSection.backButton.SetActive(!modalOpen)
	i.toolboxSection.addConversationButton.SetActive(!modalOpen)
	i.toolboxSection.refreshConversationsButton.SetActive(!modalOpen)
	i.MessageSection.sendButton.SetActive(!modalOpen)

	i.toolboxSection.isGoBackButtonPressed = i.toolboxSection.backButton.Update()
	if i.toolboxSection.addConversationButton.Update() {
		// we need this if to remain windowbox on the screen
		//in other cases we can use one line because e.g we send sth so we need this true state for short time (moment)
		i.toolboxSection.showAddConversationModal = true
	}
	i.toolboxSection.isRefreshConversationPressed = i.toolboxSection.refreshConversationsButton.Update()
	i.modalSection.isAcceptAddConversationPressed = i.modalSection.acceptAddConversationButton.Update()
	i.MessageSection.isSendButtonPressed = i.MessageSection.sendButton.Update()

	for k := range i.conversationSection.conversationsTabs {
		if i.conversationSection.conversationsTabs[k].EnterConversation.Update() {
			i.conversationSection.conversationsTabs[k].IsPressed = true
			i.conversationSection.conversationsTabs[k].EnterConversation.SetActive(!modalOpen)
		}
	}

	//ACTIONS AFTER BUTTON UPDATES
	i.MessageSection.textInput.Update()

	if i.toolboxSection.isGoBackButtonPressed {
		i.stateManager.Add(statesmanager.GoBackState)
		return
	}
	i.conversationSection.activeConversation = -1
	for k, tab := range i.conversationSection.conversationsTabs {
		if tab.IsPressed {
			i.MessageSection.nameOnTheWindow = "MESSAGES WITH " + i.conversationSection.usersConversations[tab.ID].Nametag

			//we have to mark this to know if we have to show input box with send button
			i.conversationSection.isConversationSelected = true
			i.conversationSection.activeWithID = tab.WithID
			i.conversationSection.activeConversationID = tab.ConversationID
			if tab.ID != i.conversationSection.activeConversation {
				i.conversationSection.activeConversation = tab.ID

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
					//TODO context deadline exceeded
				}
				i.MessageSection.messages = i.MessageSection.messages[:0]
				i.MessageSection.messagePanelLayout.currHeight = i.MessageSection.messagePanel.Bounds.Y + 3*i.MessageSection.messagePanelLayout.padding
				if v, ok := res.(*proto.LoadedConversation); ok {
					for _, msg := range v.Messages {
						i.AppendMessage(msg)
					}
				} else {
					//TODO err
				}
				i.conversationSection.conversationsTabs[k].IsPressed = false // to not load every time
			}
		}

	}
	if i.MessageSection.isSendButtonPressed {
		res, err := utils.MakeRequest(utils.NewRequest(i.cfg.Ctx, i.cfg.MessageServicePID, &proto.SendMessage{
			Receiver: i.conversationSection.activeWithID,
			Message: &proto.Message{
				Id:             uuid.New().String(),
				SenderID:       i.tempUserID,
				ConversationID: i.conversationSection.activeConversationID,
				Content:        i.MessageSection.textInput.GetText(),
				SentAt:         timestamppb.Now(),
			},
		}))
		if err != nil {
			// context deadline exceed
		}

		if v, ok := res.(*proto.DeliverMessage); ok {
			i.AppendMessage(v.Message)
		} else {
			//error
		}

	}

	if i.modalSection.isAcceptAddConversationPressed {
		i.addNewConversation()
	}
	if i.toolboxSection.isRefreshConversationPressed {
		i.refreshConversationsPanel()
	}

}

func (i *InboxScene) RenderInboxState() {
	//toolbox
	rl.DrawRectangle(
		int32(i.toolboxSection.toolboxArea.X),
		int32(i.toolboxSection.toolboxArea.Y),
		int32(i.toolboxSection.toolboxArea.Width),
		int32(i.toolboxSection.toolboxArea.Height),
		rl.Gray)

	//toolbox button section
	i.toolboxSection.backButton.Render()
	i.toolboxSection.addConversationButton.Render()
	i.toolboxSection.refreshConversationsButton.Render()

	//MessageSection slider
	gui.ScrollPanel(
		i.MessageSection.messagePanel.Bounds,
		i.MessageSection.nameOnTheWindow,
		i.MessageSection.messagePanel.Content,
		&i.MessageSection.messagePanel.Scroll,
		&i.MessageSection.messagePanel.View,
	)
	rl.BeginScissorMode(
		int32(i.MessageSection.messagePanel.View.X),
		int32(i.MessageSection.messagePanel.View.Y),
		int32(i.MessageSection.messagePanel.View.Width),
		int32(i.MessageSection.messagePanel.View.Height),
	)

	for k := range i.MessageSection.messages {
		movingY := i.MessageSection.messages[k].OriginalY + i.MessageSection.messagePanel.Scroll.Y
		rl.DrawRectangle(
			int32(i.MessageSection.messages[k].Bounds.X),
			int32(movingY),
			int32(i.MessageSection.messages[k].Bounds.Width),
			int32(i.MessageSection.messages[k].Bounds.Height),
			rl.SkyBlue)
		rl.DrawText(
			i.MessageSection.messages[k].Content,
			int32(i.MessageSection.messages[k].Bounds.X),
			int32(movingY),
			15,
			rl.White)

	}

	rl.EndScissorMode()

	if i.conversationSection.isConversationSelected {
		i.MessageSection.sendButton.Render()
		i.MessageSection.textInput.Render()
	}

	gui.ScrollPanel(
		i.conversationSection.conversationPanel.Bounds,
		i.conversationSection.nameOnTheWindow,
		i.conversationSection.conversationPanel.Content,
		&i.conversationSection.conversationPanel.Scroll,
		&i.conversationSection.conversationPanel.View,
	)
	rl.BeginScissorMode(
		int32(i.conversationSection.conversationPanel.View.X),
		int32(i.conversationSection.conversationPanel.View.Y),
		int32(i.conversationSection.conversationPanel.View.Width),
		int32(i.conversationSection.conversationPanel.View.Height),
	)
	for k := range i.conversationSection.conversationsTabs {
		movingYTabs := i.conversationSection.conversationsTabs[k].OriginalY + i.conversationSection.conversationPanel.Scroll.Y
		rl.DrawRectangle(
			int32(i.conversationSection.conversationsTabs[k].Bounds.X),
			int32(movingYTabs),
			int32(i.conversationSection.conversationsTabs[k].Bounds.Width),
			int32(i.conversationSection.conversationsTabs[k].Bounds.Height),
			rl.Red)
		rl.DrawText(
			i.conversationSection.conversationsTabs[k].Nametag,
			int32(i.conversationSection.conversationsTabs[k].Bounds.X),
			int32(movingYTabs),
			25,
			rl.Black)
		movingYButtons := movingYTabs

		i.conversationSection.conversationsTabs[k].EnterConversation.Bounds.Y = movingYButtons
		i.conversationSection.conversationsTabs[k].EnterConversation.Render()
	}

	rl.EndScissorMode()

	if i.toolboxSection.showAddConversationModal {
		rl.DrawRectangle(
			int32(i.modalSection.addConversationModal.Background.X),
			int32(i.modalSection.addConversationModal.Background.Y),
			int32(i.modalSection.addConversationModal.Background.Width),
			int32(i.modalSection.addConversationModal.Background.Height),
			i.modalSection.addConversationModal.BgColor)
		if gui.WindowBox(i.modalSection.addConversationModal.Core, "Add conversation") {
			i.toolboxSection.showAddConversationModal = false
			i.modalSection.usersSlider.Strings = i.modalSection.usersSlider.Strings[:0]

		}
		gui.ListViewEx(
			i.modalSection.usersSlider.Bounds,
			i.modalSection.usersSlider.Strings,
			&i.modalSection.usersSlider.IdxScroll,
			&i.modalSection.usersSlider.IdxActiveElement,
			i.modalSection.usersSlider.Focus)

		i.modalSection.acceptAddConversationButton.Render()
		if i.modalSection.isErrorModal {
			rl.DrawRectangle(
				int32(i.modalSection.errorBoxModal.X),
				int32(i.modalSection.errorBoxModal.Y),
				int32(i.modalSection.errorBoxModal.Width),
				int32(i.modalSection.errorBoxModal.Height),
				rl.LightGray)
			rl.DrawText(i.modalSection.textErrorModal, int32(i.modalSection.errorBoxModal.X), int32(i.modalSection.errorBoxModal.Y), 15, rl.Red)
		}
	}

}
