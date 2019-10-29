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
	Period       time.Duration
	QName        string
}

func (n *NotificatorService) ScanEvents(ctx context.Context) error {
	err := n.TaskQueue.DeclareQueue(ctx, n.QName, true)
	if err != nil {
		log.Printf("can't declare task quueue `%s`: %s", n.QName, err)
		return err
	}
	err = n.TaskQueue.SetQos(ctx, 1, 0, false)
	if err != nil {
		log.Printf("can't set QoS for MQ channel: %s", err)
		return err
	}
	events, err := n.EventStorage.GetEventsForNotification(ctx, time.Now(), n.Period)
	if err != nil {
		log.Printf("can't get events for notifications for period `%s`: %s", n.Period, err)
		return err
	}

	for _, e := range events {
		err = n.TaskQueue.SendTaskToQueue(ctx, n.QName, e)
		if err != nil {
			log.Printf("can't publish notification to task queue: %s", err)
		}
		err = n.EventStorage.MarkEventNotified(ctx, e.Id.String())
		if err != nil {
			log.Printf("can't mark event `%s` as notified: %s", e.Id, err)
		}
	}

	return nil
}

func (n *NotificatorService) Serve(ctx context.Context) error {
	tick := time.Tick(5 * time.Second)
	for {
		<-tick
		err := n.ScanEvents(ctx)
		if err != nil {
			log.Printf("Error during ScanEvents: %s", err)
			return err
		}
	}
}
