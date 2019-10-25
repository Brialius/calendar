package cmd

import (
	"fmt"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

const tsLayout = "2006-01-02T15:04:05"

var GrpcClientCmd = &cobra.Command{
	Use:       "grpc_client [add, delete, update, list]",
	Short:     "Run gRPC client",
	Aliases:   []string{"gc"},
	ValidArgs: []string{"add", "delete", "update", "list"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		grpcConfig = getGrpcClientConfig()
		grpcClient = getGrpcClient(grpcConfig)
		switch args[0] {
		case "add":
			runCreateRequest()
		case "delete":
			runDeleteRequest()
		case "update":
			runUpdateRequest()
		case "list":
			runListRequest()
		}
	},
}

var grpcConfig *config.GrpcClientConfig
var grpcClient api.CalendarServiceClient

func getGrpcClient(conf *config.GrpcClientConfig) api.CalendarServiceClient {
	if _, err := strconv.Atoi(conf.Port); err != nil {
		log.Fatal(err)
	}
	server := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return api.NewCalendarServiceClient(conn)
}

func getGrpcClientConfig() *config.GrpcClientConfig {
	return config.GetGrpcClientConfig()
}

func init() {
	RootCmd.AddCommand(GrpcClientCmd)
	GrpcClientCmd.Flags().StringP("id", "i", "", "event id")
	GrpcClientCmd.Flags().StringP("title", "t", "", "event title")
	GrpcClientCmd.Flags().StringP("body", "b", "", "event text body")
	GrpcClientCmd.Flags().StringP("owner", "o", "", "event owner")
	GrpcClientCmd.Flags().StringP("start-time", "s", "", "event start time, format: "+tsLayout)
	GrpcClientCmd.Flags().StringP("end-time", "e", "", "event end time, format: "+tsLayout)
	GrpcClientCmd.Flags().StringP("host", "n", "", "host name")
	GrpcClientCmd.Flags().IntP("port", "p", 0, "port to listen")
	// bind flags to viper
	_ = viper.BindPFlag("id", GrpcClientCmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("title", GrpcClientCmd.Flags().Lookup("title"))
	_ = viper.BindPFlag("body", GrpcClientCmd.Flags().Lookup("body"))
	_ = viper.BindPFlag("owner", GrpcClientCmd.Flags().Lookup("owner"))
	_ = viper.BindPFlag("start-time", GrpcClientCmd.Flags().Lookup("start-time"))
	_ = viper.BindPFlag("end-time", GrpcClientCmd.Flags().Lookup("end-time"))
	_ = viper.BindPFlag("grpc-cli-host", GrpcClientCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-cli-port", GrpcClientCmd.Flags().Lookup("port"))
	viper.Set("ts-layout", tsLayout)
}
