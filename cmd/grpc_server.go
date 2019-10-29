package cmd

import (
	"context"
	"fmt"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/grpc"
	"github.com/Brialius/calendar/internal/maindb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func constructGrpcServer(eventStorage interfaces.EventStorage) (*grpc.CalendarServer, error) {
	eventService := &services.EventService{
		EventStorage: eventStorage,
	}
	server := &grpc.CalendarServer{
		EventService: eventService,
	}
	return server, nil
}

func selectStorage(storageType, dsn string) (interfaces.EventStorage, error) {
	if storageType == "pg" {
		eventStorage, err := maindb.NewPgEventStorage(dsn)
		return eventStorage, err
	}
	return nil, errors.Errorf("storage `%s` is not implemented", storageType)
}

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc_server",
	Short: "Run gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		serverConfig := config.GetGrpcServerConfig()
		storageConfig := config.GetStorageConfig()
		isAbsentParam := false
		if serverConfig.Host == "" {
			isAbsentParam = true
			log.Println("Host is not set")
		}
		if serverConfig.Port == "" {
			isAbsentParam = true
			log.Println("Port is not set")
		}
		if storageConfig.Dsn == "" {
			isAbsentParam = true
			log.Println("Dsn is not set")
		}
		if storageConfig.StorageType == "" {
			isAbsentParam = true
			log.Println("StorageType is not set")
		}
		if isAbsentParam {
			log.Fatal("Some parameters is not set")
		}

		storage, err := selectStorage(storageConfig.StorageType, storageConfig.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer storage.Close(context.Background())

		server, err := constructGrpcServer(storage)
		if err != nil {
			log.Fatal(err)
		}
		addr := fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)
		log.Printf("Starting server on %s...", addr)
		err = server.Serve(addr)
		if err != nil {
			log.Fatal(err)
		}
	},
	Aliases: []string{"gs"},
}

func init() {
	RootCmd.AddCommand(GrpcServerCmd)
	GrpcServerCmd.Flags().StringP("host", "n", "", "host name")
	GrpcServerCmd.Flags().IntP("port", "p", 0, "port to listen")
	GrpcServerCmd.Flags().StringP("dsn", "d", "", "database connection string")
	GrpcServerCmd.Flags().StringP("storage", "s", "", "storage type")
	_ = viper.BindPFlag("grpc-srv-host", GrpcServerCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-srv-port", GrpcServerCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("dsn", GrpcServerCmd.Flags().Lookup("dsn"))
	_ = viper.BindPFlag("storage", GrpcServerCmd.Flags().Lookup("storage"))
}
