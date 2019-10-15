package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type GrpcServerConfig struct {
	Port int
	Host string
}

func GetGrpcServerConfig(cmd *cobra.Command) *GrpcServerConfig {
	viper.AutomaticEnv()
	_ = viper.BindPFlag("grpc-srv-host", cmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-srv-port", cmd.Flags().Lookup("port"))
	grpcServerValidation()
	return newGrpcServerConfig()
}

func newGrpcServerConfig() *GrpcServerConfig {
	port, err := strconv.Atoi(viper.GetString("grpc-srv-port"))
	if err != nil {
		log.Fatalf("port parameter should be a digit")
	}
	return &GrpcServerConfig{
		Port: port,
		Host: viper.GetString("grpc-srv-host"),
	}
}

func grpcServerValidation() {
	if len(viper.GetString("grpc-srv-host")) == 0 {
		log.Fatalf("host parameter is required")
	}
	if len(viper.GetString("grpc-srv-port")) == 0 {
		log.Fatalf("port parameter is required")
	}
}
