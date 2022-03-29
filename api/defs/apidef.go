package defs

//request
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

type NewVideo struct {
	AutorId int    `json:"author_id"`
	Name    string `json:"name"`
}

//response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type UserSession struct {
	Username  string `json:"user_name"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Id int `json:"id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
} 

type Comments struct {
	Comments []*Comment `json:"comments"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

// Data Model
type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

type Comment struct {
	Id         string `json:"id"`
	VideoId    string `json:"video_id"`
	AuthorName string `json:"author"`
	Content    string `json:"content"`
}

type SimpleSession struct {
	Username string
	TTL      int64
}

type User struct {
	Id        int
	LoginName string
	Pwd       string
}
