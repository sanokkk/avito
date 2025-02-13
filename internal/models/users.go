package models

import "github.com/google/uuid"

type User struct {
	Id           *uuid.UUID
	Username     string
	PasswordHash string
	Salt         []byte
	Coins        int
	Items        []*UserItem `pg:"rel:has-many"`
}
