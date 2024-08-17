CREATE TABLE IF NOT EXISTS Event_attendees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    attendance INTEGER NOT NULL,
    FOREIGN KEY (event_id) REFERENCES Events(id),
    FOREIGN KEY (user_id) REFERENCES User(id)
);