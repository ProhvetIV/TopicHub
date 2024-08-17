CREATE TABLE IF NOT EXISTS Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT,
    actor TEXT,
    notification_type INTEGER,
    group_id INTEGER,
    group_name TEXT,
    seen BOOLEAN,
    date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user) REFERENCES User(username),
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    UNIQUE(user, actor)
);
