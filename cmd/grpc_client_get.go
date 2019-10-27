package cmd

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/golang/protobuf/ptypes"
	"log"
)

func runGetRequest(ctx context.Context) {
	if grpcConfig.Id == "" {
		log.Fatal("Id is not set")
	}
	req := &api.GetEventRequest{
		Id: grpcConfig.Id,
	}
	resp, err := grpcClient.GetEvent(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.GetError() != "" {
		log.Fatal(resp.GetError())
	}
	log.Println(printEvent(resp.GetEvent()))
}

func printEvent(event *api.Event) string {
	st, _ := ptypes.Timestamp(event.StartTime)
	et, _ := ptypes.Timestamp(event.EndTime)
	res := fmt.Sprintf(`
**************************
Id: %s
title: %s
From: %s, To: %s
Owner: %s
---
%s
`, event.Id, event.Title, st, et, grpcConfig.Owner, event.Text)
	return res
}
