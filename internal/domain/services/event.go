package services

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/models"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type EventService struct {
	EventStorage interfaces.EventStorage
}

func (es *EventService) CreateEvent(ctx context.Context, owner, title, text string, startTime *time.Time, endTime *time.Time) (*models.Event, error) {
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
		log.Printf("can't create event `%s`: %s", event, err)
		return nil, err
	}
	return event, nil
}

func (es *EventService) DeleteEvent(ctx context.Context, id string) error {
	_, err := parseUuid(id)
	if err != nil {
		return err
	}
	err = es.EventStorage.DeleteEventById(ctx, id)
	if err != nil {
		log.Printf("can't delete event `%s`: %s", id, err)
		return err
	}
	return nil
}

func (es *EventService) UpdateEvent(ctx context.Context, owner, title, text, id string, startTime *time.Time, endTime *time.Time) (*models.Event, error) {
	uuidId, err := parseUuid(id)
	event := &models.Event{
		Id:        uuidId,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err = es.EventStorage.SaveEvent(ctx, event)
	if err != nil {
		log.Printf("can't update event `%s`: %s", id, err)
		return nil, err
	}
	return event, nil
}

func (es *EventService) ListEvents(ctx context.Context, owner string, startTime time.Time) ([]*models.Event, error) {
	events, err := es.EventStorage.GetEventsByOwnerStartDate(ctx, owner, startTime)
	if err != nil {
		log.Printf("can't get list of events for owner: `%s` startTime: `%s`: %s", owner, startTime, err)
		return nil, err
	}
	return events, nil
}

func parseUuid(id string) (uuid.UUID, error) {
	uuidId, err := uuid.FromString(id)
	if err != nil {
		log.Printf("can't parse UUID from string: `%s`: %s", id, err)
		return uuidId, err
	}
	return uuidId, nil
}
