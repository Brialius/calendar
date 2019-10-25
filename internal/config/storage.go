package config

import (
	"github.com/spf13/viper"
	"log"
)

type StorageConfig struct {
	Dsn         string
	StorageType string
}

func GetStorageConfig() *StorageConfig {
	log.Println("Configuring storage...")
	viper.SetDefault("dsn", "host=127.0.0.1 user=event_user password=event_pwd dbname=event_db")
	viper.SetDefault("storage", "pg")
	return newDbConfig()
}

func newDbConfig() *StorageConfig {
	return &StorageConfig{
		Dsn:         viper.GetString("dsn"),
		StorageType: viper.GetString("storage"),
	}
}
