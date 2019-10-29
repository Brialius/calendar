package mainsender

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/domain/models"
	"io"
)

type SendToStream struct {
	out io.Writer
}

func NewSendToStream(out io.Writer) (*SendToStream, error) {
	return &SendToStream{out: out}, nil
}

func (s *SendToStream) SendEvent(ctx context.Context, event *models.Event) error {
	_, err := fmt.Fprintf(s.out, "Send notification to `%s`: %s", event.Owner, event.Id)
	return err
}
