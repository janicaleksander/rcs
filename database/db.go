package database

import (
	"context"

	"github.com/janicaleksander/bcs/proto"
)

// database interface that is used in application
type Storage interface {
	InsertUser(context.Context, *proto.User) error
	LoginUser(ctx context.Context, email, password string) (string, int, error)
	GetUsersWithLVL(ctx context.Context, lvl int) ([]*proto.User, error)
	//TODO maybe in the future also role to this
	InsertUnit(ctx context.Context, nameUnit string, isConfigured bool, id string) error
	GetAllUnits(ctx context.Context) ([]*proto.Unit, error)
	GetUsersInUnit(ctx context.Context, id string) ([]*proto.User, error)
	IsUserInUnit(ctx context.Context, id string) (bool, string, error)
	AssignUserToUnit(ctx context.Context, userID string, unitID string) error
	DeleteUserFromUnit(ctx context.Context, userID string, unitID string) error
	//MESSAGE SERVICE SQL
	IsConversationExists(ctx context.Context, sender, receiver string) (bool, string, error)
	CreateAndAssignConversation(ctx context.Context, cnv *proto.CreateConversationAndAssign) error
	InsertMessage(ctx context.Context, msg *proto.Message) error
	GetUserConversations(ctx context.Context, id string) ([]*proto.ConversationSummary, error)
}
