package cmd

import (
	"context"
	"github.com/Brialius/calendar/internal/grpc/api"
	"log"
)

func runDeleteRequest() {
	if grpcConfig.Id == "" {
		log.Fatal("Id is not set")
	}
	req := &api.DeleteEventRequest{
		Id: grpcConfig.Id,
	}
	resp, err := grpcClient.DeleteEvent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.GetError() != "" {
		log.Fatal(resp.GetError())
	}
}
