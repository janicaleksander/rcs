package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/janicaleksander/bcs/types/proto"
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
		config := struct {
			Database struct {
				Dbname   string `toml:"dbname"`
				Host     string `toml:"host"`
				Port     int    `toml:"port"`
				User     string `toml:"user"`
				Password string `toml:"password"`
				Sslmode  string `toml:"sslmode"`
			}
		}{}
		_, err := toml.DecodeFile("configproduction/database.toml", &config)
		if err != nil {
			return nil, err
		}
		dbname := config.Database.Dbname
		host := config.Database.Host
		port := config.Database.Port
		user := config.Database.User
		password := config.Database.Password
		ssl := config.Database.Sslmode
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
	GetUsersWithLVL(ctx context.Context, lower, upper int) ([]*proto.User, error)
	InsertUnit(ctx context.Context, unit *proto.Unit, userID string) error
	GetAllUnits(ctx context.Context) ([]*proto.Unit, error)
	GetUsersInUnit(ctx context.Context, unitID string) ([]*proto.User, error)
	IsUserInUnit(ctx context.Context, userID string) (bool, string, error)
	AssignUserToUnit(ctx context.Context, userID string, unitID string) error
	DeleteUserFromUnit(ctx context.Context, userID string, unitID string) error
	UpdateUserLastTimeOnline(ctx context.Context, id string, time time.Time) error
	DoConversationExists(ctx context.Context, senderID, receiverID string) (bool, string, error)
	CreateConversation(ctx context.Context, cnv *proto.Conversation) error
	InsertMessage(ctx context.Context, msg *proto.Message) error
	GetUserConversations(ctx context.Context, userID string) ([]*proto.ConversationSummary, error)
	LoadConversation(ctx context.Context, id string) ([]*proto.Message, error)
	SelectUsersToNewConversation(ctx context.Context, id string) ([]*proto.User, error)
	DoesUserHaveDevice(ctx context.Context, userID string) (bool, []*proto.Device, error)
	UpdateLocation(ctx context.Context, data *proto.UpdateLocationReq) error
	GetPins(ctx context.Context) ([]*proto.Pin, error)
	GetCurrentTask(ctx context.Context, deviceID string) (*proto.CurrentTask, error)
	GetDeviceTypes(ctx context.Context) ([]int32, error)
	InsertDevice(ctx context.Context, device *proto.Device) error
}
