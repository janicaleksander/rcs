package component

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Circle struct {
	X      int32
	Y      int32
	Radius float32
	Color  rl.Color
}
type ListSlider struct {
	Strings          []string
	Bounds           rl.Rectangle
	IdxActiveElement int32
	Focus            int32
	IdxScroll        int32
}

type Modal struct {
	Background rl.Rectangle
	BgColor    rl.Color
	Core       rl.Rectangle
}

type ConversationTab struct {
	ID      int32
	Bounds  rl.Rectangle
	Nametag string
	//lastMessage       string TODO now we only have 2 states one outside conversation and one in conversation
	EnterConversation Button
	IsPressed         bool
	ConversationID    string
	WithID            string
	OriginalY         float32
}

type Message struct {
	Bounds    rl.Rectangle
	Content   string
	OriginalY float32
}

//TODO add active status to users

type ScrollPanel struct {
	Bounds  rl.Rectangle
	Content rl.Rectangle
	Scroll  rl.Vector2
	View    rl.Rectangle
}

type PinInfo struct {
	Position rl.Vector2
	Owner    string
	Time     time.Time
}
