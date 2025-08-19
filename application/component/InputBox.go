package component

import rl "github.com/gen2brain/raylib-go/raylib"

type InputBox struct {
	Bounds  rl.Rectangle
	text    []rune
	length  int32
	enabled bool
	scrollX int32
	//appearance config
	cfg InputBoxConfig
}

type InputBoxConfig struct {
	maxLength  int32
	fontsize   int32
	leftMargin int32
	lineThick  float32
	lineColor  rl.Color
	textColor  rl.Color
}

func defaultInputBoxConfig() *InputBoxConfig {
	return &InputBoxConfig{
		maxLength:  128,
		fontsize:   10,
		leftMargin: 5,
		lineThick:  1,
		lineColor:  rl.Black,
		textColor:  rl.Black,
	}
}
func NewInputBoxConfig(opts ...func(box *InputBoxConfig)) *InputBoxConfig {
	cfg := defaultInputBoxConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
func IBWithFontSize(size int32) func(config *InputBoxConfig) {
	return func(config *InputBoxConfig) {
		config.fontsize = size
	}
}
func IBWithMaxLength(length int32) func(config *InputBoxConfig) {
	return func(config *InputBoxConfig) {
		config.maxLength = length
	}
}

func IBWithLineColor(color rl.Color) func(config *InputBoxConfig) {
	return func(config *InputBoxConfig) {
		config.lineColor = color
	}
}
func IBWithTextColor(color rl.Color) func(config *InputBoxConfig) {
	return func(config *InputBoxConfig) {
		config.textColor = color
	}
}
func checkCharacterInput(key int32) bool {
	if key >= 32 && key <= 125 {
		return true
	}
	switch key {
	case 260, 261, // Ą ą
		262, 263, // Ć ć
		280, 281, // Ę ę
		321, 322, // Ł ł
		323, 324, // Ń ń
		211, 243, // Ó ó
		346, 347, // Ś ś
		377, 378, // Ź ź
		379, 380: // Ż ż
		return true
	}
	return false
}

func NewInputBox(cfg *InputBoxConfig, bounds rl.Rectangle) *InputBox {
	return &InputBox{
		Bounds:  bounds,
		text:    make([]rune, 0, cfg.maxLength),
		enabled: true,
		cfg:     *cfg,
	}
}

func (i *InputBox) Update() {
	mouse := rl.GetMousePosition()
	if i.enabled && rl.CheckCollisionPointRec(mouse, i.Bounds) {
		rl.SetMouseCursor(rl.MouseCursorIBeam)
		for {
			key := rl.GetCharPressed()
			if key == 0 {
				break
			}
			if checkCharacterInput(key) && i.length < i.cfg.maxLength {
				i.text = append(i.text, rune(key))
				i.length++

				textWidth := int32(rl.MeasureText(string(i.text), i.cfg.fontsize))
				if textWidth > int32(i.Bounds.Width)-i.cfg.leftMargin {
					i.scrollX = textWidth + i.cfg.leftMargin - int32(i.Bounds.Width)
				}
			}
		}
		if (rl.IsKeyPressed(rl.KeyBackspace)) && i.length > 0 {
			i.text = i.text[:len(i.text)-1]
			i.length--

			textWidth := int32(rl.MeasureText(string(i.text), i.cfg.fontsize))
			if textWidth <= int32(i.Bounds.Width)-i.cfg.leftMargin {
				i.scrollX = 0
			} else {
				i.scrollX = textWidth + i.cfg.leftMargin - int32(i.Bounds.Width)
			}
		}

	} else {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}
}

func (i *InputBox) Render() {
	rl.DrawRectangleRec(i.Bounds, rl.White)
	rl.DrawRectangleLinesEx(i.Bounds, i.cfg.lineThick, i.cfg.lineColor)

	rl.BeginScissorMode(int32(i.Bounds.X), int32(i.Bounds.Y), int32(i.Bounds.Width), int32(i.Bounds.Height))

	rl.DrawText(string(i.text),
		int32(i.Bounds.X+5)-i.scrollX,
		int32(i.Bounds.Y+8),
		i.cfg.fontsize,
		i.cfg.textColor)
	if i.enabled {
		cursorX := int32(i.Bounds.X) + i.cfg.leftMargin + int32(rl.MeasureText(string(i.text), i.cfg.fontsize)) - i.scrollX
		cursorY1 := int32(i.Bounds.Y + 3)
		cursorY2 := int32(i.Bounds.Y + i.Bounds.Height - 3)

		if int32(rl.GetTime()*2)%2 == 0 {
			rl.DrawRectangle(cursorX, cursorY1, 2, cursorY2, i.cfg.textColor)
		}
	}
	rl.EndScissorMode()
}

func (i *InputBox) GetText() string {
	return string(i.text)
}

func (i *InputBox) Clear() {
	i.text = i.text[:0]
	i.length = 0
	i.scrollX = 0
}

func (i *InputBox) Active() { i.enabled = true }
func (i *InputBox) Deactivate() {
	i.enabled = false
	rl.SetMouseCursor(rl.MouseCursorDefault)
}
func (i *InputBox) SetActive(b bool) { i.enabled = b }
