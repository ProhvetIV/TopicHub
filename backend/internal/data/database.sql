CREATE TABLE IF NOT EXISTS User(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    gender TEXT NOT NULL,
    age INTEGER NOT NULL,
    isPublic BOOLEAN,
    image_id INTEGER,
    aboutMe TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (image_id) REFERENCES Images(id)
);

CREATE TABLE IF NOT EXISTS Post(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user INTEGER,
    title TEXT,
    content TEXT NOT NULL,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    parent_post_id INTEGER,
    image_id INTEGER,
    group_id INTEGER,
    postIsPublic TEXT,
    allowedUsers TEXT,
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    FOREIGN KEY (user) REFERENCES User(id),
    FOREIGN KEY (image_id) REFERENCES Images(id)
);

CREATE TABLE IF NOT EXISTS Tags(
    sports BOOLEAN,
    leisure BOOLEAN,
    travelling BOOLEAN,
    gaming BOOLEAN,
    fun BOOLEAN,
    crafts BOOLEAN,
    post_id INTEGER,
    FOREIGN KEY (post_id) REFERENCES Post(id)
);

CREATE TABLE IF NOT EXISTS Session(
    user_id INTEGER PRIMARY KEY UNIQUE,
    session TEXT,
    expires INTEGER,
    FOREIGN KEY (user_id) REFERENCES User(id)
);

CREATE TABLE IF NOT EXISTS Postreaction(
    user_id INTEGER,
    post_id INTEGER, 
    reaction INTEGER,
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (post_id) REFERENCES Post(id),
    PRIMARY KEY (user_id, post_id)
);

CREATE TABLE IF NOT EXISTS Messages(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_user_username TEXT,
    to_user_username TEXT, 
    content TEXT,
    creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    seen BOOLEAN,
    FOREIGN KEY (from_user_username) REFERENCES User(username),
    FOREIGN KEY (to_user_username) REFERENCES User(username)
);


CREATE TABLE IF NOT EXISTS Images(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image_name TEXT,
    image_data BLOB
);