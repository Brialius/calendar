package cmd

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/mainmq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func constructNotificator(storage interfaces.EventStorage, taskQueue interfaces.TaskQueue,
	period time.Duration, qName string) *services.NotificatorService {
	return &services.NotificatorService{
		EventStorage: storage,
		TaskQueue:    taskQueue,
		Period:       period,
		QName:        qName,
	}
}

var NotificatorCmd = &cobra.Command{
	Use:   "notify",
	Short: "Run notifier service",
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

		nt := constructNotificator(storage, tq, 24*time.Hour, "notification.tasks")
		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGINT)
			<-stop
			log.Printf("Interrupt signal")
			cancel()
		}()
		err = nt.Serve(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
	Aliases: []string{"nt"},
}

func init() {
	RootCmd.AddCommand(NotificatorCmd)
	NotificatorCmd.Flags().StringP("url", "u", "", "amqp connection url")
	NotificatorCmd.Flags().StringP("dsn", "d", "", "database connection string")
	NotificatorCmd.Flags().StringP("storage", "s", "", "storage type")
	_ = viper.BindPFlag("dsn", NotificatorCmd.Flags().Lookup("dsn"))
	_ = viper.BindPFlag("storage", NotificatorCmd.Flags().Lookup("storage"))
	_ = viper.BindPFlag("amqp-url", NotificatorCmd.Flags().Lookup("url"))
}
