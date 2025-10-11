package model

type Message struct {
	ID      int64
	UserID  int64
	Text    string
	Title   string
	Name    string
}

type File struct {
	ID      int64
	UserID  int64
	Name    string
}