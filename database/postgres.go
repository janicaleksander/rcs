package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/types/user"
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
func (p *Postgres) GetUser(ctx context.Context, userID string) (*proto.User, error) {
	row := p.Conn.QueryRowContext(ctx,
		`SELECT
    				u.id,u.email,u.rule_level,u.last_time_online,
    				p.name,p.surname
				FROM users u 
				INNER JOIN personal p
				ON u.id = p.user_id
				WHERE (u.id = $1);
    `, userID)
	u := &proto.User{
		Id:            "",
		Email:         "",
		RuleLvl:       0,
		LasTimeOnline: nil,
		Personal: &proto.Personal{
			Name:    "",
			Surname: "",
		},
	}
	var timestamp sql.NullTime
	if err := row.Scan(&u.Id, &u.Email, &u.RuleLvl, &timestamp, &u.Personal.Name, &u.Personal.Surname); err != nil {
		return nil, err
	}
	if timestamp.Valid {
		u.LasTimeOnline = timestamppb.New(timestamp.Time)
	}
	return u, nil

}
func (p *Postgres) LoginUser(ctx context.Context, email, password string) (string, int, error) {
	row := p.Conn.QueryRowContext(ctx, `SELECT id, password,rule_level FROM users WHERE (email=$1)`, email)
	var id string
	var pwd string
	var role int
	if err := row.Scan(&id, &pwd, &role); err != nil {
		return "", -1, err
	}
	if !user.DecryptHash(password, pwd) {
		return "", -1, errors.New("invalid credentials")
	}
	return id, role, nil
}

// lower and upper are inclusive
func (p *Postgres) GetUsersWithLVL(ctx context.Context, lower, upper int) ([]*proto.User, error) {
	rows, err := p.Conn.QueryContext(ctx,
		`SELECT u.id,u.email,u.rule_level,u.last_time_online,p.name,p.surname 
				FROM users u 
				INNER JOIN personal p on u.id = p.user_id
				WHERE (rule_level>=$1 AND rule_level <= $2)`, lower, upper)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*proto.User, 0, 64)

	for rows.Next() {
		u := &proto.User{
			Id:            "",
			Email:         "",
			RuleLvl:       0,
			LasTimeOnline: nil,
			Personal: &proto.Personal{
				Name:    "",
				Surname: "",
			},
		}
		var timestamp sql.NullTime
		if err = rows.Scan(&u.Id, &u.Email, &u.RuleLvl, &timestamp, &u.Personal.Name, &u.Personal.Surname); err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		if timestamp.Valid {
			u.LasTimeOnline = timestamppb.New(timestamp.Time)
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, errors.New("no users")
	}
	return users, nil
}

func (p *Postgres) InsertUnit(ctx context.Context, unit *proto.Unit, userID string) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
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
	_, err = tx.ExecContext(ctx, `INSERT INTO unit (id,name) VALUES ($1,$2)`, unit.Id, unit.Name)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_to_unit (user_id,unit_id) VALUES ($1,$2)`, userID, unit.Id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetAllUnits(ctx context.Context) ([]*proto.Unit, error) {
	rows, err := p.Conn.QueryContext(ctx, `SELECT id,name FROM unit`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	units := make([]*proto.Unit, 0, 64)
	for rows.Next() {
		unit := &proto.Unit{
			Id:   "",
			Name: "",
		}
		err = rows.Scan(&unit.Id, &unit.Name)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		units = append(units, unit)
	}
	if len(units) == 0 {
		return nil, errors.New("no units")
	}
	return units, nil
}

func (p *Postgres) GetUsersInUnit(ctx context.Context, unitID string) ([]*proto.User, error) {
	rows, err := p.Conn.QueryContext(ctx,
		`
				SELECT u.id, u.email, u.rule_level, u.last_time_online,p.name, p.surname  
				FROM personal p
				INNER JOIN users u ON p.user_id = u.id
				INNER JOIN user_to_unit utu ON u.id = utu.user_id 
				WHERE (utu.unit_id = $1)
`, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*proto.User, 0, 64)
	for rows.Next() {
		u := &proto.User{
			Id:            "",
			Email:         "",
			RuleLvl:       0,
			LasTimeOnline: nil,
			Personal:      &proto.Personal{},
		}
		var timestamp sql.NullTime
		err = rows.Scan(&u.Id, &u.Email, &u.RuleLvl, &timestamp, &u.Personal.Name, &u.Personal.Surname)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		if timestamp.Valid {
			u.LasTimeOnline = timestamppb.New(timestamp.Time)
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, errors.New("no users")
	}
	return users, nil
}

func (p *Postgres) IsUserInUnit(ctx context.Context, userID string) (bool, string, error) {
	var unitID string
	err := p.Conn.QueryRowContext(ctx, `SELECT unit_id FROM user_to_unit WHERE (user_id = $1)`, userID).Scan(&unitID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "", sql.ErrNoRows
	}
	if err != nil {
		return false, "", err
	}

	return true, unitID, err
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

func (p *Postgres) UpdateUserLastTimeOnline(ctx context.Context, id string, t time.Time) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `UPDATE users SET last_time_online = $1 WHERE (id = $2);`, t, id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DoConversationExists(ctx context.Context, sender, receiver string) (bool, string, error) {
	var conversationID string
	err := p.Conn.QueryRowContext(ctx, `
    SELECT uc1.conversation_id
    FROM user_conversation uc1
    JOIN user_conversation uc2 
        ON uc1.conversation_id = uc2.conversation_id
    WHERE (uc1.user_id = $1 
      AND uc2.user_id = $2)
    LIMIT 1
`, sender, receiver).Scan(&conversationID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "", sql.ErrNoRows
	}

	if err != nil {
		return false, "", err
	}

	return true, conversationID, nil
}

func (p *Postgres) CreateConversation(ctx context.Context, cnv *proto.Conversation) error {
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
		cnv.SenderID, cnv.Id, nil,
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
	_, err = tx.ExecContext(ctx,
		`UPDATE user_conversation 
				SET last_seen_message_id=$1
				WHERE (conversation_id=$2) `, msg.Id, msg.ConversationID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

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
WHERE (uc.user_id = $1 AND other_uc.user_id IS NOT NULL)
ORDER BY m.sent_at DESC NULLS LAST ;`, id)

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

		var messageID sql.NullString
		var sender sql.NullString
		var content sql.NullString
		var sentAt sql.NullTime
		var name string
		var surname string
		err = rows.Scan(&cs.ConversationId, &messageID, &sender, &content, &sentAt, &cs.WithID, &name, &surname)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		if messageID.Valid {
			cs.LastMessage.Id = messageID.String
		}
		if sender.Valid {
			cs.LastMessage.SenderID = sender.String
		}
		if content.Valid {
			cs.LastMessage.Content = content.String
		}
		if sentAt.Valid {
			cs.LastMessage.SentAt = timestamppb.New(sentAt.Time)
		}
		cs.Nametag = name + " " + surname
		conversationsSummary = append(conversationsSummary, cs)
	}

	if len(conversationsSummary) == 0 {
		return nil, errors.New("no cnvs summary")
	}
	return conversationsSummary, nil

}

