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

type InputField struct {
	bounds   rl.Rectangle
	text     string
	focus    bool
	textSize int
}

type Modal struct {
	background rl.Rectangle
	bgColor    rl.Color
	core       rl.Rectangle
}

//TODO add active status to users
