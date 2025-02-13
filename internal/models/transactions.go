package models

import "github.com/google/uuid"

type History struct {
	//TableName  struct{} `pg:"histories"`
	Id         *uuid.UUID
	FromUserId uuid.UUID
	ToUserId   uuid.UUID
	Amount     int
}
