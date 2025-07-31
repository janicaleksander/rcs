package Application

import "C"
import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

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

type InputField struct {
	bounds   rl.Rectangle
	text     string
	focus    bool
	textSize int
}

//TODO add active status to users
