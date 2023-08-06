package models

import (
	"time"
)

type Task struct {
	ID        int       `bson:"_id,omitempty" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Status    string    `bson:"status" json:"status"`
	ActiveAt  time.Time `bson:"activeAt" json:"activeAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Counter struct {
	ID       string `bson:"_id"`
	Sequence int    `bson:"sequence"`
}
