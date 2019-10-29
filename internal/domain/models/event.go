package models

import (
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type Event struct {
	Id        uuid.UUID
	Owner     string
	Title     string
	Text      string
	Notified  bool
	StartTime *time.Time `db:"start_time"`
	EndTime   *time.Time `db:"end_time"`
}

func (e Event) String() string {
	return fmt.Sprintf(`
**************************
Id: %s
title: %s
From: %s, To: %s
Owner: %s
---
%s
`, e.Id, e.Title, e.StartTime, e.EndTime, e.Owner, e.Text)
}
