package interfaces

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/models"
)

type TaskQueue interface {
	DeclareQueue(ctx context.Context, qName string, durable bool) error
	SetQos(ctx context.Context, prefetchCount, prefetchSize int, global bool) error
	SendTaskToQueue(ctx context.Context, qName string, event *models.Event) error
	ConsumeTasksFromQueue(ctx context.Context, qName, consumer string, autoAck bool, task func(event *models.Event) error) error
}
