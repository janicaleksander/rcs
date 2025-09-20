package loginstate

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/application/component"
	"github.com/janicaleksander/bcs/application/statesmanager"
	"github.com/janicaleksander/bcs/utils"
)

// LOGIN STATE
type LoginScene struct {
	cfg          *utils.SharedConfig
	stateManager *statesmanager.StateManager
	loginSection LoginSection
	errorSection ErrorSection
}

type LoginSection struct {
	loginButton          component.Button
	emailInput           component.InputBox
	passwordInput        component.InputBox
	isLoginButtonPressed bool
}

type ErrorSection struct {
	errorPopup        component.Popup
	loginErrorMessage string
}

func (l *LoginScene) LoginSceneSetup(state *statesmanager.StateManager, cfg *utils.SharedConfig) {
	l.cfg = cfg
	l.stateManager = state
	l.Reset()

	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	componentWidth := float32(300)
	inputHeight := float32(40)
	buttonHeight := float32(50)

	xPos := screenWidth/2 - componentWidth/2
	yPos := screenHeight/2 - 100

	l.loginSection.emailInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(xPos, yPos, componentWidth, inputHeight),
	)

	yPos += inputHeight + 20
	l.loginSection.passwordInput = *component.NewInputBox(
		component.NewInputBoxConfig(),
		rl.NewRectangle(xPos, yPos, componentWidth, inputHeight),
	)

	yPos += inputHeight + 30
	l.loginSection.loginButton = *component.NewButton(
		component.NewButtonConfig(),
		rl.NewRectangle(xPos, yPos, componentWidth, buttonHeight),
		"LOGIN",
		false,
	)

	l.errorSection.errorPopup = *component.NewPopup(
		component.NewPopupConfig(component.WithBgColor(utils.POPUPERRORBG)),
		rl.NewRectangle(
			xPos,
			yPos+buttonHeight+20,
			componentWidth,
			40,
		),
		&l.errorSection.loginErrorMessage,
	)
}

func (l *LoginScene) UpdateLoginState() {
	l.loginSection.emailInput.Update()
	l.loginSection.passwordInput.Update()
	l.loginSection.isLoginButtonPressed = l.loginSection.loginButton.Update()

	if l.loginSection.isLoginButtonPressed {
		l.Login()
	}
}

func (l *LoginScene) RenderLoginState() {
	rl.ClearBackground(utils.LOGINBGCOLOR)

	titleText := "LOGIN PAGE"
	titleFontSize := int32(80)
	titleWidth := rl.MeasureText(titleText, titleFontSize)
	titleX := int32(rl.GetScreenWidth()/2) - titleWidth/2
	rl.DrawText(titleText, titleX, 80, titleFontSize, rl.DarkGray)

	subtitleText := "remote command system"
	subtitleFontSize := int32(24)
	subtitleWidth := rl.MeasureText(subtitleText, subtitleFontSize)
	subtitleX := int32(rl.GetScreenWidth()/2) - subtitleWidth/2
	rl.DrawText(subtitleText, subtitleX, 170, subtitleFontSize, rl.Gray)

	l.drawInputAreaDecorations()

	l.errorSection.errorPopup.Render()
	l.loginSection.emailInput.Render()
	l.loginSection.passwordInput.Render()
	l.loginSection.loginButton.Render()
}

