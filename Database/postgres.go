package Database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	_, err = tx.ExecContext(ctx, "INSERT INTO personal (user_id,name,surname) VALUES ($1,$2,$3)", user.Id, user.Personal.Name, user.Personal.Surname)
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
	rows, err := p.conn.Query(
		`SELECT users.id,users.email,users.rule_level,name,surname 
				FROM users
				INNER JOIN personal p on users.id = p.user_id
				WHERE (rule_level>=$1)`, lvl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*Proto.User, 0, 64)

	for rows.Next() {
		user := &Proto.User{Personal: &Proto.Personal{}}
		if err := rows.Scan(&user.Id, &user.Email, &user.RuleLvl, &user.Personal.Name, &user.Personal.Surname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *Postgres) InsertUnit(ctx context.Context, nameUnit string, isConfigured bool, userID string) error {
	// TODO user can be in the same time in one unit

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//check if user is unique in table
	var exists bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM user_to_unit WHERE user_id = $1)`, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user is already assigned to a unit")
	}

	//making row in unit table
	ID := uuid.New().String()
	_, err = tx.ExecContext(ctx, `INSERT INTO unit (id,name,is_configured) VALUES ($1,$2,$3)`, ID, nameUnit, isConfigured)
	if err != nil {
		return err
	}
	//making row in user to generated unit
	_, err = tx.ExecContext(ctx, `INSERT INTO user_to_unit (user_id,unit_id) VALUES ($1,$2)`, userID, ID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil

}

func (p *Postgres) GetAllUnits(ctx context.Context) ([]*Proto.Unit, error) {
	rows, err := p.conn.Query(`SELECT id,name,is_configured FROM unit`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	units := make([]*Proto.Unit, 0, 64)
	for rows.Next() {
		unit := &Proto.Unit{}
		err = rows.Scan(&unit.Id, &unit.Name, &unit.IsConfigured)
		if err != nil {
			//log error
			continue
		}
		units = append(units, unit)
	}
	return units, nil
}

// TODO maybe add more personal infos
func (p *Postgres) GetUsersInUnit(ctx context.Context, id string) ([]*Proto.User, error) {
	rows, err := p.conn.Query(
		`
	SELECT u.id, u.email, u.rule_level, personal.name, personal.surname  
	FROM personal
	INNER JOIN users u ON personal.user_id = u.id
	INNER JOIN user_to_unit utu ON u.id = utu.user_id 
	WHERE utu.unit_id = $1
`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*Proto.User, 0, 64)
	for rows.Next() {
		user := &Proto.User{Personal: &Proto.Personal{}}
		err = rows.Scan(&user.Id, &user.Email, &user.RuleLvl, &user.Personal.Name, &user.Personal.Surname)
		if err != nil {
			//log error
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *Postgres) IsUserInUnit(ctx context.Context, id string) (bool, string, error) {
	var exists bool
	var unitID string
	err := p.conn.QueryRow(`SELECT unit_id FROM user_to_unit WHERE user_id = $1 LIMIT 1`, id).Scan(&unitID)
	if errors.Is(err, sql.ErrNoRows) || err != nil {
		exists = false
	} else {
		exists = true
	}
	return exists, unitID, err
}
func (p *Postgres) AssignUserToUnit(ctx context.Context, userID string, unitID string) error {
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `INSERT INTO user_to_unit (user_id, unit_id) VALUES ($1,$2)`, userID, unitID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (p *Postgres) DeleteUserFromUnit(ctx context.Context, userID string, unitID string) error {
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `DELETE FROM user_to_unit WHERE (user_id=$1 AND unit_id=$2 )`, userID, unitID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (p *Postgres) GetUsersUnits(ctx context.Context, userID string) ([]string, error) {
	rows, err := p.conn.Query(`SELECT unit_id FROM user_to_unit WHERE (user_id=$1);`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	unitsID := make([]string, 0, 64)
	var id string
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			//TODO error
			return nil, err
		}
		unitsID = append(unitsID, id)
	}
	return unitsID, nil
}
