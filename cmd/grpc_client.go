package cmd

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
)

const tsLayout = "2006-01-02T15:04:05"

var GrpcClientCmd = &cobra.Command{
	Use:   "grpc_client",
	Short: "Run gRPC client",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetGrpcClientConfig(cmd, tsLayout)
		server := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		conn, err := grpc.Dial(server, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewCalendarServiceClient(conn)
		switch {
		case conf.Op == "add":
			runCreateRequest(conf, client)
		case conf.Op == "delete":
			runDeleteRequest(conf, client)
		case conf.Op == "update":
			runUpdateRequest(conf, client)
		case conf.Op == "list":
			runListRequest(conf, client)
		default:
			log.Fatalf("Unknown command")
		}
	},
}

func runCreateRequest(conf *config.GrpcClientConfig, client api.CalendarServiceClient) {
	req := &api.CreateEventRequest{
		Title:     conf.Title,
		Text:      conf.Text,
		StartTime: conf.StartTime,
		EndTime:   conf.EndTime,
	}
	resp, err := client.CreateEvent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.GetError() != "" {
		log.Fatal(resp.GetError())
	} else {
		log.Println(resp.GetEvent().Id)
	}
}

func runListRequest(conf *config.GrpcClientConfig, client api.CalendarServiceClient) {
	req := &api.ListEventsRequest{
		StartTime: conf.StartTime,
	}
	resp, err := client.ListEvents(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(printEventsList(resp.GetEvents()))
}

func printEventsList(events []*api.Event) string {
	var res string
	for _, e := range events {
		res += fmt.Sprintf(`
**************************
%s (%s)
From: %s, To: %s
Owner: %s
---
%s
---
`, e.Title, e.Id, e.StartTime.String(), e.EndTime.String(), e.Owner, e.Text)
	}
	return res
}

func runUpdateRequest(conf *config.GrpcClientConfig, client api.CalendarServiceClient) {
	req := &api.UpdateEventRequest{
		Id:        conf.Id,
		Title:     conf.Title,
		Text:      conf.Text,
		StartTime: conf.StartTime,
		EndTime:   conf.EndTime,
	}
	resp, err := client.UpdateEvent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.GetError() != "" {
		log.Fatal(resp.GetError())
	} else {
		log.Println(resp.GetEvent())
	}
}

func runDeleteRequest(conf *config.GrpcClientConfig, client api.CalendarServiceClient) {
	req := &api.DeleteEventRequest{
		Id: conf.Id,
	}
	resp, err := client.DeleteEvent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.GetError() != "" {
		log.Fatal(resp.GetError())
	}
}

func init() {
	GrpcClientCmd.Flags().StringP("host", "n", "localhost", "host name")
	GrpcClientCmd.Flags().StringP("port", "p", "8080", "port to listen")
	GrpcClientCmd.Flags().StringP("title", "t", "", "event title")
	GrpcClientCmd.Flags().BoolP("list", "l", true, "events list")
	GrpcClientCmd.Flags().BoolP("add", "a", false, "add event")
	GrpcClientCmd.Flags().BoolP("delete", "d", false, "delete event")
	GrpcClientCmd.Flags().BoolP("update", "u", false, "update event")
	GrpcClientCmd.Flags().StringP("id", "i", "", "event id")
	GrpcClientCmd.Flags().StringP("body", "b", "", "event text body")
	GrpcClientCmd.Flags().StringP("start-time", "s", "", "event start time, format: "+tsLayout)
	GrpcClientCmd.Flags().StringP("end-time", "e", "", "event end time, format: "+tsLayout)
}
