CREATE TABLE IF NOT EXISTS Followers (
    follower TEXT,
    beingFollowed TEXT,
    currentstate TEXT,
    FOREIGN KEY (follower) REFERENCES User(username),
    FOREIGN KEY (beingFollowed) REFERENCES User(username),
    PRIMARY KEY (follower, beingFollowed)
);
