package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Event struct {
	Id        uuid.UUID
	Owner     string
	Title     string
	Text      string
	StartTime *time.Time `db:"start_time"`
	EndTime   *time.Time `db:"end_time"`
}
