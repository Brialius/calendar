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
	"github.com/streadway/amqp"
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
	mq               *mqStruct
}

type mqStruct struct {
	server   string
	queue    string
	routeKey string
	exchange string
	ch       *amqp.Channel
	conn     *amqp.Connection
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
	if err != nil {
		return err
	}
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
		return fmt.Errorf("events are different: %v != %v", a.eventToVerify, a.getResponse.GetEvent())
	}
	return nil
}

func (a *apiStruct) iUpdateCreatedEvent(eventJSON *gherkin.DocString) error {
	event := &models.Event{}
	err := json.Unmarshal([]byte(eventJSON.Content), event)
	if err != nil {
		return err
	}
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

func (a *apiStruct) iDeleteEventByPreviousId() error {
	return a.iDeleteEventById(a.eventToVerify.Id)
}

func (a *apiStruct) iDeleteEventById(id string) (err error) {
	a.deleteResponse, err = a.apiCli.DeleteEvent(ctx, &api.DeleteEventRequest{
		Id: id,
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
	st, err := ptypes.TimestampProto(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
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
		return fmt.Errorf("list contains wrong number of records: %d but expect: %d",
			len(a.listResponse.Events), len(a.createdEventsIds))
	}
	for _, e := range a.listResponse.Events {
		if !contains(a.createdEventsIds, e.Id) {
			return fmt.Errorf("list has wrong event: %s", e.Id)
		}
	}
	return nil
}

func (a *apiStruct) iDeleteAllCreatedEvents() error {
	for _, id := range a.createdEventsIds {
		err := a.iDeleteEventById(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *apiStruct) iPurgeOldEvents() error {
	err := a.iDeleteEventById("")
	if err != nil {
		return err
	}
	// Assume that all events except latest are old
	a.createdEventsIds = []string{a.createdEventsIds[len(a.createdEventsIds)-1]}
	fmt.Printf("a.createdEventsIds: %v", a.createdEventsIds)
	return nil
}

func (a *apiStruct) thereIsMQServer(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	a.mq.conn = conn
	a.mq.ch = ch
	return nil
}

func (a *apiStruct) DeclareQueue(qName string, durable bool) error {
	_, err := a.mq.ch.QueueDeclare(
		qName,
		durable,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (a *apiStruct) BindQueue(qName, routingKey, exchange string, durable bool) error {
	err := a.mq.ch.QueueBind(
		qName,
		routingKey,
		exchange,
		false,
		nil,
	)
	return err
}

func (a *apiStruct) DeclareExchange(name, kind string, durable bool) error {
	err := a.mq.ch.ExchangeDeclare(
		name,
		kind,
		durable,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (a *apiStruct) mQExchange(exchange string) error {
	err := a.mq.ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = a.mq.ch.QueueBind(
		a.mq.queue,
		a.mq.routeKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	_, err = a.mq.ch.QueuePurge(a.mq.queue, false)
	return err
}

func (a *apiStruct) mQRouteKey(key string) error {
	a.mq.routeKey = key
	return nil
}

func (a *apiStruct) mQueue(qName string) error {
	_, err := a.mq.ch.QueueDeclare(
		qName,
		false,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (a *apiStruct) iGetTaskFromQueue() error {
	msg, err := a.mq.ch.Consume(a.mq.queue, "", true, false, false, false, nil)
	e := &models.Event{}

	select {
	case <-time.After(10 * time.Second):
		return errors.New("no messages in queue for 10 seconds")
	case task := <-msg:
		err = json.Unmarshal(task.Body, e)
		if err != nil {
			return err
		}
	}
	event, err := grpcsrv.EventToProto(e)
	if err != nil {
		return err
	}
	a.eventToVerify = event
	return nil
}

func (a *apiStruct) eventShouldBeTheSameAsCreated() error {
	if !reflect.DeepEqual(a.eventToVerify, a.createResponse.GetEvent()) {
		return fmt.Errorf("events are different: %v != %v", a.eventToVerify, a.createResponse.GetEvent())
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &apiStruct{}
	m := &mqStruct{}
	a.mq = m
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
	s.Step(`^I delete all created events$`, a.iDeleteAllCreatedEvents)
	s.Step(`^I purge old events$`, a.iPurgeOldEvents)
	s.Step(`^there is MQ server "([^"]*)"$`, a.thereIsMQServer)
	s.Step(`^MQ exchange "([^"]*)"$`, a.mQExchange)
	s.Step(`^MQ route key "([^"]*)"$`, a.mQRouteKey)
	s.Step(`^MQ queue "([^"]*)"$`, a.mQueue)
	s.Step(`^I get task from queue$`, a.iGetTaskFromQueue)
	s.Step(`^Event should be the same as created$`, a.eventShouldBeTheSameAsCreated)
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
