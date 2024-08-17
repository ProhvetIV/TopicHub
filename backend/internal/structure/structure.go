package structure

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     []byte    `json:"password"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Gender       string    `json:"gender"`
	Age          int       `json:"age"`
	CreationDate string    `json:"creation_date"`
}

type Post struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Creator      string   `json:"creator"`
	CreationDate string   `json:"creation_date"`
	Upvotes      Likes    `json:"upvotes"`
	Parent       string   `json:"parent"`
	Reply        []Post   `json:"reply"`
	ImagePath    string   `json:"image_path"`
	Categories   []string `json:"categories"`
}

/* type PostAllowedUsers struct {
	ID 	int 	`json:id`
	Username 	string	`json:username`
} */

type Users struct {
	Users []User `json:"users"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Likes struct {
	Total    int  `json:"total"`
	Liked    bool `json:"liked"`
	Disliked bool `json:"disliked"`
}

// not needed?
type RegisterUser struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
}

type LogInInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