func (p *Postgres) LoadConversation(ctx context.Context, cnvID string) ([]*proto.Message, error) {
	rows, err := p.Conn.QueryContext(ctx,
		`SELECT id,user_id,conversation_id,content,sent_at 
				FROM message 
				WHERE conversation_id=$1 ORDER BY sent_at`, cnvID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
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
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		m.SentAt = timestamppb.New(timestamp)
		messages = append(messages, m)
	}
	if len(messages) == 0 {
		return nil, errors.New("no messages")
	}
	return messages, nil
}
func (p *Postgres) SelectUsersToNewConversation(ctx context.Context, userID string) ([]*proto.User, error) {
	rows, err := p.Conn.QueryContext(ctx,
		`SELECT u.id, u.email, p.name, p.surname 
				FROM users u
				INNER JOIN personal p ON u.id = p.user_id
				WHERE u.id <> $1
				AND NOT EXISTS (
					SELECT 1 
					FROM user_conversation uc1
					INNER JOIN user_conversation uc2 ON uc1.conversation_id = uc2.conversation_id
					WHERE uc1.user_id = u.id AND uc2.user_id = $1);
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*proto.User, 0, 64)
	for rows.Next() {

		usr := &proto.User{
			Id:    "",
			Email: "",
			Personal: &proto.Personal{
				Name:    "",
				Surname: "",
			},
		}
		err = rows.Scan(&usr.Id, &usr.Email, &usr.Personal.Name, &usr.Personal.Surname)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		users = append(users, usr)
	}
	if len(users) == 0 {
		return nil, errors.New("no users")
	}
	return users, nil
}

// THIS IS ONLY WHEN users is in one unit
func (p *Postgres) DoesUserHaveDevice(ctx context.Context, userID string) (bool, []*proto.Device, error) {
	rows, err := p.Conn.QueryContext(ctx, `SELECT d.id,d.name,d.last_time_online,d.owner,d.type FROM device d 
						    INNER JOIN users u 
						    ON d.owner = u.id
						    WHERE (u.id = $1);`, userID)
	if err != nil {
		return false, nil, err
	}
	devices := make([]*proto.Device, 0, 8)
	for rows.Next() {
		d := &proto.Device{
			Id:             "",
			Name:           "",
			Owner:          "",
			LastTimeOnline: nil,
			Type:           0,
		}
		var lastTimeOnline sql.NullTime
		err = rows.Scan(&d.Id, &d.Name, &lastTimeOnline, &d.Owner, &d.Type)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		if lastTimeOnline.Valid {
			d.LastTimeOnline = timestamppb.New(lastTimeOnline.Time)
		}
		devices = append(devices, d)
	}

	if err != nil {
		return false, nil, err
	}
	if len(devices) == 0 {
		return false, nil, errors.New("no devices")
	}
	return true, devices, nil
}

func (p *Postgres) UpdateLocation(ctx context.Context, data *proto.UpdateLocationReq) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `     
        INSERT INTO device_location (device_id, location, changed_at)
        VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326),$4)`,
		data.DeviceID,
		data.Location.Longitude,
		data.Location.Latitude,
		time.Now(),
	)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetPins(ctx context.Context) ([]*proto.Pin, error) {
	rows, err := p.Conn.QueryContext(ctx, `SELECT DISTINCT ON (device_id) p.name,p.surname, device_id, 
    								ST_Y(location::geometry) as lat,
    								ST_X(location::geometry) as lng, 
                               		changed_at
									FROM device_location
									INNER JOIN device d ON d.id = device_id
									INNER JOIN users u ON d.owner = u.id
									INNER JOIN personal p on u.id = p.user_id
									ORDER BY device_id, changed_at DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pins := make([]*proto.Pin, 0, 32)
	for rows.Next() {
		pin := &proto.Pin{
			DeviceID:     "",
			OwnerName:    "",
			OwnerSurname: "",
			Location:     &proto.Location{},
			LastOnline:   nil,
		}
		var tmp time.Time
		err = rows.Scan(&pin.OwnerName, &pin.OwnerSurname, &pin.DeviceID, &pin.Location.Latitude, &pin.Location.Longitude, &tmp)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		pin.LastOnline = timestamppb.New(tmp)
		pins = append(pins, pin)
	}
	if len(pins) == 0 {
		return nil, errors.New("no pins")
	}
	return pins, nil
}

