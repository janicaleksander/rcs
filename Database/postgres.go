package Database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/User"
	_ "github.com/lib/pq"
	"strconv"
)

// Struct that is made to implement Storage interface
// It holds connection to make queries to PostgreSQL.
type Postgres struct {
	conn *sql.DB
}

// Constructor of  new PostgreSQL. It can be modified by optional options.
// It has to be filled up with Database credentials to perform connection.
// If no errors occur it returns pointer to struct
func NewPostgres(
	dbname string,
	user string,
	password string,
	host string,
	port string,
	sslmode string,
	options ...func(*string)) (*Postgres, error) {

	// default connection string, to which we can add some extras
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v&", user, password, host, port, dbname, sslmode)
	for _, o := range options {
		o(&connStr)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{conn: db}, nil
}

// Database setting: Timeout in seconds
func WithConnectionTimeout(timeout uint) func(*string) {
	return func(s *string) {
		*s += "connect_timeout=" + strconv.Itoa(int(timeout)) + "&"
	}
}

// Database setting
func WithSSLCert(cert string) func(*string) {
	return func(s *string) {
		*s += "sslcert=" + cert + "&"
	}
}

// Database setting
func WithSSLKey(key string) func(*string) {
	return func(s *string) {
		*s += "sslkey=" + key + "&"
	}
}

// Database setting
func WithSSLRootCert(sslRootCert string) func(*string) {
	return func(s *string) {
		*s += "sslrootcert=" + sslRootCert + "&"
	}
}

// Implemented function from Storage interface.
// It takes context (e.g. to set max time for query), User and if no errors insert into Database.
// If error occurs -> return it, and rollback transaction.
func (p *Postgres) InsertUser(ctx context.Context, user *Proto.User) error {
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO users (id,email,password,rule_level) VALUES ($1,$2,$3,$4)", user.Id, user.Email, user.Password, user.RuleLvl)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) LoginUser(ctx context.Context, email, password string) (string, int, error) {
	rows, err := p.conn.Query(`SELECT id, password,rule_level FROM users WHERE (email=$1)`, email)
	if err != nil {
		return "", -1, err
	}
	defer rows.Close()
	var id string
	var pwd string
	var role int
	for rows.Next() {
		if err = rows.Scan(&id, &pwd, &role); err != nil {
			return "", -1, err
		}
	}
	if !User.DecryptHash(password, pwd) {
		return "", -1, errors.New("invalid credentials")
	}

	return id, role, nil
}

func (p *Postgres) GetUsersWithLVL(ctx context.Context, lvl int) ([]*Proto.User, error) {
	rows, err := p.conn.Query(`SELECT id,email,password,rule_level FROM users WHERE (rule_level>=$1)`, lvl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*Proto.User, 0, 64)

	for rows.Next() {
		user := &Proto.User{}
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.RuleLvl); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
