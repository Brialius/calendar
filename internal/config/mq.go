package config

import (
	"github.com/spf13/viper"
	"log"
)

type MqConfig struct {
	Url string
}

func GetMqConfig() *MqConfig {
	log.Println("Configuring message queue broker...")
	viper.SetDefault("amqp-url", "amqp://queue_user:queue-super-password@localhost:5672/")
	return newMqrConfig()
}

func newMqrConfig() *MqConfig {
	return &MqConfig{
		Url: viper.GetString("amqp-url"),
	}
}
