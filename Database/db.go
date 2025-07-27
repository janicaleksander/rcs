package Database

import (
	"context"
	"github.com/janicaleksander/bcs/Proto"
)

// Database interface that is used in Application
type Storage interface {
	InsertUser(context.Context, *Proto.User) error
	LoginUser(ctx context.Context, email, password string) (string, int, error)
	GetUsersWithLVL(ctx context.Context, lvl int) ([]*Proto.User, error)
	//TODO maybe in the future also role to this
	InsertUnit(ctx context.Context, nameUnit string, isConfigured bool, id string) error

	GetAllUnits(ctx context.Context) ([]*Proto.Unit, error)
	GetUsersInUnit(ctx context.Context, id string) ([]*Proto.User, error)
	//InsertUnit(context.Context, Unit.Unit) error
	//InsertRole(context.Context, string) error
	//InsertDevice(context.Context, Device.Device) error
	//InsertSquad(context.Context, Squad.Squad) error
}
