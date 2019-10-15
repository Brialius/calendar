package services

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/models"
	"github.com/satori/go.uuid"
	"time"
)

type EventService struct {
	EventStorage interfaces.EventStorage
}

func (es *EventService) CreateEvent(ctx context.Context, owner, title, text string, startTime *time.Time, endTime *time.Time) (*models.Event, error) {
	// TODO: persistence, validation
	event := &models.Event{
		Id:        uuid.NewV4(),
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := es.EventStorage.SaveEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (es *EventService) DeleteEvent(ctx context.Context, id string) error {
	// TODO: persistence, validation
	err := es.EventStorage.DeleteEventById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
