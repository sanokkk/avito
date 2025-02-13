package models

import "github.com/google/uuid"

type Item struct {
	Id    uuid.UUID
	Title string
	Cost  int
}

type UserItem struct {
	Id       uuid.UUID
	ItemId   uuid.UUID
	Title    string
	Quantity int
	UserId   uuid.UUID
}