func (l *LoginScene) drawInputAreaDecorations() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	frameWidth := float32(360)  // Larger than components
	frameHeight := float32(220) // Covers all inputs + button + spacing
	frameX := screenWidth/2 - frameWidth/2
	frameY := screenHeight/2 - 130 // Slightly above first input

	shadowOffset := float32(8)
	shadowColor := rl.NewColor(0, 0, 0, 1)

	rl.DrawRectangle(
		int32(frameX+shadowOffset),
		int32(frameY+shadowOffset),
		int32(frameWidth),
		int32(frameHeight),
		shadowColor,
	)

	outerShadowOffset := float32(12)
	outerShadowColor := rl.NewColor(0, 0, 0, 5)
	rl.DrawRectangle(
		int32(frameX+outerShadowOffset),
		int32(frameY+outerShadowOffset),
		int32(frameWidth),
		int32(frameHeight),
		outerShadowColor,
	)

	bgColor := rl.NewColor(20, 30, 20, 20) // Dark green background
	rl.DrawRectangle(
		int32(frameX),
		int32(frameY),
		int32(frameWidth),
		int32(frameHeight),
		bgColor,
	)

	// Subtle inner glow
	innerGlowColor := rl.NewColor(120, 150, 90, 15)
	innerGlowSize := float32(4)
	rl.DrawRectangle(
		int32(frameX+innerGlowSize),
		int32(frameY+innerGlowSize),
		int32(frameWidth-innerGlowSize*2),
		int32(frameHeight-innerGlowSize*2),
		innerGlowColor,
	)

	// Outer frame (subtle)
	rl.DrawRectangleLinesEx(
		rl.NewRectangle(frameX, frameY, frameWidth, frameHeight),
		2.0,
		rl.NewColor(120, 150, 90, 60), // Semi-transparent green
	)

	// Inner frame (more prominent)
	innerPadding := float32(10)
	rl.DrawRectangleLinesEx(
		rl.NewRectangle(
			frameX+innerPadding,
			frameY+innerPadding,
			frameWidth-innerPadding*2,
			frameHeight-innerPadding*2,
		),
		1.5,
		rl.NewColor(120, 150, 90, 100),
	)

	// Corner accents (military style) with glow effect
	cornerSize := float32(20)
	cornerColor := rl.NewColor(120, 150, 90, 150)
	cornerGlowColor := rl.NewColor(120, 150, 90, 50)

	// Top-left corner with glow
	rl.DrawLineEx(
		rl.NewVector2(frameX-1, frameY+cornerSize+1),
		rl.NewVector2(frameX-1, frameY-1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX-1, frameY-1),
		rl.NewVector2(frameX+cornerSize+1, frameY-1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX, frameY+cornerSize),
		rl.NewVector2(frameX, frameY),
		3.0,
		cornerColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX, frameY),
		rl.NewVector2(frameX+cornerSize, frameY),
		3.0,
		cornerColor,
	)

	// Top-right corner with glow
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+1, frameY+cornerSize+1),
		rl.NewVector2(frameX+frameWidth+1, frameY-1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+1, frameY-1),
		rl.NewVector2(frameX+frameWidth-cornerSize-1, frameY-1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth, frameY+cornerSize),
		rl.NewVector2(frameX+frameWidth, frameY),
		3.0,
		cornerColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth, frameY),
		rl.NewVector2(frameX+frameWidth-cornerSize, frameY),
		3.0,
		cornerColor,
	)

	// Bottom-left corner with glow
	rl.DrawLineEx(
		rl.NewVector2(frameX-1, frameY+frameHeight-cornerSize-1),
		rl.NewVector2(frameX-1, frameY+frameHeight+1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX-1, frameY+frameHeight+1),
		rl.NewVector2(frameX+cornerSize+1, frameY+frameHeight+1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX, frameY+frameHeight-cornerSize),
		rl.NewVector2(frameX, frameY+frameHeight),
		3.0,
		cornerColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX, frameY+frameHeight),
		rl.NewVector2(frameX+cornerSize, frameY+frameHeight),
		3.0,
		cornerColor,
	)

	// Bottom-right corner with glow
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+1, frameY+frameHeight-cornerSize-1),
		rl.NewVector2(frameX+frameWidth+1, frameY+frameHeight+1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+1, frameY+frameHeight+1),
		rl.NewVector2(frameX+frameWidth-cornerSize-1, frameY+frameHeight+1),
		5.0,
		cornerGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth, frameY+frameHeight-cornerSize),
		rl.NewVector2(frameX+frameWidth, frameY+frameHeight),
		3.0,
		cornerColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth, frameY+frameHeight),
		rl.NewVector2(frameX+frameWidth-cornerSize, frameY+frameHeight),
		3.0,
		cornerColor,
	)

	// Enhanced corner dots with glow
	dotSize := float32(4)
	dotColor := rl.NewColor(120, 150, 90, 200)
	dotGlowColor := rl.NewColor(120, 150, 90, 80)

	positions := []rl.Vector2{
		{frameX + 15, frameY + 15},
		{frameX + frameWidth - 15, frameY + 15},
		{frameX + 15, frameY + frameHeight - 15},
		{frameX + frameWidth - 15, frameY + frameHeight - 15},
	}

	for _, pos := range positions {
		rl.DrawCircleV(pos, dotSize+2, dotGlowColor) // Glow
		rl.DrawCircleV(pos, dotSize, dotColor)       // Main dot
	}

	lineLength := float32(30)
	lineColor := rl.NewColor(120, 150, 90, 120)
	lineGlowColor := rl.NewColor(120, 150, 90, 40)

	leftY1 := frameY + frameHeight/3
	leftY2 := frameY + 2*frameHeight/3

	rl.DrawLineEx(
		rl.NewVector2(frameX-6, leftY1),
		rl.NewVector2(frameX-6+lineLength, leftY1),
		4.0,
		lineGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX-5, leftY1),
		rl.NewVector2(frameX-5+lineLength, leftY1),
		2.0,
		lineColor,
	)

	rl.DrawLineEx(
		rl.NewVector2(frameX-6, leftY2),
		rl.NewVector2(frameX-6+lineLength, leftY2),
		4.0,
		lineGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX-5, leftY2),
		rl.NewVector2(frameX-5+lineLength, leftY2),
		2.0,
		lineColor,
	)

	rightY1 := frameY + frameHeight/3
	rightY2 := frameY + 2*frameHeight/3

	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+6-lineLength, rightY1),
		rl.NewVector2(frameX+frameWidth+6, rightY1),
		4.0,
		lineGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+5-lineLength, rightY1),
		rl.NewVector2(frameX+frameWidth+5, rightY1),
		2.0,
		lineColor,
	)

	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+6-lineLength, rightY2),
		rl.NewVector2(frameX+frameWidth+6, rightY2),
		4.0,
		lineGlowColor,
	)
	rl.DrawLineEx(
		rl.NewVector2(frameX+frameWidth+5-lineLength, rightY2),
		rl.NewVector2(frameX+frameWidth+5, rightY2),
		2.0,
		lineColor,
	)
}
