package utils

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/janicaleksander/bcs/customerror"
)

//COLORS

var (
	POPUPERRORBG   = rl.NewColor(250, 120, 120, 255)
	POPUPINFOBG    = rl.NewColor(141, 235, 166, 255)
	LOGINBGCOLOR   = rl.NewColor(235, 237, 216, 255)
	HCPARTSBG      = rl.NewColor(207, 209, 190, 255)
	HCMENUBG       = rl.NewColor(235, 237, 216, 255)
	CREATEUNITBG   = rl.NewColor(235, 237, 216, 255)
	CREATEUSERBG   = rl.NewColor(235, 237, 216, 255)
	USERDESCBG     = rl.NewColor(235, 237, 216, 255)
	USERBUTTONSBG  = rl.NewColor(207, 209, 190, 255)
	CREATEDEVICEBG = rl.NewColor(235, 237, 216, 255)
)
var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
var STARTPLACE = "WRO"

const (
	WaitTime = 3 * time.Second
)

type Request struct {
	ctx *actor.Context
	pid *actor.PID
	msg any
}

func NewRequest(ctx *actor.Context, pid *actor.PID, msg any) *Request {
	return &Request{
		ctx: ctx,
		pid: pid,
		msg: msg,
	}
}

func MakeRequest(r *Request) (any, error) {
	resp := r.ctx.Request(r.pid, r.msg, WaitTime)
	res, err := resp.Result()
	if err != nil {
		return nil, customerror.ContextDeadlineExceed
	}
	return res, nil
}

func WrapText(maxWidth int32, input string, fontSize int32) string {
	var output strings.Builder
	var line strings.Builder

	words := strings.Fields(input)

	for i, word := range words {
		testLine := line.String()
		if line.Len() > 0 {
			testLine += " " + word
		} else {
			testLine = word
		}

		width := rl.MeasureText(testLine, fontSize)
		if width > maxWidth && line.Len() > 0 {
			output.WriteString(line.String())
			output.WriteString("\n")
			line.Reset()
			line.WriteString(word)
		} else {
			if line.Len() > 0 {
				line.WriteString(" ")
			}
			line.WriteString(word)
		}

		if i == len(words)-1 {
			output.WriteString(line.String())
		}
	}
	return output.String()
}

type SharedConfig struct {
	ServerPID         *actor.PID
	MessageServicePID *actor.PID
	Ctx               *actor.Context
}
