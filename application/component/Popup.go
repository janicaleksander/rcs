package component

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/utils"
)

type Popup struct {
	show   bool
	text   *string
	bounds rl.Rectangle
	cfg    PopupConfig
}

type PopupConfig struct {
	bgColor   rl.Color
	textColor rl.Color
	fontSize  int32
}

func defaultPopupConfig() *PopupConfig {
	return &PopupConfig{
		bgColor:   rl.Red,
		textColor: rl.White,
		fontSize:  15,
	}
}
func NewPopupConfig(opts ...func(*PopupConfig)) *PopupConfig {
	cfg := defaultPopupConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func NewPopup(cfg *PopupConfig, bounds rl.Rectangle, text *string) *Popup {
	return &Popup{
		show:   false,
		text:   text,
		bounds: bounds,
		cfg:    *cfg,
	}

}
func (p *Popup) Show() { p.show = true }
func (p *Popup) Hide() { p.show = false }
func (p *Popup) Render() {
	if p.show {
		rl.DrawRectangle(
			int32(p.bounds.X),
			int32(p.bounds.Y),
			int32(p.bounds.Width),
			int32(p.bounds.Height),
			p.cfg.bgColor)
		rl.DrawText(utils.WrapText(int32(p.bounds.Width), *p.text, p.cfg.fontSize), int32(p.bounds.X), int32(p.bounds.Y), p.cfg.fontSize, p.cfg.textColor)
	}
}

func (p *Popup) GetText() string {
	return *p.text
}
