package component

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/types/proto"
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

type LocationMapInformation struct {
	MapPinInformation map[string]*PinInformation // device id to PinInformation
	MapCurrentTask    map[string]*CurrentTaskTab // deviceID to CurrentTaskTab
}
type PinInformation struct {
	Position       rl.Vector2
	DeviceID       string
	OwnerName      string
	OwnerSurname   string
	LastTimeOnline time.Time
}

type CurrentTaskTab struct {
	OwnerID        string
	OwnerName      string
	OwnerSurname   string
	DeviceID       string
	LastTimeOnline time.Time
	Task           *proto.Task
}

type ToggleGroup struct {
	Selected int
	Labels   []string
	Bounds   []rl.Rectangle
}
