package integration_tests

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/Brialius/calendar/internal/domain/models"
	grpcsrv "github.com/Brialius/calendar/internal/grpc"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"reflect"
	"testing"
	"time"
)

type apiStruct struct {
	ctx              context.Context
	apiCli           api.CalendarServiceClient
	getRequest       *api.GetEventRequest
	getResponse      *api.GetEventResponse
	createRequest    *api.CreateEventRequest
	createResponse   *api.CreateEventResponse
	updateRequest    *api.UpdateEventRequest
	updateResponse   *api.UpdateEventResponse
	deleteRequest    *api.DeleteEventRequest
	deleteResponse   *api.DeleteEventResponse
	listRequest      *api.ListEventsRequest
	listResponse     *api.ListEventsResponse
	eventToVerify    *api.Event
	createdEventsIds []string
}

var ctx context.Context

func (a *apiStruct) thereIsUser(owner string) error {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("owner", owner))
	return nil
}

func (a *apiStruct) thereIsServer(url string) error {
	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithUserAgent("calendar integration_tests"))
	if err != nil {
		return err
	}
	a.apiCli = api.NewCalendarServiceClient(conn)
	if a.apiCli == nil {
		return err
	}
	return nil
}

func (a *apiStruct) iCreateEvent(eventJSON *gherkin.DocString) error {
	event := &models.Event{}
	err := json.Unmarshal([]byte(eventJSON.Content), event)
	eventProto, err := grpcsrv.EventToProto(event)
	if err != nil {
		return err
	}
	a.createResponse, err = a.apiCli.CreateEvent(ctx, &api.CreateEventRequest{
		Title:     eventProto.Title,
		Text:      eventProto.Text,
		StartTime: eventProto.StartTime,
		EndTime:   eventProto.EndTime,
	})
	if err != nil {
		return err
	}
	if errResponce := a.createResponse.GetError(); errResponce != "" {
		return errors.New(errResponce)
	}
	a.eventToVerify = a.createResponse.GetEvent()
	a.createdEventsIds = append(a.createdEventsIds, a.createResponse.GetEvent().Id)
	return nil
}

func (a *apiStruct) iGetEventById() (err error) {
	a.getResponse, err = a.apiCli.GetEvent(ctx, &api.GetEventRequest{
		Id: a.eventToVerify.Id,
	})
	if err != nil {
		return err
	}
	if errResponce := a.getResponse.GetError(); errResponce != "" {
		return errors.New(errResponce)
	}
	return
}

func (a *apiStruct) eventsShouldBeTheSame() error {
	if !reflect.DeepEqual(a.eventToVerify, a.getResponse.GetEvent()) {
		return errors.New("events are different")
	}
	return nil
}

func (a *apiStruct) iUpdateCreatedEvent(eventJSON *gherkin.DocString) error {
	event := &models.Event{}
	err := json.Unmarshal([]byte(eventJSON.Content), event)
	eventProto, err := grpcsrv.EventToProto(event)
	if err != nil {
		return err
	}
	a.updateResponse, err = a.apiCli.UpdateEvent(ctx, &api.UpdateEventRequest{
		Id:        a.createResponse.GetEvent().Id,
		Title:     eventProto.Title,
		Text:      eventProto.Text,
		StartTime: eventProto.StartTime,
		EndTime:   eventProto.EndTime,
	})
	if err != nil {
		return err
	}
	if errResponce := a.createResponse.GetError(); errResponce != "" {
		return errors.New(errResponce)
	}
	a.eventToVerify = a.updateResponse.GetEvent()
	return nil
}

func (a *apiStruct) iDeleteEventByPreviousId() (err error) {
	a.deleteResponse, err = a.apiCli.DeleteEvent(ctx, &api.DeleteEventRequest{
		Id: a.eventToVerify.Id,
	})
	if err != nil {
		return err
	}
	if errResponce := a.deleteResponse.GetError(); errResponce != "" {
		return errors.New(errResponce)
	}
	return
}

func (a *apiStruct) eventByPreviousIdShouldBeAbsent() (err error) {
	a.getResponse, err = a.apiCli.GetEvent(ctx, &api.GetEventRequest{
		Id: a.eventToVerify.Id,
	})
	if err != nil {
		return err
	}
	if errResponce := a.getResponse.GetError(); errResponce != "event not found" {
		return errors.New("event exists")
	}
	return
}

func (a *apiStruct) iGetEventList() (err error) {
	st, _ := ptypes.TimestampProto(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC))
	a.listResponse, err = a.apiCli.ListEvents(ctx, &api.ListEventsRequest{
		StartTime: st,
	})
	if err != nil {
		return err
	}
	if errResponce := a.getResponse.GetError(); errResponce != "event not found" {
		return errors.New("event exists")
	}
	return
}

func (a *apiStruct) eventListShouldContainCreatedEvents() error {
	if len(a.listResponse.Events) != len(a.createdEventsIds) {
		fmt.Printf("a.listResponse.Events: %d, a.createdEventsIds: %d", len(a.listResponse.Events), len(a.createdEventsIds))
		return errors.New("list contains wrong number of records")
	}
	for _, e := range a.listResponse.Events {
		if !contains(a.createdEventsIds, e.Id) {
			return errors.New("list has wrong records")
		}
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &apiStruct{}
	s.BeforeScenario(func(interface{}) {
		a.createdEventsIds = make([]string, 0)
	})
	s.Step(`^there is user "([^"]*)"$`, a.thereIsUser)
	s.Step(`^there is server "([^"]*)"$`, a.thereIsServer)
	s.Step(`^I create event$`, a.iCreateEvent)
	s.Step(`^I get event by previous id$`, a.iGetEventById)
	s.Step(`^Events should be the same$`, a.eventsShouldBeTheSame)
	s.Step(`^I update created event$`, a.iUpdateCreatedEvent)
	s.Step(`^I delete event by previous id$`, a.iDeleteEventByPreviousId)
	s.Step(`^Event by previous id should be absent$`, a.eventByPreviousIdShouldBeAbsent)
	s.Step(`^I get event list$`, a.iGetEventList)
	s.Step(`^Event list should contain created events$`, a.eventListShouldContainCreatedEvents)
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
