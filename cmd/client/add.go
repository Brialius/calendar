package main

import (
	"context"
	"github.com/Brialius/calendar/internal/grpc/api"
	"log"
)

func runCreateRequest(ctx context.Context) {
	isAbsentParam := false
	if grpcConfig.Title == "" {
		isAbsentParam = true
		log.Println("Title is not set")
	}
	if grpcConfig.Text == "" {
		isAbsentParam = true
		log.Println("Text is not set")
	}
	if grpcConfig.Owner == "" {
		isAbsentParam = true
		log.Println("Owner is not set")
	}
	if grpcConfig.StartTime == "" {
		isAbsentParam = true
		log.Println("StartTime is not set")
	}
	if grpcConfig.EndTime == "" {
		isAbsentParam = true
		log.Println("EndTime is not set")
	}
	if isAbsentParam {
		log.Fatal("Some parameters is not set")
	}
	st, err := grpcConfig.GetStartTime()
	if err != nil {
		log.Fatal(err)
	}
	et, err := grpcConfig.GetEndTime()
	if err != nil {
		log.Fatal(err)
	}
	req := &api.CreateEventRequest{
		Title:     grpcConfig.Title,
		Text:      grpcConfig.Text,
		StartTime: st,
		EndTime:   et,
	}
	resp, err := grpcClient.CreateEvent(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.GetError() != "" {
		log.Fatal(resp.GetError())
	}
	log.Println(resp.GetEvent().Id)
}
