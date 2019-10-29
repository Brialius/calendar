package interfaces

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/models"
)

type EventSender interface {
	SendEvent(ctx context.Context, event *models.Event) error
}
