package services

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"log"
	"time"
)

type NotificatorService struct {
	EventStorage interfaces.EventStorage
	TaskQueue    interfaces.TaskQueue
}

func (n *NotificatorService) ScanEvents(ctx context.Context, period time.Duration, qName string) error {
	err := n.TaskQueue.DeclareQueue(ctx, qName, true)
	if err != nil {
		log.Printf("can't declare task quueue `%s`: %s", qName, err)
		return err
	}
	err = n.TaskQueue.SetQos(ctx, 1, 0, false)
	if err != nil {
		log.Printf("can't set QoS for MQ channel: %s", err)
		return err
	}
	events, err := n.EventStorage.GetEventsForNotification(ctx, time.Now(), period)
	if err != nil {
		log.Printf("can't get events for notifications for period `%s`: %s", period, err)
		return err
	}
	for _, e := range events {
		err = n.TaskQueue.SendTaskToQueue(ctx, qName, e)
		if err != nil {
			log.Printf("can't publish notification to task queue: %s", err)
			return err
		}
		err = n.EventStorage.MarkEventNotified(ctx, e.Id.String())
		if err != nil {
			log.Printf("can't mark event `%s` as notified: %s", e.Id, err)
			return err
		}
	}
	return nil
}
