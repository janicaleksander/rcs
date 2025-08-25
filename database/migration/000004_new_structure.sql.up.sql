CREATE TABLE IF NOT EXISTS users
(
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    rule_level INT NOT NULL,
    last_time_online TIMESTAMP
);
CREATE TABLE IF NOT EXISTS personal
(
    user_id UUID PRIMARY KEY REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS unit
(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    is_configured BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS device_type
(
    type VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS device
(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    last_time_online TIMESTAMP,
    owner UUID REFERENCES users(id) ON DELETE SET NULL,
    type VARCHAR(255) REFERENCES device_type(type) NOT NULL


);

CREATE TABLE IF NOT EXISTS user_to_unit
(
    user_id UUID REFERENCES users(id) UNIQUE NOT NULL,
    unit_id UUID REFERENCES unit(id) NOT NULL
);


CREATE TABLE IF NOT EXISTS device_to_unit
(
    unit_id UUID REFERENCES unit(id) NOT NULL,
    device_id UUID REFERENCES device(id) NOT NULL,
    PRIMARY KEY(unit_id, device_id)

);

CREATE TABLE IF NOT EXISTS conversation
(
    id UUID PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS message
(
    id UUID PRIMARY KEY ,
    user_id UUID REFERENCES users(id) NOT NULL,
    conversation_id UUID REFERENCES conversation(id)  NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL

);
CREATE TABLE IF NOT EXISTS user_conversation
(
    user_id UUID REFERENCES users(id) NOT NULL,
    conversation_id UUID REFERENCES conversation(id)  NOT NULL,
    last_seen_message_id UUID REFERENCES message(id),
    UNIQUE(user_id, conversation_id)

);


CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_unit_name ON unit(name);
CREATE INDEX IF NOT EXISTS idx_device_name ON device(name);
CREATE INDEX IF NOT EXISTS idx_device_owner ON device(owner);
CREATE INDEX IF NOT EXISTS idx_user_to_unit_unit ON user_to_unit(unit_id);
CREATE INDEX IF NOT EXISTS idx_device_to_unit_device ON device_to_unit(device_id);
CREATE INDEX IF NOT EXISTS idx_message_conversation_sent
    ON message(conversation_id, sent_at);
CREATE INDEX IF NOT EXISTS idx_user_conversation_user ON user_conversation(user_id);
CREATE INDEX IF NOT EXISTS idx_user_conversation_conv ON user_conversation(conversation_id);
CREATE INDEX IF NOT EXISTS idx_user_conversation_last_seen ON user_conversation(last_seen_message_id);

