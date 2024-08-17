CREATE TABLE IF NOT EXISTS PostAllowedUsers(
    post_id INTEGER,
    user_id INTEGER,
    username TEXT,
    FOREIGN KEY (post_id) REFERENCES Post(id),
    FOREIGN KEY (user_id) REFERENCES User(id)
);