package application

import (
	"sync"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/proto"
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
	addConversationModal           Modal
	usersSlider                    ListSlider       // user slider inside modals
	users                          []*proto.User    // users from DB without logged in user
	acceptAddConversationButton    component.Button // button inside modal to confirm
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

	messageChan        chan *proto.Message //chan to transport messages from window actor
	messages           []Message           // messages for current conversation
	messagePanel       ScrollPanel         //slider with all messages
	messagePanelLayout messagePanelLayout  // messagePanel configuration
	nameOnTheWindow    string
}
type ConversationSection struct {
	conversationArea        rl.Rectangle // area where messages display
	isConversationSelected  bool
	usersConversations      []*proto.ConversationSummary
	conversationPanelLayout conversationPanelLayout
	conversationsTabs       []ConversationTab
	activeConversation      int32
	activeConversationID    string
	activeWithID            string
}

type InboxScene struct {
	tempUserID          string // id of current logged in a user
	toolboxSection      ToolboxSection
	modalSection        ModalSection
	messageSection      MessageSection
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
func (w *Window) setupInboxScene() {
	w.inboxScene.Reset()

	w.GetLoggedID()
	w.GetUserConversation()
	w.GetUserToNewConversation()

	w.inboxScene.toolboxSection.toolboxArea = rl.NewRectangle(
		0,
		0,
		(2.0/5.0)*float32(w.width),
		(1.0/8.0)*float32(w.height))

	//go back button
	var toolboxButtonPadding float32 = 10
	var toolBoxButtonWidth float32 = 150
	var toolBoxButtonHeight float32 = 50

	w.inboxScene.toolboxSection.backButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			w.inboxScene.toolboxSection.toolboxArea.X+toolboxButtonPadding,
			w.inboxScene.toolboxSection.toolboxArea.Y+toolboxButtonPadding,
			toolBoxButtonWidth,
			toolBoxButtonHeight),
		"GO BACK",
		false)

	w.inboxScene.toolboxSection.addConversationButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			w.inboxScene.toolboxSection.backButton.Bounds.X+w.inboxScene.toolboxSection.backButton.Bounds.Width+toolboxButtonPadding,
			w.inboxScene.toolboxSection.toolboxArea.Y+toolboxButtonPadding,
			toolBoxButtonWidth,
			toolBoxButtonHeight),
		"Add \n conversation",
		false)

	w.inboxScene.modalSection.addConversationModal = Modal{
		background: rl.NewRectangle(
			0,
			0,
			float32(w.width),
			float32(w.height)),
		bgColor: rl.Fade(rl.LightGray, 0.2),
		core:    rl.NewRectangle(float32(w.width/2-150.0), float32(w.height/2-150.0), 300, 280),
	}
	var addConversationModalPadding float32 = 20
	w.inboxScene.modalSection.usersSlider = ListSlider{
		strings: make([]string, 0, 64),
		bounds: rl.NewRectangle(
			w.inboxScene.modalSection.addConversationModal.core.X+w.inboxScene.modalSection.addConversationModal.core.Width/2-(w.inboxScene.modalSection.addConversationModal.core.Width/2)/2,
			w.inboxScene.modalSection.addConversationModal.core.Y+2*addConversationModalPadding,
			w.inboxScene.modalSection.addConversationModal.core.Width/2,
			80),
		idxActiveElement: 0,
		focus:            0,
		idxScroll:        0,
	}
	//fill above slice with users
	for _, user := range w.inboxScene.modalSection.users {
		w.inboxScene.modalSection.usersSlider.strings = append(w.inboxScene.modalSection.usersSlider.strings,
			user.Email+"\n"+user.Personal.Name+user.Personal.Name)
	}

	w.inboxScene.modalSection.acceptAddConversationButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			w.inboxScene.modalSection.usersSlider.bounds.X,
			w.inboxScene.modalSection.usersSlider.bounds.Y+5*addConversationModalPadding,
			w.inboxScene.modalSection.usersSlider.bounds.Width,
			25),
		"Add conversation", false)

	w.inboxScene.modalSection.errorBoxModal = rl.NewRectangle(
		w.inboxScene.modalSection.usersSlider.bounds.X,
		w.inboxScene.modalSection.usersSlider.bounds.Y+8*addConversationModalPadding,
		w.inboxScene.modalSection.usersSlider.bounds.Width,
		75)

	w.inboxScene.toolboxSection.refreshConversationsButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			w.inboxScene.toolboxSection.addConversationButton.Bounds.X+w.inboxScene.toolboxSection.addConversationButton.Bounds.Width+toolboxButtonPadding,
			w.inboxScene.toolboxSection.toolboxArea.Y+toolboxButtonPadding,
			toolBoxButtonWidth,
			toolBoxButtonHeight),
		"Fetch \n conversations", false)

	w.inboxScene.conversationSection.conversationArea = rl.NewRectangle(
		w.inboxScene.toolboxSection.toolboxArea.X,
		w.inboxScene.toolboxSection.toolboxArea.Height,
		w.inboxScene.toolboxSection.toolboxArea.Width,
		(7.0/8.0)*float32(w.height))

	w.inboxScene.messageSection.messagePanel.bounds = rl.NewRectangle(
		w.inboxScene.toolboxSection.toolboxArea.Width,
		w.inboxScene.toolboxSection.toolboxArea.Y,
		(3.0/5.0)*float32(w.width),
		float32(w.height))
	w.inboxScene.messageSection.messagePanel.content = rl.NewRectangle(
		w.inboxScene.toolboxSection.toolboxArea.Width+5,
		w.inboxScene.toolboxSection.toolboxArea.Y+5,
		(3.0/5.0)*float32(w.width)-15,
		float32(w.height)*10)

	w.inboxScene.messageSection.nameOnTheWindow = "MESSAGES"
	w.inboxScene.messageSection.messagePanel.view = rl.Rectangle{}
	w.inboxScene.messageSection.messagePanel.scroll = rl.Vector2{}

	w.inboxScene.messageSection.messagePanelLayout.padding = 20.0
	w.inboxScene.messageSection.messagePanelLayout.middle = w.inboxScene.messageSection.messagePanel.bounds.X + (w.inboxScene.messageSection.messagePanel.bounds.Width)/2.0
	w.inboxScene.messageSection.messagePanelLayout.currHeight = w.inboxScene.messageSection.messagePanel.bounds.Y + 3*w.inboxScene.messageSection.messagePanelLayout.padding
	w.inboxScene.messageSection.messagePanelLayout.messageWidth = 150
	w.inboxScene.messageSection.messagePanelLayout.messageFontSize = 20
	w.inboxScene.messageSection.messagePanelLayout.leftSide = w.inboxScene.messageSection.messagePanel.bounds.X + w.inboxScene.messageSection.messagePanelLayout.padding
	w.inboxScene.messageSection.messagePanelLayout.rightSide = w.inboxScene.messageSection.messagePanelLayout.middle + w.inboxScene.messageSection.messagePanelLayout.padding

	w.inboxScene.messageSection.textInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(
			w.inboxScene.messageSection.messagePanel.bounds.X,
			w.inboxScene.messageSection.messagePanel.bounds.Height-30, //height
			w.inboxScene.messageSection.messagePanel.bounds.Width-70,  //for button// TODO vars
			30),
		false)

	w.inboxScene.messageSection.sendButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(
			w.inboxScene.messageSection.messagePanel.bounds.X+w.inboxScene.messageSection.textInput.Bounds.Width,
			w.inboxScene.messageSection.messagePanel.bounds.Height-30,
			65,
			30),
		"SEND",
		false)
	w.inboxScene.conversationSection.conversationPanelLayout = conversationPanelLayout{
		currHeight:  80,
		panelHeight: 80,
		startHeight: 100,
	}
	w.refreshConversationsPanel()
	w.inboxScene.messageSection.messageChan = make(chan *proto.Message, 1024)
	go func() {
		for msg := range w.inboxScene.messageSection.messageChan {
			w.AppendMessage(msg)
		}
	}()

}

