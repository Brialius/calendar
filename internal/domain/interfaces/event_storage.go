package interfaces

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/models"
	"time"
)

type EventStorage interface {
	SaveEvent(ctx context.Context, event *models.Event) error
	GetEventByIdOwner(ctx context.Context, id, owner string) (*models.Event, error)
	GetEventsForNotification(ctx context.Context, startTime time.Time, period time.Duration) ([]*models.Event, error)
	GetEventsByOwnerStartDate(ctx context.Context, owner string, startTime *time.Time) ([]*models.Event, error)
	GetEventsCountByOwnerStartDateEndDate(ctx context.Context, owner string, startTime, endTime *time.Time) (int, error)
	DeleteEventByIdOwner(ctx context.Context, id, owner string) error
	UpdateEventByIdOwner(ctx context.Context, id string, event *models.Event) error
	MarkEventNotified(ctx context.Context, id string) error
}
