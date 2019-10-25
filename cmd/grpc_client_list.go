package cmd

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/grpc/api"
	"log"
)

func runListRequest() {
	isAbsentParam := false
	if grpcConfig.Owner == "" {
		isAbsentParam = true
		log.Println("Owner is not set")
	}
	if grpcConfig.StartTime == "" {
		grpcConfig.StartTime = "0000-01-01T00:00:00"
		log.Println("Set StartTime to default")
	}
	if isAbsentParam {
		log.Fatal("Some parameters is not set")
	}
	st, err := grpcConfig.GetStartTime()
	if err != nil {
		log.Fatal(err)
	}
	req := &api.ListEventsRequest{
		Owner:     grpcConfig.Owner,
		StartTime: st,
	}
	resp, err := grpcClient.ListEvents(context.Background(), req)
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
