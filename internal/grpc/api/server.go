package api

import (
	"context"
	"github.com/Brialius/calendar/internal/domain/models"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"net"
	"time"

	"github.com/Brialius/calendar/internal/domain/errors"
	"github.com/Brialius/calendar/internal/domain/services"
)

type CalendarServer struct {
	EventService *services.EventService
}

// implements CalendarServiceServer
func (cs *CalendarServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	owner := "admin"
	if o := ctx.Value("owner"); o != nil {
		owner, _ = o.(string)
	}

	startTime := new(time.Time)
	if req.GetStartTime() != nil {
		st, err := ptypes.Timestamp(req.GetStartTime())
		if err != nil {
			return nil, err
		}
		startTime = &st
	}

	endTime := new(time.Time)
	if req.GetEndTime() != nil {
		et, err := ptypes.Timestamp(req.GetEndTime())
		if err != nil {
			return nil, err
		}
		endTime = &et
	}

	event, err := cs.EventService.CreateEvent(ctx, owner, req.GetTitle(), req.GetText(), startTime, endTime)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			resp := &CreateEventResponse{
				Result: &CreateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, err
		}
	}
	protoEvent, err := eventToProto(event)
	resp := &CreateEventResponse{
		Result: &CreateEventResponse_Event{
			Event: protoEvent,
		},
	}
	return resp, nil
}

func eventToProto(event *models.Event) (*Event, error) {
	protoEvent := &Event{
		Id:    event.Id.String(),
		Title: event.Title,
		Text:  event.Text,
		Owner: event.Owner,
	}
	var err error
	if protoEvent.StartTime, err = ptypes.TimestampProto(*event.StartTime); err != nil {
		return nil, err
	}
	if protoEvent.EndTime, err = ptypes.TimestampProto(*event.EndTime); err != nil {
		return nil, err
	}
	return protoEvent, nil
}

func (cs *CalendarServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	err := cs.EventService.DeleteEvent(ctx, req.GetId())
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			resp := &DeleteEventResponse{
				Result: &DeleteEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, err
		}
	}
	return &DeleteEventResponse{}, nil
}

func (cs *CalendarServer) ListEvents(ctx context.Context, req *ListEventsRequest) (*ListEventsResponse, error) {
	owner := "admin"
	if o := ctx.Value("owner"); o != nil {
		owner, _ = o.(string)
	}
	st, err := ptypes.Timestamp(req.GetStartTime())
	if err != nil {
		return nil, err
	}
	events, err := cs.EventService.ListEvents(ctx, owner, st)
	if err != nil {
		return nil, err
	}
	protoEvents := make([]*Event, 0, len(events))
	for _, e := range events {
		sTime, err := ptypes.TimestampProto(*e.StartTime)
		if err != nil {
			return nil, err
		}
		eTime, err := ptypes.TimestampProto(*e.EndTime)
		if err != nil {
			return nil, err
		}
		protoEvents = append(protoEvents, &Event{
			Id:        e.Id.String(),
			Title:     e.Title,
			Text:      e.Text,
			Owner:     e.Owner,
			StartTime: sTime,
			EndTime:   eTime,
		})
	}
	resp := &ListEventsResponse{
		Events: protoEvents,
	}
	return resp, nil
}

func (cs *CalendarServer) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*UpdateEventResponse, error) {
	// TODO
	return nil, nil
}

func (cs *CalendarServer) Serve(addr string) error {
	s := grpc.NewServer()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	RegisterCalendarServiceServer(s, cs)
	return s.Serve(l)
}
