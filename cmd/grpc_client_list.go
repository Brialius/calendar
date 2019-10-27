package cmd

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/golang/protobuf/ptypes"
	"log"
)

func runListRequest(ctx context.Context) {
	isAbsentParam := false
	if grpcConfig.Owner == "" {
		isAbsentParam = true
		log.Println("Owner is not set")
	}
	if grpcConfig.StartTime == "" {
		grpcConfig.StartTime = "0001-01-01T00:00:00"
		log.Println("List all events")
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
	resp, err := grpcClient.ListEvents(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(printEventsList(resp.GetEvents()))
}

func printEventsList(events []*api.Event) string {
	var res string
	for _, e := range events {
		st, _ := ptypes.Timestamp(e.StartTime)
		et, _ := ptypes.Timestamp(e.EndTime)
		res += fmt.Sprintf(`
**************************
Id: %s
title: %s
From: %s, To: %s
Owner: %s
---
%s
`, e.Id, e.Title, st, et, grpcConfig.Owner, e.Text)
	}
	return res
}
