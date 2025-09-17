package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/janicaleksander/bcs/token"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO refactor to use err from actor. now i sue unsuported but
// i cast to Err proto msg ang then read that error
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
	if v, ok := res.(*proto.AcceptSpawnAndRunDevice); !ok {
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
	json.NewEncoder(w).Encode(&res)
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
		render.Render(w, r, ErrInvalidRequestBody(err))
		return
	}
	res, err := utils.MakeRequest(utils.NewRequest(h.ctx, h.ctx.PID(), l))
	if err != nil {
		render.Render(w, r, ErrActorMakeRequest(err))
		return
	}
	if _, ok := res.(*proto.AcceptUpdateLocationReq); !ok {
		render.Render(w, r, ErrUpdateLocation(errors.ErrUnsupported))
		return
	}
	resp := &proto.UpdateLocationRes{
		Message: "Updated",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

/*
GET req /task/{taskID}?deviceID=123
{}
*/

func (h *Handler) userTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	deviceID := r.URL.Query().Get("deviceID")
	if len(strings.TrimSpace(taskID)) == 0 || len(strings.TrimSpace(deviceID)) == 0 {
		render.Render(w, r, ErrBadQueryParam(errors.New("taskID or device param is empty")))
		return
	}
	res, err := utils.MakeRequest(
		utils.NewRequest(
			h.ctx,
			h.ctx.PID(),
			&proto.UserTaskReq{
				DeviceID: deviceID,
				TaskID:   taskID}),
	)
	if err != nil {
		render.Render(w, r, ErrActorMakeRequest(err))
		return
	}
	var resp *proto.UserTaskRes
	if v, ok := res.(*proto.UserTaskRes); !ok {
		render.Render(w, r, ErrGetUserTask(errors.ErrUnsupported))
		return
	} else {
		resp = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

/*
GET req: /tasks/{deviceID}
{}
*/

func (h *Handler) userTasks(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceID")
	if len(strings.TrimSpace(deviceID)) == 0 {
		render.Render(w, r, ErrBadQueryParam(errors.New("userID param is empty")))
		return
	}
	res, err := utils.MakeRequest(
		utils.NewRequest(
			h.ctx,
			h.ctx.PID(),
			&proto.UserTasksReq{
				DeviceID: deviceID},
		))
	if err != nil {
		render.Render(w, r, ErrActorMakeRequest(err))
		return
	}
	var resp *proto.UserTasksRes
	if v, ok := res.(*proto.UserTasksRes); !ok {
		render.Render(w, r, ErrGetUserTask(errors.ErrUnsupported))
		return
	} else {
		resp = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

// task/current/{userID}?deviceID=123
/*
taskID: ""
*/
func (h *Handler) updateCurrentTask(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	deviceID := r.URL.Query().Get("deviceID")
	reqBody := struct {
		TaskID string `json:"taskID"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		render.Render(w, r, ErrInvalidRequestBody(err))
		return
	}
	if len(strings.TrimSpace(userID)) == 0 || len(strings.TrimSpace(deviceID)) == 0 {
		render.Render(w, r, ErrBadQueryParam(errors.New("params are empty")))
		return
	}
	res, err := utils.MakeRequest(utils.NewRequest(h.ctx, h.ctx.PID(), &proto.UpdateCurrentTaskReq{
		DeviceID: deviceID,
		TaskID:   reqBody.TaskID,
		UserID:   userID,
	}))
	if err != nil {
		render.Render(w, r, ErrActorMakeRequest(errors.ErrUnsupported))
		return
	}
	if _, ok := res.(*proto.UpdateCurrentTaskRes); !ok {
		render.Render(w, r, ErrUpdateCurrentTask(errors.New("bad actor response")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)

}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	deviceID := r.URL.Query().Get("deviceID")
	if len(strings.TrimSpace(taskID)) == 0 || len(strings.TrimSpace(deviceID)) == 0 {
		render.Render(w, r, ErrBadQueryParam(errors.New("params are empty")))
		return
	}
	res, err := utils.MakeRequest(utils.NewRequest(h.ctx, h.ctx.PID(), &proto.DeleteTaskReq{
		DeviceID: deviceID,
		TaskID:   taskID,
	}))
	if err != nil {
		render.Render(w, r, ErrActorMakeRequest(errors.ErrUnsupported))
		return
	}
	if _, ok := res.(*proto.DeleteTaskRes); !ok {
		render.Render(w, r, ErrDeleteTask(errors.New("bad actor response")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}
