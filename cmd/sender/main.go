package main

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/mainmq"
	"github.com/Brialius/calendar/internal/mainsender"
	"github.com/Brialius/calendar/internal/monitoring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func constructSender(taskQueue interfaces.TaskQueue,
	qName string, sender interfaces.EventSender) *services.SenderService {
	return &services.SenderService{
		TaskQueue: taskQueue,
		QName:     qName,
		Sender:    sender,
	}
}

var RootCmd = &cobra.Command{
	Use:   "sender",
	Short: "Run sender service",
	Run: func(cmd *cobra.Command, args []string) {
		mqConf := config.GetMqConfig()
		ctx, cancel := context.WithCancel(context.Background())

		if mqConf.Url == "" {
			log.Println("MQ URL is not set")
		}

		tq, err := mainmq.NewRabbitMqQueue(mqConf.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer tq.Close(ctx)

		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGINT)
			<-stop
			log.Printf("Interrupt signal")
			cancel()
		}()
		sender, err := mainsender.NewSendToStream(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
		s := constructSender(tq, "notification.tasks", sender)
		m := &monitoring.PrometheusService{
			Port: viper.GetString("metrics-port"),
		}
		log.Printf("Starting monitoring server on %s...", m.Port)
		m.Serve()
		err = s.Serve(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cobra.OnInitialize(config.SetConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Config file location")
	RootCmd.PersistentFlags().StringP("metrics-port", "m", "9002", "Port for metrics server")
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("metrics-port", RootCmd.PersistentFlags().Lookup("metrics-port"))
	RootCmd.Flags().StringP("url", "u", "", "amqp connection url")
	_ = viper.BindPFlag("amqp-url", RootCmd.Flags().Lookup("url"))
}

var (
	version = "dev"
	build   = "local"
)

func main() {
	log.Printf("Started calendar sender service %s-%s", version, build)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
