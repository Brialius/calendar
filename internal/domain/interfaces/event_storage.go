package interfaces

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/models"
	"time"
)

type EventStorage interface {
	SaveEvent(ctx context.Context, event *models.Event) error
	GetEventById(ctx context.Context, id string) (*models.Event, error)
	GetEventsByOwnerStartDate(ctx context.Context, owner string, startTime time.Time) []*models.Event
	DeleteEventById(ctx context.Context, id string) error
	UpdateEventById(ctx context.Context, id string, event *models.Event) error
}
