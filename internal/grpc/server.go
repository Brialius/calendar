package grpc

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/errors"
	"github.com/Brialius/calendar/internal/domain/models"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/grpc/api"
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
	"time"
)

type CalendarServer struct {
	EventService *services.EventService
}

// implements CalendarServiceServer
func (cs *CalendarServer) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
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
	event, err := cs.EventService.CreateEvent(ctx, &models.Event{
		Owner:     owner,
		Title:     req.GetTitle(),
		Text:      req.GetText(),
		StartTime: &st,
		EndTime:   &et,
	})
	if err != nil {
		log.Printf("Error during event creation: `%s` -  %s", req.GetTitle(), err)
		if berr, ok := err.(errors.EventError); ok {
			resp := &api.CreateEventResponse{
				Result: &api.CreateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("Event created: `%s` -  %s", req.GetTitle(), event.Id)
	protoEvent, err := EventToProto(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &api.CreateEventResponse{
		Result: &api.CreateEventResponse_Event{
			Event: protoEvent,
		},
	}
	return resp, nil
}

func EventToProto(event *models.Event) (*api.Event, error) {
	protoEvent := &api.Event{
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

func (cs *CalendarServer) DeleteEvent(ctx context.Context, req *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetId() == "" {
		log.Println("Cleaning up old events..")
		date := time.Now().AddDate(-1, 0, 0)
		err = cs.EventService.DeleteEventsOlderDate(ctx, &date, owner)
		if err != nil {
			if berr, ok := err.(errors.EventError); ok {
				log.Printf("Error during event cleaning up: %s", berr)
				resp := &api.DeleteEventResponse{
					Result: &api.DeleteEventResponse_Error{
						Error: string(berr),
					},
				}
				return resp, nil
			}
			log.Printf("Error during event cleaning up: %s", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &api.DeleteEventResponse{}, nil
	}
	log.Printf("Deleting event: `%s`...", req.GetId())
	err = cs.EventService.DeleteEvent(ctx, req.GetId(), owner)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			log.Printf("Error during event deletion: `%s` -  %s", req.GetId(), berr)
			resp := &api.DeleteEventResponse{
				Result: &api.DeleteEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		}
		log.Printf("Error during event deletion: `%s` -  %s", req.GetId(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if config.Verbose {
		log.Printf("Event Deleted: `%s`", req.GetId())
	}
	return &api.DeleteEventResponse{}, nil
}

func (cs *CalendarServer) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.GetEventResponse, error) {
	log.Printf("Getting event: `%s`...", req.GetId())
	owner, err := getOwner(ctx)
	if err != nil {
		return nil, err
	}
	event, err := cs.EventService.GetEvent(ctx, req.GetId(), owner)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			log.Printf("Error during getting event: `%s` -  %s", req.GetId(), berr)
			resp := &api.GetEventResponse{
				Result: &api.GetEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		}
		log.Printf("Error during getting event: `%s` -  %s", req.GetId(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if config.Verbose {
		log.Printf("Event received: `%s`", req.GetId())
	}
	protoEvent, err := EventToProto(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.GetEventResponse{
		Result: &api.GetEventResponse_Event{Event: protoEvent},
	}, nil
}

func (cs *CalendarServer) ListEvents(ctx context.Context, req *api.ListEventsRequest) (*api.ListEventsResponse, error) {
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
	protoEvents := make([]*api.Event, 0, len(events))
	for _, e := range events {
		protoEvent, err := EventToProto(e)
		if err != nil {
			return nil, err
		}
		protoEvents = append(protoEvents, protoEvent)
	}
	resp := &api.ListEventsResponse{
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

func (cs *CalendarServer) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
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
			resp := &api.UpdateEventResponse{
				Result: &api.UpdateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	protoEvent, err := EventToProto(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &api.UpdateEventResponse{
		Result: &api.UpdateEventResponse_Event{
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
	api.RegisterCalendarServiceServer(s, cs)
	return s.Serve(l)
}
