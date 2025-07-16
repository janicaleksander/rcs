package Database

import (
	"context"
	"github.com/janicaleksander/bcs/User"
)

// Database interface that is used in Application
type Storage interface {
	InsertUser(context.Context, User.User) error
	LoginUser(ctx context.Context, email, password string) (string, error)
	//InsertUnit(context.Context, Unit.Unit) error
	//InsertRole(context.Context, string) error
	//InsertDevice(context.Context, Device.Device) error
	//InsertSquad(context.Context, Squad.Squad) error
}
