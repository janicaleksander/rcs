DROP TABLE IF EXISTS  device_to_unit;
DROP TABLE IF EXISTS users CASCADE ;
CREATE TABLE IF NOT EXISTS users
(
    id UUID,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    rule_level INT NOT NULL,
    last_time_online TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (email)
);
DROP TABLE IF EXISTS personal CASCADE ;
CREATE TABLE IF NOT EXISTS personal
(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    PRIMARY KEY (user_id)
);
DROP TABLE IF EXISTS unit CASCADE ;
CREATE TABLE IF NOT EXISTS unit
(
    id UUID,
    name TEXT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);
DROP TABLE IF EXISTS user_to_unit CASCADE ;
CREATE TABLE IF NOT EXISTS user_to_unit
(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    unit_id UUID REFERENCES unit(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (user_id)
);
DROP TABLE IF EXISTS device_type CASCADE ;
CREATE TABLE IF NOT EXISTS device_type
(
    type INT,
    PRIMARY KEY (type)
);
DROP TABLE IF EXISTS device CASCADE ;
CREATE TABLE IF NOT EXISTS device
(
    id UUID,
    name TEXT NOT NULL,
    last_time_online TIMESTAMP NULL,
    owner UUID REFERENCES users(id) ON DELETE RESTRICT NOT NULL ,
    type INT REFERENCES device_type(type) ON DELETE RESTRICT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

DROP TABLE IF EXISTS device_location CASCADE ;
CREATE TABLE IF NOT EXISTS device_location
(
    device_id UUID REFERENCES device(id) ON DELETE CASCADE NOT NULL,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    changed_at TIMESTAMP NOT NULL
);



DROP TABLE IF EXISTS conversation CASCADE ;
CREATE TABLE IF NOT EXISTS conversation
(
    id UUID,
    PRIMARY KEY(id)
);

DROP TABLE IF EXISTS message CASCADE ;
CREATE TABLE IF NOT EXISTS message
(
    id UUID,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    conversation_id UUID REFERENCES conversation(id) ON DELETE CASCADE NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)

);

DROP TABLE IF EXISTS user_conversation CASCADE ;
CREATE TABLE IF NOT EXISTS user_conversation
(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    conversation_id UUID REFERENCES conversation(id) ON DELETE CASCADE  NOT NULL,
    last_seen_message_id UUID REFERENCES message(id) ON DELETE SET NULL NULL,
    PRIMARY KEY (user_id, conversation_id)

);

DROP TABLE IF EXISTS task CASCADE ;
CREATE TABLE IF NOT EXISTS task (
    id UUID ,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    state INT NOT NULL,
    deadline TIMESTAMP NULL,
    completion_date TIMESTAMP NULL,
    PRIMARY KEY(id)
);

DROP TABLE IF EXISTS user_to_task CASCADE ;
CREATE TABLE IF NOT EXISTS user_to_task (
    user_id UUID REFERENCES users(id) ON DELETE RESTRICT  NOT NULL,
    task_id UUID REFERENCES task(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id)
);


DROP TABLE IF EXISTS current_user_task CASCADE ;
CREATE TABLE IF NOT EXISTS current_user_task (
    task_id UUID REFERENCES task(id) ON DELETE CASCADE  NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE RESTRICT   NOT NULL,
    PRIMARY KEY (user_id)

);
CREATE INDEX IF NOT EXISTS idx_users_rule_level ON users(rule_level);
CREATE INDEX IF NOT EXISTS idx_personal_name_surname ON personal(name, surname);
CREATE INDEX IF NOT EXISTS idx_device_type ON device(type);
CREATE INDEX IF NOT EXISTS idx_message_user_id ON message(user_id);
CREATE INDEX IF NOT EXISTS idx_message_sent_at ON message(sent_at);
CREATE INDEX IF NOT EXISTS idx_device_location_device_id ON device_location(device_id);
CREATE INDEX IF NOT EXISTS idx_device_location_changed_at ON device_location(changed_at);
CREATE INDEX IF NOT EXISTS idx_device_location_device_time ON device_location(device_id, changed_at);
CREATE INDEX IF NOT EXISTS idx_device_location_geography ON device_location USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_task_state ON task(state);
CREATE INDEX IF NOT EXISTS idx_task_completion_date ON task(completion_date);
CREATE INDEX IF NOT EXISTS idx_user_to_task_user_id ON user_to_task(user_id);
CREATE INDEX IF NOT EXISTS idx_task_state_user ON task(state) INCLUDE (completion_date);

