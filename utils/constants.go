package utils

import (
	"log/slog"
	"os"
	"time"

	"github.com/anthdm/hollywood/actor"
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
