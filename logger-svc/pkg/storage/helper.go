package storage

import (
	"time"
)

type LogEntry struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Data      string    `bson:"data"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// type LogEntry struct {
// 	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
// 	Name      string    `bson:"name" json:"name"`
// 	Data      string    `bson:"data" json:"data"`
// 	CreatedAt time.Time `bson:"created_at" json:"created_at"`
// 	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
// }
