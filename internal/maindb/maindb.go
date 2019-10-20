package maindb

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

// implements domain.interfaces.EventStorage
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

func (pges *PgEventStorage) GetEventById(ctx context.Context, id string) (*models.Event, error) {
	// TODO
	return nil, nil
}

func (pges *PgEventStorage) GetEventsByOwnerStartDate(ctx context.Context, owner string, startTime time.Time) ([]*models.Event, error) {
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

func (pges *PgEventStorage) DeleteEventById(ctx context.Context, id string) error {
	query := `
		DELETE FROM events WHERE id=$1
	`
	_, err := pges.db.ExecContext(ctx, query, id)
	return err
}

func (pges *PgEventStorage) UpdateEventById(ctx context.Context, id string, event *models.Event) error {
	// TODO
	return nil
}
