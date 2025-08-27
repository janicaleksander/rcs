package api

import (
	"encoding/json"
	"errors"
	"fmt"
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

/*
POST req:
{
email:" ",
password: ""
}
*/
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u *proto.LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		render.Render(w, r, ErrInvalidRequestBody(err))
		return
	}
	res, err := utils.MakeRequest(utils.NewRequest(h.ctx, h.serverPID, &proto.HTTPSpawnDevice{
		Email:    u.Email,
		Password: u.Password,
	}))
	if err != nil { //error context deadline
		render.Render(w, r, ErrActorMakeRequest(err))
		return
	}
	var userID string
	var deviceID string
	if v, ok := res.(*proto.SuccessSpawnDevice); !ok {
		render.Render(w, r, ErrInvalidCredentials(errors.New("invalid credentials")))
		return
	} else {
		h.ctx.Send(h.ctx.PID(), &proto.ConnectHDeviceToADevice{
			DeviceID:  v.DeviceID,
			DevicePID: v.DevicePID,
		})
		userID = v.UserID
		deviceID = v.DeviceID

	}

	token, err := token.CreateToken(userID, u.Email)
	if err != nil {
		render.Render(w, r, ErrCreateJWT(err))
		return
	}
	res = &proto.LoginUserRes{
		UserID:      userID,
		DeviceID:    deviceID,
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nice work"))
}

/*
POST req
{
location:{
			 latitude:" "
			 longitude:" "
		 },

deviceID: " "
}
*/

func (h *Handler) updateLocation(w http.ResponseWriter, r *http.Request) {
	var l *proto.UpdateLocationReq
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		return
	}
	fmt.Println(l)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(l)
}
