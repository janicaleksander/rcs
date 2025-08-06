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
CREATE TABLE IF NOT EXISTS device
(
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);


CREATE TABLE IF NOT EXISTS user_to_unit
(
    user_id UUID REFERENCES users(id) UNIQUE NOT NULL,
    unit_id UUID REFERENCES unit(id) NOT NULL
);


CREATE TABLE IF NOT EXISTS device_to_unit
(
    unit_id UUID REFERENCES unit(id) NOT NULL,
    device_id INT REFERENCES device(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS conversation
(
    id UUID PRIMARY KEY ,
    name VARCHAR(255) NOT NULL
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


/*
id uuid
user_id uuid (ref) who sent
conversation_id uuid (ref) to which conv sent
content varchar
sent_at timestamp
seen_at timestamp
*/