package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/janicaleksander/bcs/types/proto"
	"github.com/joho/godotenv"
)

type Config struct {
	url string
}
type DBManager struct {
	db   *sql.DB
	once sync.Once
}

var dbManager *DBManager

func GetDBManager(options ...func(*string)) (*DBManager, error) {
	if dbManager == nil {
		err := godotenv.Load()
		if err != nil {
			//make error
			return nil, err
		}
		dbname := os.Getenv("DBNAME")
		host := os.Getenv("HOST")
		port := os.Getenv("PORT")
		user := os.Getenv("USER")
		password := os.Getenv("PASSWORD")
		ssl := os.Getenv("SSLMODE")
		dbManager = &DBManager{}
		dbManager.once.Do(func() {
			cfg := &Config{
				url: fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v&", user, password, host, port, dbname, ssl),
			}
			for _, opt := range options {
				opt(&cfg.url)
			}
			dbManager.init(cfg)

		})
	}
	return dbManager, nil
}

// database setting: Timeout in seconds
func WithConnectionTimeout(timeout uint) func(*string) {
	return func(s *string) {
		*s += "connect_timeout=" + strconv.Itoa(int(timeout)) + "&"
	}
}

// database setting
func WithSSLCert(cert string) func(*string) {
	return func(s *string) {
		*s += "sslcert=" + cert + "&"
	}
}

// database setting
func WithSSLKey(key string) func(*string) {
	return func(s *string) {
		*s += "sslkey=" + key + "&"
	}
}

// database setting
func WithSSLRootCert(sslRootCert string) func(*string) {
	return func(s *string) {
		*s += "sslrootcert=" + sslRootCert + "&"
	}
}

func (d *DBManager) init(cfg *Config) {
	conn, err := sql.Open("postgres", cfg.url)
	if err != nil {
		panic(err)
	}
	if err = conn.Ping(); err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(10)
	d.db = conn
}

func (d *DBManager) GetDB() *sql.DB {
	return d.db
}

// database interface that is used in application
type Storage interface {
	InsertUser(context.Context, *proto.User) error
	GetUser(ctx context.Context, id string) (*proto.User, error)
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
	DoConversationExists(ctx context.Context, sender, receiver string) (bool, string, error)
	CreateConversation(ctx context.Context, cnv *proto.Conversation) error
	InsertMessage(ctx context.Context, msg *proto.Message) error
	GetUserConversations(ctx context.Context, id string) ([]*proto.ConversationSummary, error)
	LoadConversation(ctx context.Context, id string) ([]*proto.Message, error)
	SelectUsersToNewConversation(ctx context.Context, id string) ([]*proto.User, error)
	DoUserHaveDevice(ctx context.Context, userID string) (bool, []*proto.Device, error)
}
