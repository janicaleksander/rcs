package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/janicaleksander/bcs/token"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

type Handler struct {
	ctx        *actor.Context
	serverPID  *actor.PID
	listenAddr string
	router     *chi.Mux
}

func NewHandler(addr string, ctx *actor.Context, serverPID *actor.PID) *Handler {
	return &Handler{
		ctx:        ctx,
		serverPID:  serverPID,
		listenAddr: addr,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u *proto.LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	res, err := utils.MakeRequest(utils.NewRequest(h.ctx, h.serverPID, &proto.HTTPSpawnDevice{
		Email:    u.Email,
		Password: u.Password,
	}))
	if err != nil { //error context deadline
		render.Render(w, r, ErrMakeRequest(err))
		return
	}
	var userID string
	if v, ok := res.(*proto.SuccessSpawnDevice); !ok {
		render.Render(w, r, ErrInvalidRespondMessage(errors.New("type is other than:  proto.SuccessSpawnDevice")))
		return
	} else {
		h.ctx.Send(h.ctx.PID(), &proto.ConnectHDeviceToADevice{
			Id:        v.UserID,
			DevicePID: v.DevicePID,
		})
		userID = v.UserID
	}

	token, err := token.CreateToken(userID, u.Email)
	if err != nil {
		render.Render(w, r, ErrCreateJWT(err))
		return
	}
	res = &proto.LoginUserRes{
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nice work"))
}
