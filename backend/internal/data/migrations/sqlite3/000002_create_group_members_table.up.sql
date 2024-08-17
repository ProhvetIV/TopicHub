CREATE TABLE IF NOT EXISTS Group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    state INTEGER,
    FOREIGN KEY (username) REFERENCES User(username)
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    FOREIGN KEY (user_id) REFERENCES User(id)
);