func (w *Window) updateInboxState() {
	if w.inboxScene.toolboxSection.isGoBackButtonPressed {
		w.goSceneBack()
	}

	w.inboxScene.messageSection.textInput.Update()
	w.inboxScene.toolboxSection.isGoBackButtonPressed = w.inboxScene.toolboxSection.backButton.Update()
	if w.inboxScene.toolboxSection.addConversationButton.Update() {
		// we need this if to remain windowbox on the screen
		//in other cases we can use one line because e.g we send sth so we need this true state for short time (moment)
		w.inboxScene.toolboxSection.showAddConversationModal = true
	}
	w.inboxScene.toolboxSection.isRefreshConversationPressed = w.inboxScene.toolboxSection.refreshConversationsButton.Update()
	w.inboxScene.modalSection.isAcceptAddConversationPressed = w.inboxScene.modalSection.acceptAddConversationButton.Update()
	w.inboxScene.messageSection.isSendButtonPressed = w.inboxScene.messageSection.sendButton.Update()
	for i := range w.inboxScene.conversationSection.conversationsTabs {
		if w.inboxScene.conversationSection.conversationsTabs[i].enterConversation.Update() {
			w.inboxScene.conversationSection.conversationsTabs[i].isPressed = true
		}
	}
	w.inboxScene.conversationSection.activeConversation = -1
	for i, tab := range w.inboxScene.conversationSection.conversationsTabs {
		if tab.isPressed {
			w.inboxScene.messageSection.nameOnTheWindow = "MESSAGES WITH " + w.inboxScene.conversationSection.usersConversations[tab.ID].Nametag

			//we have to mark this to know if we have to show input box with send button
			w.inboxScene.conversationSection.isConversationSelected = true
			w.inboxScene.conversationSection.activeWithID = tab.withID
			w.inboxScene.conversationSection.activeConversationID = tab.conversationID
			if tab.ID != w.inboxScene.conversationSection.activeConversation {
				w.inboxScene.conversationSection.activeConversation = tab.ID

				//open conversation
				w.ctx.Send(w.messageServicePID, &proto.UpdatePresence{
					Id: w.inboxScene.tempUserID,
					Presence: &proto.PresenceType{
						Type: &proto.PresenceType_Inbox{
							Inbox: &proto.Inbox{
								WithID: tab.withID}}},
				})
				res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.OpenAndLoadConversation{
					UserID:         w.inboxScene.tempUserID,
					ReceiverID:     tab.withID,
					ConversationID: tab.conversationID,
				},
				))
				if err != nil {
					//TODO context deadline exceeded
				}
				w.inboxScene.messageSection.messages = w.inboxScene.messageSection.messages[:0]
				w.inboxScene.messageSection.messagePanelLayout.currHeight = w.inboxScene.messageSection.messagePanel.bounds.Y + 3*w.inboxScene.messageSection.messagePanelLayout.padding
				if v, ok := res.(*proto.SuccessOpenAndLoadConversation); ok {
					for _, msg := range v.Messages {
						w.AppendMessage(msg)
					}
				} else {
					//TODO err
				}
				w.inboxScene.conversationSection.conversationsTabs[i].isPressed = false // to not load every time
			}
		}

	}
	if w.inboxScene.messageSection.isSendButtonPressed {
		res, err := utils.MakeRequest(utils.NewRequest(w.ctx, w.messageServicePID, &proto.SendMessage{
			Receiver: w.inboxScene.conversationSection.activeWithID,
			Message: &proto.Message{
				Id:             uuid.New().String(),
				SenderID:       w.inboxScene.tempUserID,
				ConversationID: w.inboxScene.conversationSection.activeConversationID,
				Content:        w.inboxScene.messageSection.textInput.GetText(),
				SentAt:         timestamppb.Now(),
			},
		}))
		if err != nil {
			// context deadline exceed
		}

		if v, ok := res.(*proto.DeliverMessage); ok {
			w.AppendMessage(v.Message)
		} else {
			//error
		}

	}

	if w.inboxScene.modalSection.isAcceptAddConversationPressed {
		w.addNewConversation()
	}
	if w.inboxScene.toolboxSection.isRefreshConversationPressed {
		w.refreshConversationsPanel()
	}

}

