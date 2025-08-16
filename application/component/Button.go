package component

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	text   string
	bounds rl.Rectangle
	focus  bool
	cfg    ButtonConfig
}

type ButtonConfig struct {
	buttonColor  rl.Color
	textColor    rl.Color
	textFontSize int32
}

func defaultButtonConfig() *ButtonConfig {
	return &ButtonConfig{
		buttonColor:  rl.Blue,
		textColor:    rl.Black,
		textFontSize: 12,
	}
}
func NewButtonConfig(opts ...func(box *ButtonConfig)) *ButtonConfig {
	cfg := defaultButtonConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
func WithFontSize(fontSize int32) func(*ButtonConfig) {
	return func(config *ButtonConfig) {
		config.textFontSize = fontSize
	}
}
func WithTextColor(color rl.Color) func(*ButtonConfig) {
	return func(config *ButtonConfig) {
		config.textColor = color
	}

}

func WithButtonColor(color rl.Color) func(*ButtonConfig) {
	return func(config *ButtonConfig) {
		config.buttonColor = color
	}
}
func NewButton(cfg *ButtonConfig, bounds rl.Rectangle, text string, focus bool) *Button {
	return &Button{
		text:   text,
		bounds: bounds,
		focus:  focus,
		cfg:    *cfg,
	}
}

func (b *Button) Update() bool {
	if b.focus {
		active := false
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, b.bounds) {
			if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
				active = true
			}
		}
		return active
	}
	return false
}

func (b *Button) Render() {
	rl.DrawRectangle(
		int32(b.bounds.X),
		int32(b.bounds.Y),
		int32(b.bounds.Width),
		int32(b.bounds.Height),
		b.cfg.buttonColor,
	)
	textWidth := rl.MeasureText(b.text, b.cfg.textFontSize)
	textX := int32(b.bounds.X + (b.bounds.Width-float32(textWidth))/2)
	textY := int32(b.bounds.Y + (b.bounds.Height-float32(b.cfg.textFontSize))/2)
	rl.DrawText(
		b.text,
		textX,
		textY,
		b.cfg.textFontSize,
		b.cfg.textColor,
	)
}

func (b *Button) Active() { b.focus = true }
func (b *Button) Deactivate() {
	b.focus = false
	rl.SetMouseCursor(rl.MouseCursorDefault)
}
