package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/User"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Struct that is made to implement Storage interface
// It holds connection to make queries to PostgreSQL.
type Postgres struct {
	Conn *sql.DB
}

// Implemented function from Storage interface.
// It takes context (e.g. to set max time for query), User and if no errors insert into database.
// If error occurs -> return it, and rollback transaction.
func (p *Postgres) InsertUser(ctx context.Context, user *proto.User) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
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
	rows, err := p.Conn.Query(`SELECT id, password,rule_level FROM users WHERE (email=$1)`, email)
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
	if !user.DecryptHash(password, pwd) {
		return "", -1, errors.New("invalid credentials")
	}

	return id, role, nil
}

func (p *Postgres) GetUsersWithLVL(ctx context.Context, lvl int) ([]*proto.User, error) {
	rows, err := p.Conn.Query(
		`SELECT users.id,users.email,users.rule_level,name,surname 
				FROM users
				INNER JOIN personal p on users.id = p.user_id
				WHERE (rule_level>=$1)`, lvl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*proto.User, 0, 64)

	for rows.Next() {
		user := &proto.User{Personal: &proto.Personal{}}
		if err := rows.Scan(&user.Id, &user.Email, &user.RuleLvl, &user.Personal.Name, &user.Personal.Surname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *Postgres) InsertUnit(ctx context.Context, nameUnit string, isConfigured bool, userID string) error {
	// TODO user can be in the same time in one unit

	tx, err := p.Conn.BeginTx(ctx, nil)
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

func (p *Postgres) GetAllUnits(ctx context.Context) ([]*proto.Unit, error) {
	rows, err := p.Conn.Query(`SELECT id,name,is_configured FROM unit`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	units := make([]*proto.Unit, 0, 64)
	for rows.Next() {
		unit := &proto.Unit{}
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
func (p *Postgres) GetUsersInUnit(ctx context.Context, id string) ([]*proto.User, error) {
	rows, err := p.Conn.Query(
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
	users := make([]*proto.User, 0, 64)
	for rows.Next() {
		user := &proto.User{Personal: &proto.Personal{}}
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
	err := p.Conn.QueryRow(`SELECT unit_id FROM user_to_unit WHERE user_id = $1 LIMIT 1`, id).Scan(&unitID)
	if errors.Is(err, sql.ErrNoRows) || err != nil {
		exists = false
	} else {
		exists = true
	}
	return exists, unitID, err
}
func (p *Postgres) AssignUserToUnit(ctx context.Context, userID string, unitID string) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
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
	tx, err := p.Conn.BeginTx(ctx, nil)
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

// MESSAGE SERVICE SQL
// Return:
// - true if exists false otherwise
// - string id of conversation between two user
// - error
func (p *Postgres) IsConversationExists(ctx context.Context, sender, receiver string) (bool, string, error) {
	var conversationID string
	err := p.Conn.QueryRow(`
    SELECT uc1.conversation_id
    FROM user_conversation uc1
    JOIN user_conversation uc2 
        ON uc1.conversation_id = uc2.conversation_id
    WHERE uc1.user_id = $1 
      AND uc2.user_id = $2
    LIMIT 1
`, sender, receiver).Scan(&conversationID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "", err
	}

	return true, conversationID, nil
}

func (p *Postgres) CreateAndAssignConversation(ctx context.Context, cnv *proto.CreateConversationAndAssign) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO conversation (id) VALUES ($1)`, cnv.Id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_conversation (user_id, conversation_id, last_seen_message_id) 
         VALUES ($1, $2, $3), ($4, $2, $3)`,
		cnv.SenderID, cnv.Id, -1,
		cnv.ReceiverID,
	)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) InsertMessage(ctx context.Context, msg *proto.Message) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`INSERT INTO message (id, user_id, conversation_id, content, sent_at) 
				VALUES ($1,$2,$3,$4,$5)`, msg.Id, msg.SenderID, msg.ConversationID, msg.Content, msg.SentAt.AsTime())
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

//TODO make a one place where i set a uuid of sth e.g message or conversation and other
//every ID i setting in backend, to do this we have to give to fb funcs a proto structs (full with everything during creating)

func (p *Postgres) GetUserConversations(ctx context.Context, id string) ([]*proto.ConversationSummary, error) {
	rows, err := p.Conn.QueryContext(ctx,
		`SELECT
					uc.conversation_id,
					m.id as message_id,
					m.user_id,
					m.content,
					m.sent_at,
					other_uc.user_id as other_user_id,
					pc.name,
					pc.surname
				FROM user_conversation uc
				LEFT JOIN user_conversation other_uc ON other_uc.conversation_id = uc.conversation_id AND other_uc.user_id != uc.user_id
				LEFT JOIN personal pc ON pc.user_id = other_uc.user_id
				    LEFT JOIN message m ON m.id = (
					SELECT id FROM message
					WHERE conversation_id = uc.conversation_id
					ORDER BY sent_at DESC
					LIMIT 1
					)
WHERE (uc.user_id = $1 AND other_uc IS NOT NULL)
ORDER BY m.sent_at DESC;`, id)

	if err != nil {
		return nil, err
	}
	conversationsSummary := make([]*proto.ConversationSummary, 0, 64)
	for rows.Next() {
		cs := &proto.ConversationSummary{
			ConversationId: "",
			WithID:         "",
			Nametag:        "",
			LastMessage: &proto.Message{
				Id:             "",
				SenderID:       "",
				ConversationID: "",
				Content:        "",
				SentAt:         nil,
			},
		}
		var timestamp time.Time
		var name string
		var surname string
		err = rows.Scan(&cs.ConversationId, &cs.LastMessage.Id, &cs.LastMessage.SenderID, &cs.LastMessage.Content, &timestamp, &cs.WithID, &name, &surname)
		if err != nil {
			//TODO
		}
		cs.LastMessage.SentAt = timestamppb.New(timestamp)
		cs.Nametag = name + " " + surname
		conversationsSummary = append(conversationsSummary, cs)
	}
	return conversationsSummary, nil

}

func (p *Postgres) LoadConversation(ctx context.Context, id string) ([]*proto.Message, error) {
	rows, err := p.Conn.Query(
		`SELECT id,user_id,conversation_id,content,sent_at 
				FROM message 
				WHERE conversation_id=$1 ORDER BY sent_at`, id)
	if err != nil {
		return nil, err
	}
	messages := make([]*proto.Message, 0, 64)
	for rows.Next() {
		m := &proto.Message{
			Id:             "",
			SenderID:       "",
			ConversationID: "",
			Content:        "",
			SentAt:         nil,
		}
		var timestamp time.Time
		err = rows.Scan(&m.Id, &m.SenderID, &m.ConversationID, &m.Content, &timestamp)
		m.SentAt = timestamppb.New(timestamp)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		messages = append(messages, m)
	}
	return messages, nil
}
