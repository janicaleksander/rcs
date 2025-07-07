package Database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/janicaleksander/bcs/User"
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
func (p *Postgres) InsertUser(ctx context.Context, user User.User) error {
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO users (id,email,password) VALUES ($1,$2,$3)", user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
