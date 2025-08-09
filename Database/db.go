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
	IsUserInUnit(ctx context.Context, id string) (bool, string, error)
	AssignUserToUnit(ctx context.Context, userID string, unitID string) error
	DeleteUserFromUnit(ctx context.Context, userID string, unitID string) error
	//MESSAGE SERVICE SQL
	IsConversationExists(ctx context.Context, sender, receiver string) (bool, string, error)
	CreateAndAssignConversation(ctx context.Context, cnv *Proto.CreateConversationAndAssign) error
	InsertMessage(ctx context.Context, msg *Proto.Message) error
	GetUserConversations(ctx context.Context, id string) ([]*Proto.ConversationSummary, error)
}
