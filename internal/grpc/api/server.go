package api

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/errors"
	"github.com/Brialius/calendar/internal/domain/models"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type CalendarServer struct {
	EventService *services.EventService
}

// implements CalendarServiceServer
func (cs *CalendarServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	log.Printf("Creating new event: `%s`...", req.GetTitle())
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	st, err := ptypes.Timestamp(req.GetStartTime())
	if err != nil {
		log.Printf("start time is incorrect: %s", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	et, err := ptypes.Timestamp(req.GetEndTime())
	if err != nil {
		log.Printf("end time is incorrect: %s", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	event, err := cs.EventService.CreateEvent(ctx, owner, req.GetTitle(), req.GetText(), &st, &et)
	if err != nil {
		log.Printf("Error during event creation: `%s` -  %s", req.GetTitle(), err)
		if berr, ok := err.(errors.EventError); ok {
			resp := &CreateEventResponse{
				Result: &CreateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	log.Printf("Event created: `%s` -  %s", req.GetTitle(), event.Id)
	protoEvent, err := eventToProto(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
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
	}
	var err error
	if protoEvent.StartTime, err = ptypes.TimestampProto(*event.StartTime); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if protoEvent.EndTime, err = ptypes.TimestampProto(*event.EndTime); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return protoEvent, nil
}

func (cs *CalendarServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	log.Printf("Deleting event: `%s`...", req.GetId())
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	err = cs.EventService.DeleteEvent(ctx, req.GetId(), owner)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			log.Printf("Error during event deletion: `%s` -  %s", req.GetId(), berr)
			resp := &DeleteEventResponse{
				Result: &DeleteEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			log.Printf("Error during event deletion: `%s` -  %s", req.GetId(), err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if config.Verbose {
		log.Printf("Event Deleted: `%s`", req.GetId())
	}
	return &DeleteEventResponse{}, nil
}

func (cs *CalendarServer) GetEvent(ctx context.Context, req *GetEventRequest) (*GetEventResponse, error) {
	log.Printf("Getting event: `%s`...", req.GetId())
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	event, err := cs.EventService.GetEvent(ctx, req.GetId(), owner)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			log.Printf("Error during getting event: `%s` -  %s", req.GetId(), berr)
			resp := &GetEventResponse{
				Result: &GetEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			log.Printf("Error during getting event: `%s` -  %s", req.GetId(), err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if config.Verbose {
		log.Printf("Event received: `%s`", req.GetId())
	}
	protoEvent, err := eventToProto(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &GetEventResponse{
		Result: &GetEventResponse_Event{Event: protoEvent},
	}, nil
}

func (cs *CalendarServer) ListEvents(ctx context.Context, req *ListEventsRequest) (*ListEventsResponse, error) {
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	st, err := ptypes.Timestamp(req.GetStartTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Printf("Getting events list: Owner: `%s`, start date: %s ...", owner, st)
	events, err := cs.EventService.ListEvents(ctx, owner, &st)
	if err != nil {
		log.Printf("Error during event list preparing for user: `%s` since:  %s - %s", owner, req.GetStartTime(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("Events list received for user: `%s` since:  %s", owner, st)
	protoEvents := make([]*Event, 0, len(events))
	for _, e := range events {
		protoEvent, err := eventToProto(e)
		if err != nil {
			return nil, err
		}
		protoEvents = append(protoEvents, protoEvent)
	}
	resp := &ListEventsResponse{
		Events: protoEvents,
	}
	return resp, nil
}

func getOwner(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if o := md.Get("owner"); len(o) > 0 {
			return o[0], nil
		}
	}
	return "", status.Errorf(codes.Unauthenticated, "Unauthenticated")
}

func (cs *CalendarServer) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*UpdateEventResponse, error) {
	log.Printf("Updating event: `%s`...", req.GetId())
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	st, err := ptypes.Timestamp(req.GetStartTime())
	if err != nil {
		log.Printf("start time is incorrect: %s", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	et, err := ptypes.Timestamp(req.GetEndTime())
	if err != nil {
		log.Printf("end time is incorrect: %s", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	event, err := cs.EventService.UpdateEvent(ctx, owner, req.GetTitle(), req.GetText(), req.GetId(), &st, &et)
	if err != nil {
		log.Printf("Error during event creation: `%s` -  %s", req.GetTitle(), err)
		if berr, ok := err.(errors.EventError); ok {
			resp := &UpdateEventResponse{
				Result: &UpdateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	protoEvent, err := eventToProto(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &UpdateEventResponse{
		Result: &UpdateEventResponse_Event{
			Event: protoEvent,
		},
	}
	return resp, nil
}

func (cs *CalendarServer) Serve(addr string) error {
	s := grpc.NewServer()
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGINT)
		<-stop
		log.Printf("Interrupt signal")
		log.Printf("Gracefully shutdown")
		s.GracefulStop()
	}()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	RegisterCalendarServiceServer(s, cs)
	return s.Serve(l)
}
