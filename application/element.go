package application

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Circle struct {
	x      int32
	y      int32
	radius float32
	color  rl.Color
}
type Button struct {
	bounds rl.Rectangle
	text   string
}

type ListSlider struct {
	strings          []string
	bounds           rl.Rectangle
	idxActiveElement int32
	focus            int32
	idxScroll        int32
}

type Modal struct {
	background rl.Rectangle
	bgColor    rl.Color
	core       rl.Rectangle
}

type ConversationTab struct {
	bounds  rl.Rectangle
	nametag string
	//lastMessage       string TODO now we only have 2 states one outside conversation and one in conversation
	enterConversation Button
	isClicked         bool
	conversationID    string // to remove in the future, when i made own scrollbar i will use tab[currIDX]
	withID            string // to remove in the future, when i made own scrollbar i will use tab[currIDX]
}

type Message struct {
	bounds    rl.Rectangle
	content   string
	originalY float32
}

//TODO add active status to users

type ScrollPanel struct {
	bounds  rl.Rectangle
	content rl.Rectangle
	scroll  rl.Vector2
	view    rl.Rectangle
}
