CREATE TABLE IF NOT EXISTS task (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    state INT NOT NULL,
    completion_date TIMESTAMP NULL
);
/*
0 - not done
1 - done
*/
CREATE TABLE IF NOT EXISTS user_to_task (
    user_id UUID REFERENCES  users(id) NOT NULL,
    task_id UUID REFERENCES task(id),
    UNIQUE(task_id)
);

CREATE TABLE IF NOT EXISTS current_user_task (
    task_id UUID REFERENCES task(id),
    user_id UUID REFERENCES users(id),
    PRIMARY KEY (user_id,task_id)
)
