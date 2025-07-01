package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

// TODO
// 1. Migrations
// 2. Come up with basic base structure
// 3.Create basic method in interface
// 4. Implement interface
type Postgres struct {
	Conn *sql.DB
}

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
	return &Postgres{Conn: db}, nil
}

// Timeout in seconds
func WithConnectionTimeout(timeout uint) func(*string) {
	return func(s *string) {
		*s += "connect_timeout=" + strconv.Itoa(int(timeout)) + "&"
	}
}

func WithSSLCert(cert string) func(*string) {
	return func(s *string) {
		*s += "sslcert=" + cert + "&"
	}
}

func WithSSLKey(key string) func(*string) {
	return func(s *string) {
		*s += "sslkey=" + key + "&"
	}
}

func WithSSLRootCert(sslRootCert string) func(*string) {
	return func(s *string) {
		*s += "sslrootcert=" + sslRootCert + "&"
	}
}
