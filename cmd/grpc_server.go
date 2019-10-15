package cmd

import (
	"fmt"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/grpc/api"
	"github.com/Brialius/calendar/internal/maindb"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
)

func construct(eventStorage interfaces.EventStorage) (*api.CalendarServer, error) {
	eventService := &services.EventService{
		EventStorage: eventStorage,
	}
	server := &api.CalendarServer{
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
		serverConfig := config.GetGrpcServerConfig(cmd)
		storageConfig := config.GetStorageConfig(cmd)
		storage, err := selectStorage(storageConfig.StorageType, storageConfig.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		server, err := construct(storage)
		if err != nil {
			log.Fatal(err)
		}
		addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
		err = server.Serve(addr)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	GrpcServerCmd.Flags().StringP("host", "n", "localhost", "host name")
	GrpcServerCmd.Flags().StringP("port", "p", "8080", "port to listen")
	GrpcServerCmd.Flags().StringP("dsn", "d", "host=127.0.0.1 user=event_user password=event_pwd dbname=event_db", "database connection string")
	GrpcServerCmd.Flags().StringP("storage", "s", "pg", "storage type")
}
