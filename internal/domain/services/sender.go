package services

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/models"
	"log"
	"time"
)

type SenderService struct {
	Sender    interfaces.EventSender
	TaskQueue interfaces.TaskQueue
	QName     string
}

func (s *SenderService) SendNotification(ctx context.Context, event *models.Event) error {
	log.Printf("Sending Notification %v", event)
	err := s.Sender.SendEvent(ctx, event)
	return err
}

func (s *SenderService) Serve(ctx context.Context) error {
	tick := time.Tick(5 * time.Second)
	for {
		<-tick
		err := s.TaskQueue.ConsumeTasksFromQueue(ctx, s.QName, "", false, s.SendNotification)
		if err != nil {
			log.Printf("Can't send notification: %s", err)
		}
	}
}
