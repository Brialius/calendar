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
	Exchange     string
}

func (n *NotificatorService) ScanEvents(ctx context.Context) error {
	events, err := n.EventStorage.GetEventsForNotification(ctx, time.Now(), n.Period)
	if err != nil {
		log.Printf("can't get events for notifications for period `%s`: %s", n.Period, err)
		return err
	}

	for _, e := range events {
		log.Printf("sending notification to `%s` about event `%s`", e.Owner, e.Id)
		if err := n.TaskQueue.SendTaskToQueue(ctx, n.Exchange, n.QName, e); err != nil {
			log.Printf("can't publish notification to task queue: %s", err)
			break
		}
		if err := n.EventStorage.MarkEventNotified(ctx, e.Id.String()); err != nil {
			log.Printf("can't mark event `%s` as notified: %s", e.Id, err)
		}
	}

	return nil
}

func (n *NotificatorService) ServeNotificator(ctx context.Context) error {
	err := n.TaskQueue.DeclareQueue(ctx, n.QName, false)
	if err != nil {
		log.Printf("can't declare task quueue `%s`: %s", n.QName, err)
		return err
	}
	err = n.TaskQueue.DeclareExchange(ctx, n.Exchange, "fanout", true)
	if err != nil {
		log.Printf("can't declare task exchange `%s`: %s", n.Exchange, err)
		return err
	}
	err = n.TaskQueue.BindQueue(ctx, n.QName, n.QName, n.Exchange, false)
	if err != nil {
		log.Printf("can't declare task exchange `%s`: %s", n.Exchange, err)
		return err
	}
	err = n.TaskQueue.SetQos(ctx, 1, 0, false)
	if err != nil {
		log.Printf("can't set QoS for MQ channel: %s", err)
		return err
	}
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
