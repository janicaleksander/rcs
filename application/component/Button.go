package component

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	Bounds  rl.Rectangle
	Text    string
	Focus   bool
	Enabled bool
	Cfg     ButtonConfig
}

type ButtonConfig struct {
	ButtonColor  rl.Color
	TextColor    rl.Color
	TextFontSize int32
}

func defaultButtonConfig() *ButtonConfig {
	return &ButtonConfig{
		ButtonColor:  rl.Blue,
		TextColor:    rl.Black,
		TextFontSize: 12,
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
		config.TextFontSize = fontSize
	}
}
func WithTextColor(color rl.Color) func(*ButtonConfig) {
	return func(config *ButtonConfig) {
		config.TextColor = color
	}

}

func WithButtonColor(color rl.Color) func(*ButtonConfig) {
	return func(config *ButtonConfig) {
		config.ButtonColor = color
	}
}
func NewButton(cfg *ButtonConfig, bounds rl.Rectangle, text string, focus bool) *Button {
	return &Button{
		Bounds:  bounds,
		Text:    text,
		Focus:   focus,
		Enabled: true,
		Cfg:     *cfg,
	}
}

func (b *Button) Update() bool {
	mouse := rl.GetMousePosition()
	hovered := rl.CheckCollisionPointRec(mouse, b.Bounds)
	if !b.Enabled {
		return false
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && hovered {
		b.Focus = true
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		click := b.Focus && hovered
		b.Focus = false
		return click
	}
	if b.Focus && !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		b.Focus = false
	}
	return false
}

func (b *Button) Render() {
	rl.DrawRectangle(
		int32(b.Bounds.X),
		int32(b.Bounds.Y),
		int32(b.Bounds.Width),
		int32(b.Bounds.Height),
		b.Cfg.ButtonColor,
	)
	textWidth := rl.MeasureText(b.Text, b.Cfg.TextFontSize)
	textX := int32(b.Bounds.X + (b.Bounds.Width-float32(textWidth))/2)
	textY := int32(b.Bounds.Y + (b.Bounds.Height-float32(b.Cfg.TextFontSize))/2)
	rl.DrawText(
		b.Text,
		textX,
		textY,
		b.Cfg.TextFontSize,
		b.Cfg.TextColor,
	)
}
func (b *Button) SetActive(bl bool) { b.Enabled = bl }
func (b *Button) Active()           { b.Enabled = true }
func (b *Button) Deactivate() {
	b.Focus = false
	b.Enabled = false
	rl.SetMouseCursor(rl.MouseCursorDefault)
}
