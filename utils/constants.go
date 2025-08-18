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

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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
	for _, char := range input {
		line.WriteString(string(char))
		width := rl.MeasureText(line.String(), fontSize)
		if width >= maxWidth {
			output.WriteString("\n")
			line.Reset()
		}
		output.WriteString(string(char))
	}

	return output.String()
}

type Task struct {
	remaining float64
	handler   func()
}

type Scheduler struct {
	tasks []Task
}

func (s *Scheduler) After(seconds float64, fn func()) {
	s.tasks = append(s.tasks, Task{remaining: seconds, handler: fn})
}
func (s *Scheduler) Update(dt float64) {
	remaining := s.tasks[:0]
	for _, t := range s.tasks {
		t.remaining -= dt
		if t.remaining <= 0 {
			t.handler()
		} else {
			remaining = append(remaining, t)
		}
	}
	s.tasks = remaining
}
