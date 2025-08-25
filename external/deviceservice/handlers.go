package deviceservice

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/janicaleksander/bcs/token"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func (d *DeviceHTTP) Login(w http.ResponseWriter, r *http.Request) {
	var u *proto.LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	res, err := utils.MakeRequest(utils.NewRequest(d.ctx, d.serverPID, &proto.HTTPSpawnDevice{
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
		d.ctx.Send(d.ctx.PID(), &proto.ConnectHDeviceToADevice{
			Id:        v.UserID,
			DevicePID: v.DevicePID,
		})
		return
	} else {
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

func (d *DeviceHTTP) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nice work"))
}
