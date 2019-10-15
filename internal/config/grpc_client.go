package config

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

type GrpcClientConfig struct {
	Port      int
	Host      string
	Title     string
	Text      string
	StartTime *timestamp.Timestamp
	EndTime   *timestamp.Timestamp
}

func parseTs(s, tsLayout string) (*timestamp.Timestamp, error) {
	t, err := time.Parse(tsLayout, s)
	if err != nil {
		return nil, err
	}
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func GetGrpcClientConfig(cmd *cobra.Command, tsLayout string) *GrpcClientConfig {
	viper.AutomaticEnv()
	_ = viper.BindPFlag("grpc-cli-host", cmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("grpc-cli-port", cmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("title", cmd.Flags().Lookup("title"))
	_ = viper.BindPFlag("body", cmd.Flags().Lookup("body"))
	_ = viper.BindPFlag("start-time", cmd.Flags().Lookup("start-time"))
	_ = viper.BindPFlag("end-time", cmd.Flags().Lookup("end-time"))
	grpcClientValidation()
	return newGrpcClientConfig(tsLayout)
}

func newGrpcClientConfig(tsLayout string) *GrpcClientConfig {
	port, err := strconv.Atoi(viper.GetString("grpc-cli-port"))
	if err != nil {
		log.Fatalf("port parameter should be a digit")
	}
	st, err := parseTs(viper.GetString("start-time"), tsLayout)
	if err != nil {
		log.Fatal(err)
	}
	et, err := parseTs(viper.GetString("end-time"), tsLayout)
	if err != nil {
		log.Fatal(err)
	}
	return &GrpcClientConfig{
		Port:      port,
		Host:      viper.GetString("grpc-cli-host"),
		Title:     viper.GetString("title"),
		Text:      viper.GetString("body"),
		StartTime: st,
		EndTime:   et,
	}
}

func grpcClientValidation() {
	if viper.GetString("grpc-cli-host") == "" {
		log.Fatalf("host parameter is required")
	}
	if len(viper.GetString("grpc-cli-port")) == 0 {
		log.Fatalf("host parameter is required")
	}
	if len(viper.GetString("title")) == 0 {
		log.Fatalf("title parameter is required")
	}
	if len(viper.GetString("body")) == 0 {
		log.Fatalf("body parameter is required")
	}
	if len(viper.GetString("start-time")) == 0 {
		log.Fatalf("start-time parameter is required")
	}
	if len(viper.GetString("end-time")) == 0 {
		log.Fatalf("end-time parameter is required")
	}
}