func (w *Window) renderInboxState() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		w.inboxScene.toolboxSection.backButton.Deactivate()
		w.inboxScene.toolboxSection.addConversationButton.Deactivate()
		w.inboxScene.toolboxSection.refreshConversationsButton.Deactivate()
		w.inboxScene.modalSection.acceptAddConversationButton.Deactivate()
		w.inboxScene.messageSection.sendButton.Deactivate()

		if rl.CheckCollisionPointRec(mousePos, w.inboxScene.toolboxSection.backButton.Bounds) && !w.inboxScene.toolboxSection.showAddConversationModal {
			w.inboxScene.toolboxSection.backButton.Active()
		}
		if rl.CheckCollisionPointRec(mousePos, w.inboxScene.toolboxSection.addConversationButton.Bounds) && !w.inboxScene.toolboxSection.showAddConversationModal {
			w.inboxScene.toolboxSection.addConversationButton.Active()
		}
		if rl.CheckCollisionPointRec(mousePos, w.inboxScene.toolboxSection.refreshConversationsButton.Bounds) && !w.inboxScene.toolboxSection.showAddConversationModal {
			w.inboxScene.toolboxSection.refreshConversationsButton.Active()
		}
		if rl.CheckCollisionPointRec(mousePos, w.inboxScene.modalSection.acceptAddConversationButton.Bounds) {
			w.inboxScene.modalSection.acceptAddConversationButton.Active()
		}
		if rl.CheckCollisionPointRec(mousePos, w.inboxScene.messageSection.sendButton.Bounds) && !w.inboxScene.toolboxSection.showAddConversationModal {
			w.inboxScene.messageSection.sendButton.Active()
		}
		if w.inboxScene.toolboxSection.showAddConversationModal {
			for i := range w.inboxScene.conversationSection.conversationsTabs {
				w.inboxScene.conversationSection.conversationsTabs[i].enterConversation.Deactivate()
			}
			w.inboxScene.messageSection.textInput.Deactivate()
		} else {
			w.inboxScene.messageSection.textInput.Active()
			for i := range w.inboxScene.conversationSection.conversationsTabs {
				w.inboxScene.conversationSection.conversationsTabs[i].enterConversation.Active()
			}
		}

	}
	//toolbox
	rl.DrawRectangle(
		int32(w.inboxScene.toolboxSection.toolboxArea.X),
		int32(w.inboxScene.toolboxSection.toolboxArea.Y),
		int32(w.inboxScene.toolboxSection.toolboxArea.Width),
		int32(w.inboxScene.toolboxSection.toolboxArea.Height),
		rl.Gray)

	//toolbox button section
	w.inboxScene.toolboxSection.backButton.Render()
	w.inboxScene.toolboxSection.addConversationButton.Render()
	w.inboxScene.toolboxSection.refreshConversationsButton.Render()

	rl.DrawRectangle(
		int32(w.inboxScene.conversationSection.conversationArea.X),
		int32(w.inboxScene.conversationSection.conversationArea.Y),
		int32(w.inboxScene.conversationSection.conversationArea.Width),
		int32(w.inboxScene.conversationSection.conversationArea.Height),
		rl.White)

	gui.ScrollPanel(
		w.inboxScene.messageSection.messagePanel.bounds,
		w.inboxScene.messageSection.nameOnTheWindow,
		w.inboxScene.messageSection.messagePanel.content,
		&w.inboxScene.messageSection.messagePanel.scroll,
		&w.inboxScene.messageSection.messagePanel.view,
	)
	rl.BeginScissorMode(
		int32(w.inboxScene.messageSection.messagePanel.view.X),
		int32(w.inboxScene.messageSection.messagePanel.view.Y),
		int32(w.inboxScene.messageSection.messagePanel.view.Width),
		int32(w.inboxScene.messageSection.messagePanel.view.Height),
	)

	for i := range w.inboxScene.messageSection.messages {
		movingY := w.inboxScene.messageSection.messages[i].originalY + w.inboxScene.messageSection.messagePanel.scroll.Y
		rl.DrawRectangle(
			int32(w.inboxScene.messageSection.messages[i].bounds.X),
			int32(movingY),
			int32(w.inboxScene.messageSection.messages[i].bounds.Width),
			int32(w.inboxScene.messageSection.messages[i].bounds.Height),
			rl.SkyBlue)
		rl.DrawText(
			w.inboxScene.messageSection.messages[i].content,
			int32(w.inboxScene.messageSection.messages[i].bounds.X),
			int32(movingY),
			15,
			rl.White)

	}

	rl.EndScissorMode()

	if w.inboxScene.conversationSection.isConversationSelected {
		w.inboxScene.messageSection.sendButton.Render()
		w.inboxScene.messageSection.textInput.Render()
	}
	for i := range w.inboxScene.conversationSection.conversationsTabs {
		rl.DrawRectangle(
			int32(w.inboxScene.conversationSection.conversationsTabs[i].bounds.X),
			int32(w.inboxScene.conversationSection.conversationsTabs[i].bounds.Y),
			int32(w.inboxScene.conversationSection.conversationsTabs[i].bounds.Width),
			int32(w.inboxScene.conversationSection.conversationsTabs[i].bounds.Height),
			rl.Red)
		rl.DrawText(
			w.inboxScene.conversationSection.conversationsTabs[i].nametag,
			int32(w.inboxScene.conversationSection.conversationsTabs[i].bounds.X),
			int32(w.inboxScene.conversationSection.conversationsTabs[i].bounds.Y),
			25,
			rl.Black)
		w.inboxScene.conversationSection.conversationsTabs[i].enterConversation.Render()
	}

	if w.inboxScene.toolboxSection.showAddConversationModal {
		rl.DrawRectangle(
			int32(w.inboxScene.modalSection.addConversationModal.background.X),
			int32(w.inboxScene.modalSection.addConversationModal.background.Y),
			int32(w.inboxScene.modalSection.addConversationModal.background.Width),
			int32(w.inboxScene.modalSection.addConversationModal.background.Height),
			w.inboxScene.modalSection.addConversationModal.bgColor)
		if gui.WindowBox(w.inboxScene.modalSection.addConversationModal.core, "Add conversation") {
			w.inboxScene.toolboxSection.showAddConversationModal = false
			w.inboxScene.modalSection.usersSlider.strings = w.inboxScene.modalSection.usersSlider.strings[:0]

		}
		gui.ListViewEx(
			w.inboxScene.modalSection.usersSlider.bounds,
			w.inboxScene.modalSection.usersSlider.strings,
			&w.inboxScene.modalSection.usersSlider.idxScroll,
			&w.inboxScene.modalSection.usersSlider.idxActiveElement,
			w.inboxScene.modalSection.usersSlider.focus)

		w.inboxScene.modalSection.acceptAddConversationButton.Render()
		if w.inboxScene.modalSection.isErrorModal {
			rl.DrawRectangle(
				int32(w.inboxScene.modalSection.errorBoxModal.X),
				int32(w.inboxScene.modalSection.errorBoxModal.Y),
				int32(w.inboxScene.modalSection.errorBoxModal.Width),
				int32(w.inboxScene.modalSection.errorBoxModal.Height),
				rl.LightGray)
			rl.DrawText(w.inboxScene.modalSection.textErrorModal, int32(w.inboxScene.modalSection.errorBoxModal.X), int32(w.inboxScene.modalSection.errorBoxModal.Y), 15, rl.Red)
		}
	}

}
