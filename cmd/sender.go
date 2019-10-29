package cmd

import (
	"context"
	"github.com/Brialius/calendar/internal/config"
	"github.com/Brialius/calendar/internal/domain/interfaces"
	"github.com/Brialius/calendar/internal/domain/services"
	"github.com/Brialius/calendar/internal/mainmq"
	"github.com/Brialius/calendar/internal/mainsender"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func constructSender(taskQueue interfaces.TaskQueue, qName string, sender interfaces.EventSender) *services.SenderService {
	return &services.SenderService{
		TaskQueue: taskQueue,
		QName:     qName,
		Sender:    sender,
	}
}

var SenderCmd = &cobra.Command{
	Use:   "send",
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
		s := constructSender(tq, "notification.tasks", sender)
		err = s.Serve(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
	Aliases: []string{"sn"},
}

func init() {
	RootCmd.AddCommand(SenderCmd)
	SenderCmd.Flags().StringP("url", "u", "", "amqp connection url")
	_ = viper.BindPFlag("amqp-url", SenderCmd.Flags().Lookup("url"))
}
