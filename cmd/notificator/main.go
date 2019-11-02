package main

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/maindb"
	"github.com/Brialius/calendar/internal/mainmq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func constructNotificator(storage interfaces.EventStorage, taskQueue interfaces.TaskQueue,
	period time.Duration, qName, exchange string) *services.NotificatorService {
	return &services.NotificatorService{
		EventStorage: storage,
		TaskQueue:    taskQueue,
		Period:       period,
		QName:        qName,
		Exchange:     exchange,
	}
}

func selectStorage(storageType, dsn string) (interfaces.EventStorage, error) {
	if storageType == "pg" {
		eventStorage, err := maindb.NewPgEventStorage(dsn)
		return eventStorage, err
	}

	return nil, errors.Errorf("storage `%s` is not implemented", storageType)
}

var RootCmd = &cobra.Command{
	Use:   "notificator",
	Short: "Run notificator service",
	Run: func(cmd *cobra.Command, args []string) {
		mqConf := config.GetMqConfig()
		storageConfig := config.GetStorageConfig()
		ctx, cancel := context.WithCancel(context.Background())

		var isAbsentParam bool
		if mqConf.Url == "" {
			isAbsentParam = true
			log.Println("MQ URL is not set")
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

		tq, err := mainmq.NewRabbitMqQueue(mqConf.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer tq.Close(ctx)

		storage, err := selectStorage(storageConfig.StorageType, storageConfig.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer storage.Close(ctx)

		nt := constructNotificator(storage, tq, 24*time.Hour, "notification.tasks", "calendar")
		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGINT)
			<-stop
			log.Printf("Interrupt signal")
			cancel()
		}()
		err = nt.ServeNotificator(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
	Aliases: []string{"nt"},
}

func init() {
	cobra.OnInitialize(config.SetConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Config file location")
	RootCmd.Flags().StringP("url", "u", "", "amqp connection url")
	RootCmd.Flags().StringP("dsn", "d", "", "database connection string")
	RootCmd.Flags().StringP("storage", "s", "", "storage type")
	_ = viper.BindPFlag("dsn", RootCmd.Flags().Lookup("dsn"))
	_ = viper.BindPFlag("storage", RootCmd.Flags().Lookup("storage"))
	_ = viper.BindPFlag("amqp-url", RootCmd.Flags().Lookup("url"))
}

var (
	version = "dev"
	build   = "local"
)

func main() {
	log.Printf("Started calendar notificator service %s-%s", version, build)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
