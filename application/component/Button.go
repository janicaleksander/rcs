package component

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	Bounds  rl.Rectangle
	text    string
	focus   bool
	enabled bool
	cfg     ButtonConfig
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
		Bounds:  bounds,
		text:    text,
		focus:   focus,
		enabled: true,
		cfg:     *cfg,
	}
}

func (b *Button) Update() bool {
	mouse := rl.GetMousePosition()
	hovered := rl.CheckCollisionPointRec(mouse, b.Bounds)
	if !b.enabled {
		return false
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && hovered {
		b.focus = true
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		click := b.focus && hovered
		b.focus = false
		return click
	}
	if b.focus && !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		b.focus = false
	}
	return false
}

func (b *Button) Render() {
	rl.DrawRectangle(
		int32(b.Bounds.X),
		int32(b.Bounds.Y),
		int32(b.Bounds.Width),
		int32(b.Bounds.Height),
		b.cfg.buttonColor,
	)
	textWidth := rl.MeasureText(b.text, b.cfg.textFontSize)
	textX := int32(b.Bounds.X + (b.Bounds.Width-float32(textWidth))/2)
	textY := int32(b.Bounds.Y + (b.Bounds.Height-float32(b.cfg.textFontSize))/2)
	rl.DrawText(
		b.text,
		textX,
		textY,
		b.cfg.textFontSize,
		b.cfg.textColor,
	)
}
func (b *Button) SetActive(bl bool) { b.enabled = bl }
func (b *Button) Active()           { b.enabled = true }
func (b *Button) Deactivate() {
	b.focus = false
	b.enabled = false
	rl.SetMouseCursor(rl.MouseCursorDefault)
}
