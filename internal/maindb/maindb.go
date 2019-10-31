package maindb

import (
	"context"
	"database/sql"
	"github.com/Brialius/calendar/internal/domain/errors"
	"github.com/Brialius/calendar/internal/domain/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type PgEventStorage struct {
	db *sqlx.DB
}

func NewPgEventStorage(dsn string) (*PgEventStorage, error) {
	db, err := sqlx.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PgEventStorage{db: db}, nil
}

func (pges *PgEventStorage) SaveEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events(id, owner, title, text, start_time, end_time)
		VALUES (:id, :owner, :title, :text, :start_time, :end_time)
	`
	_, err := pges.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         event.Id.String(),
		"owner":      event.Owner,
		"title":      event.Title,
		"text":       event.Text,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
	})
	return err
}

func (pges *PgEventStorage) GetEventByIdOwner(ctx context.Context, id, owner string) (*models.Event, error) {
	query := `
		SELECT * FROM events WHERE id=$1 AND owner=$2
`
	event := &models.Event{}
	err := pges.db.GetContext(ctx, event, query, id, owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return event, nil
}

func (pges *PgEventStorage) GetEventsByOwnerStartDate(ctx context.Context, owner string, startTime *time.Time) ([]*models.Event, error) {
	query := `
		SELECT * FROM events WHERE owner=$1 AND start_time>=$2
`
	var events []*models.Event
	err := pges.db.SelectContext(ctx, &events, query, owner, startTime)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (pges *PgEventStorage) GetEventsForNotification(ctx context.Context, startTime time.Time, period time.Duration) ([]*models.Event, error) {
	query := `
		SELECT * FROM events WHERE start_time<=$1 AND notified = false
`
	var events []*models.Event
	err := pges.db.SelectContext(ctx, &events, query, startTime.Add(period))
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (pges *PgEventStorage) GetEventsCountByOwnerStartDateEndDate(ctx context.Context, owner string, startTime, endTime *time.Time) (int, error) {
	query := `
SELECT count(*)
FROM events
WHERE owner = $1
  AND (start_time BETWEEN $2 AND $3
    OR end_time BETWEEN $2 AND $3)
`
	var eventsCount int
	err := pges.db.GetContext(ctx, &eventsCount, query, owner, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return eventsCount, nil
}

func (pges *PgEventStorage) DeleteEventByIdOwner(ctx context.Context, id, owner string) error {
	query := `
		DELETE FROM events WHERE id=$1 AND owner=$2
	`
	res, err := pges.db.ExecContext(ctx, query, id, owner)
	if res != nil {
		if c, _ := res.RowsAffected(); c == 0 {
			return errors.ErrNotFound
		}
	}
	return err
}

func (pges *PgEventStorage) UpdateEventByIdOwner(ctx context.Context, id string, event *models.Event) error {
	query := `
		UPDATE events SET title=$3, text=$4, start_time=$5, end_time=$6 WHERE id=$1 AND owner=$2
`
	res, err := pges.db.ExecContext(ctx, query, id, event.Owner, event.Title, event.Text, event.StartTime, event.EndTime)
	if res != nil {
		if c, _ := res.RowsAffected(); c == 0 {
			return errors.ErrNotFound
		}
	}
	return err
}

func (pges *PgEventStorage) MarkEventNotified(ctx context.Context, id string) error {
	query := `
		UPDATE events SET notified=true WHERE id=$1
`
	_, err := pges.db.ExecContext(ctx, query, id)
	return err
}

func (pges *PgEventStorage) Close(ctx context.Context) {
	_ = pges.db.Close()
}
