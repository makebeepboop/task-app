package models

type Task struct {
	ID     int64
	Title  string
	Flow   string
	Number uint64
	Status string
}
