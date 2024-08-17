package datahandler

/*
Inner Join: Returns rows when there is at least one match in both tables.

Left (Outer) Join: Returns all rows from the left table and matching rows
from the right table. If there is no match, it returns NULL for the missing values from the right table.

Right (Outer) Join: Returns all rows from the right table and matching rows
from the left table. If there is no match, it returns NULL for the missing values from the left table.

Full (Outer) Join: Returns all rows when combining data from both tables,
and NULL values are filled in for missing matches on either side.
*/
var FullPost = `
SELECT 
	p.id,
	p.title,
	p.content,
	p.user AS creatorId,
	p.creation_date,
	p.parent_post_id,
    p.image_id,
    i.image_data,
    p.postIsPublic,
    p.allowedUsers,
	(SELECT COUNT(*) FROM Postreaction WHERE post_id = p.id AND reaction = 1) AS upvotes,
    (SELECT COUNT(*) FROM Postreaction WHERE post_id = p.id AND reaction = 2) AS downvotes,
    u.username AS creator,
	r.reaction,
	JSON_OBJECT(
        'sports', CAST(t.sports AS BOOLEAN),
        'leisure', CAST(t.leisure AS BOOLEAN),
        'travelling', CAST(t.travelling AS BOOLEAN),
        'gaming', CAST(t.gaming AS BOOLEAN),
        'fun', CAST(t.fun AS BOOLEAN),
        'crafts', CAST(t.crafts AS BOOLEAN)
    ) AS tags
FROM 
	Post p
LEFT JOIN
	User u ON p.user = u.id
LEFT JOIN 
	Postreaction r ON r.post_id = p.id AND r.user_id = ? 
LEFT JOIN
    Images i ON i.id = p.image_id
LEFT JOIN
	Tags t ON p.id = t.post_id
`

var FullUser = `
SELECT
    u.username,
    u.email,
    u.firstname,
    u.lastname,
    u.gender,
    u.age,
    u.creation_date,
    u.isPublic, 
    u.image_id,
    i.image_data,
    u.aboutMe
FROM 
    User u
LEFT JOIN
    Images i ON i.id = u.image_id
WHERE
    username = ?
`

var isPublic = `
SELECT
    u.isPublic
FROM
    User u
WHERE username = ?
`

var GetChatMessage = `
SELECT 
    m.id,
    m.from_user_username AS senderUsername,
    m.to_user_username AS recieverUsername, 
    m.content,
    m.creation_date,
    m.seen
FROM 
    Messages m
WHERE
    m.from_user_username = ? or m.to_user_username = ?
`

var PostChatMessage = `
INSERT INTO Messages (
    from_user_username,
    to_user_username, 
    content,
    seen
)
VALUES (?, ?, ?, false);
`

var PostReaction = `
INSERT OR REPLACE INTO Postreaction (
    user_id,
    post_id, 
    reaction
)
VALUES (?, ?, ?)
`

// Group queries.
var CreateGroup = `
INSERT INTO Groups (
    title,
    description,
    creator_id
)
VALUES (?, ?, ?)
`

var GetGroup = `
SELECT
    id,
    title,
    description,
    creator_id
FROM
    Groups
WHERE
    id = ?
`

var GetGroups = `
SELECT
    id,
    title,
    description,
    creator_id
FROM
    Groups
`

// Group_members queries.
var GetGroupMembers = `
SELECT
    user_id,
    username,
    state
FROM
    Group_members
WHERE
    group_id = ?
`

var AddGroupMember = `
INSERT INTO Group_members (
    group_id,
    user_id,
    username,
    state
)
VALUES (?, ?, ?, ?)
`

var DeleteGroupMember = `
DELETE FROM Group_members
WHERE
    group_id = ?
    AND user_id = ?
`

var AcceptGroupMember = `
UPDATE Group_members
SET state = ?
WHERE group_id = ? AND user_id = ?
`

var GetUserGroups = `
SELECT
    group_id,
    state
FROM
    Group_members
WHERE
    user_id = ?
`

// Event queries.
var CreateEvent = `
INSERT INTO Events (
    title,
    content,
    date,
    creator_id,
    group_id
)
VALUES (?, ?, ?, ?, ?)
`

var GetEvents = `
SELECT
    id,
    title,
    content,
    date,
    creator_id
FROM
    Events
WHERE
    group_id = ?
`

var GetEvent = `
SELECT
    id,
    title,
    content,
    date,
    creator_id
FROM
    Events
WHERE
    id = ?
`

// Event_attendee queries.
var AddEventAttendee = `
INSERT INTO Event_attendees (
    event_id,
    user_id,
    username,
    attendance
)
VALUES (?, ?, ?, ?)
`

var UpdateEventAttendee = `
UPDATE Event_attendees
SET attendance = ?
WHERE event_id = ? AND user_id = ?
`

var DeleteEventAttendee = `
DELETE FROM Event_attendees
WHERE
    event_id = ?
    AND user_id = ?
`

var GetEventAttendees = `
SELECT
    attendance,
    username,
    user_id
FROM
    Event_attendees
WHERE
    event_id = ?
`

var PostFollowStatus = `
INSERT OR REPLACE INTO Followers (
    follower,
    beingFollowed, 
    currentstate
)
VALUES (?, ?, ?)
`

var PostNotification = `
INSERT INTO Notifications (user, actor, notification_type, date, group_id, group_name, seen)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(user, actor) DO UPDATE SET
    notification_type = excluded.notification_type,
    date = excluded.date,
    group_id = excluded.group_id,
    group_name = excluded.group_name,
    seen = excluded.seen;
`

// var PostNotification2 = `
// INSERT INTO Notifications (user, actor, notification_type, date, group_id, group_name, seen)
// VALUES (?, ?, ?, ?, ?, ?, ?)
// `

/*
user TEXT,

	actor TEXT,
*/
var GetFollowers = `
SELECT    
    f.follower,
    f.beingFollowed,
    f.currentstate
FROM 
    Followers f
WHERE
    f.beingFollowed = ? and f.currentstate = 'following'
`

var GetFollowing = `
SELECT    
    f.follower,
    f.beingFollowed,
    f.currentstate
FROM 
    Followers f
WHERE
    f.follower = ? and f.currentstate = 'following'
`
var GetNotifications = `
SELECT    
    n.id,
    n.user,
    n.actor,
    n.notification_type AS type,
    n.group_id,
    n.group_name,
    n.date
FROM 
    Notifications n
WHERE
    n.user = ? AND n.seen = false
`