func (p *Postgres) GetCurrentTask(ctx context.Context, deviceID string) (*proto.CurrentTask, error) {
	row := p.Conn.QueryRowContext(ctx,
		`SELECT
    					d.owner,t.id,t.name,t.description,t.state,t.completion_date,t.deadline
				FROM device d
				INNER JOIN current_user_task cut ON d.owner = cut.user_id
				INNER JOIN task t ON t.id = cut.task_id
				WHERE d.id = $1`, deviceID)
	t := &proto.CurrentTask{
		Task: &proto.Task{
			Id:             "",
			Name:           "",
			Description:    "",
			State:          0,
			CompletionDate: nil,
			Deadline:       nil,
		},
		UserID: "",
	}
	var taskCompletionDate sql.NullTime
	var deadline sql.NullTime
	err := row.Scan(&t.UserID, &t.Task.Id, &t.Task.Name, &t.Task.Description, &t.Task.State, &taskCompletionDate, &deadline)
	if err != nil {
		return nil, err
	}
	if taskCompletionDate.Valid {
		t.Task.CompletionDate = timestamppb.New(taskCompletionDate.Time)
	}
	if deadline.Valid {
		t.Task.Deadline = timestamppb.New(deadline.Time)
	}
	return t, nil
}

func (p *Postgres) GetTask(ctx context.Context, taskID string) (*proto.Task, error) {
	row := p.Conn.QueryRowContext(ctx, `SELECT t.id,t.name,t.description,t.state,t.completion_date,t.deadline
				FROM task t 
				WHERE t.id = $1 ORDER BY t.deadline`, taskID)

	t := &proto.Task{
		Id:             "",
		Name:           "",
		Description:    "",
		State:          0,
		CompletionDate: nil,
		Deadline:       nil,
	}
	var completion sql.NullTime
	var deadline sql.NullTime
	err := row.Scan(&t.Id, &t.Name, &t.Description, &t.State, &completion, &deadline)
	if err != nil {
		return nil, err
	}
	if completion.Valid {
		t.CompletionDate = timestamppb.New(completion.Time)
	}
	if deadline.Valid {
		t.Deadline = timestamppb.New(deadline.Time)
	}
	return t, err
}

func (p *Postgres) GetUserTasks(ctx context.Context, deviceID string) ([]*proto.Task, error) {
	rows, err := p.Conn.QueryContext(ctx,
		`SELECT t.id,t.name,t.description,t.state,t.completion_date,t.deadline
				FROM task t
				INNER JOIN user_to_task utt ON utt.task_id  = t.id
				INNER JOIN device d ON d.owner = utt.user_id
				WHERE d.id = $1 ORDER BY t.deadline;
				`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]*proto.Task, 0, 16)
	var completion sql.NullTime
	var deadline sql.NullTime
	for rows.Next() {
		t := &proto.Task{
			Id:             "",
			Name:           "",
			Description:    "",
			State:          0,
			CompletionDate: nil,
			Deadline:       nil,
		}
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.State, &completion, &deadline)
		if completion.Valid {
			t.CompletionDate = timestamppb.New(completion.Time)
		}
		if deadline.Valid {
			t.Deadline = timestamppb.New(deadline.Time)
		}
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		return nil, errors.New("no tasks")
	}
	return tasks, nil
}

func (p *Postgres) UpdateCurrentTask(ctx context.Context, newTaskID, userID string) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`INSERT INTO current_user_task (user_id,task_id) 
				VALUES ($1,$2) ON CONFLICT (user_id)
				DO UPDATE SET user_id=excluded.user_id,task_id=excluded.task_id`, userID, newTaskID)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeleteTask(ctx context.Context, taskID string) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `DELETE FROM task WHERE task.id = $1`, taskID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil

}

func (p *Postgres) GetDeviceTypes(ctx context.Context) ([]int32, error) {
	rows, err := p.Conn.QueryContext(ctx, `SELECT type FROM device_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	types := make([]int32, 0, 2)
	for rows.Next() {
		var t int32
		err = rows.Scan(&t)
		if err != nil {
			utils.Logger.Error(err.Error())
			continue
		}
		types = append(types, t)
	}
	if len(types) == 0 {
		return nil, errors.New("no types")
	}
	return types, nil
}

func (p *Postgres) InsertDevice(ctx context.Context, device *proto.Device) error {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`INSERT INTO device (id,name,last_time_online,owner,type) 
				   VALUES ($1,$2,$3,$4,$5)`,
		device.Id,
		device.Name,
		nil,
		device.Owner,
		device.Type,
	)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
