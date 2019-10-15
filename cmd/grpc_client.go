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
	},
}

func init() {
	GrpcClientCmd.Flags().StringP("host", "n", "localhost", "host name")
	GrpcClientCmd.Flags().StringP("port", "p", "8080", "port to listen")
	GrpcClientCmd.Flags().StringP("title", "t", "", "event title")
	GrpcClientCmd.Flags().StringP("body", "b", "", "event text body")
	GrpcClientCmd.Flags().StringP("start-time", "s", "", "event start time, format: "+tsLayout)
	GrpcClientCmd.Flags().StringP("end-time", "e", "", "event end time, format: "+tsLayout)
}
