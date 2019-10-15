package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

type StorageConfig struct {
	Dsn         string
	StorageType string
}

func GetStorageConfig(cmd *cobra.Command) *StorageConfig {
	viper.AutomaticEnv()
	_ = viper.BindPFlag("dsn", cmd.Flags().Lookup("dsn"))
	_ = viper.BindPFlag("storage", cmd.Flags().Lookup("storage"))
	storageValidation()
	return newDbConfig()
}

func storageValidation() {
	if len(viper.GetString("dsn")) == 0 {
		log.Fatalf("dsn parameter is required")
	}
	if len(viper.GetString("storage")) == 0 {
		log.Fatalf("storage parameter is required")
	}
}

func newDbConfig() *StorageConfig {
	return &StorageConfig{
		Dsn:         viper.GetString("dsn"),
		StorageType: viper.GetString("storage"),
	}
}
