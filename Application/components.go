package Application

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	position rl.Rectangle
	text     string
}

type ListSlider struct {
	strings          []string
	bounds           rl.Rectangle
	idxActiveElement int32
	focus            int32
	idxScroll        int32
}
