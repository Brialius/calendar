package main

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const tsLayout = "2006-01-02T15:04:05"
const ReqTimeout = time.Second * 10

var RootCmd = &cobra.Command{
	Use:       "client [add, delete, update, list]",
	Short:     "Run gRPC client",
	ValidArgs: []string{"add", "delete", "update", "list", "get", "del", "upd", "ls"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		grpcConfig = getGrpcClientConfig()
		ctx, cancel := context.WithTimeout(context.Background(), ReqTimeout)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("owner", grpcConfig.Owner))
		grpcClient = getGrpcClient(ctx, grpcConfig)
		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGINT)
			<-stop
			log.Printf("Interrupt signal")
			cancel()
		}()
		switch args[0] {
		case "add":
			runCreateRequest(ctx)
		case "delete":
			runDeleteRequest(ctx)
		case "del":
			runDeleteRequest(ctx)
		case "update":
			runUpdateRequest(ctx)
		case "upd":
			runUpdateRequest(ctx)
		case "list":
			runListRequest(ctx)
		case "ls":
			runListRequest(ctx)
		case "get":
			runGetRequest(ctx)
		}
	},
}

var grpcConfig *config.GrpcClientConfig
var grpcClient api.CalendarServiceClient

func getGrpcClient(ctx context.Context, conf *config.GrpcClientConfig) api.CalendarServiceClient {
	if _, err := strconv.Atoi(conf.Port); err != nil {
		log.Fatal(err)
	}
	server := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	conn, err := grpc.DialContext(ctx, server, grpc.WithInsecure(), grpc.WithUserAgent("calendar client"))
	if err != nil {
		log.Fatal(err)
	}
	return api.NewCalendarServiceClient(conn)
}

func getGrpcClientConfig() *config.GrpcClientConfig {
	return config.GetGrpcClientConfig()
}

func init() {
	cobra.OnInitialize(config.SetConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Config file location")
	RootCmd.Flags().StringP("id", "i", "", "event id")
	RootCmd.Flags().StringP("title", "t", "", "event title")
	RootCmd.Flags().StringP("body", "b", "", "event text body")
	RootCmd.Flags().StringP("owner", "o", "", "event owner")
	RootCmd.Flags().StringP("start-time", "s", "", "event start time, format: "+tsLayout)
	RootCmd.Flags().StringP("end-time", "e", "", "event end time, format: "+tsLayout)
	RootCmd.Flags().StringP("host", "n", "", "host name")
	RootCmd.Flags().IntP("port", "p", 0, "port to listen")
	// bind flags to viper
	_ = viper.BindPFlag("id", RootCmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("title", RootCmd.Flags().Lookup("title"))
	_ = viper.BindPFlag("body", RootCmd.Flags().Lookup("body"))
	_ = viper.BindPFlag("owner", RootCmd.Flags().Lookup("owner"))
	_ = viper.BindPFlag("start-time", RootCmd.Flags().Lookup("start-time"))
	_ = viper.BindPFlag("end-time", RootCmd.Flags().Lookup("end-time"))
	_ = viper.BindPFlag("grpc-cli-host", RootCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-cli-port", RootCmd.Flags().Lookup("port"))
	viper.Set("ts-layout", tsLayout)
}

var (
	version = "dev"
	build   = "local"
)

func main() {
	log.Printf("Started calendar gRPC client %s-%s", version, build)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
