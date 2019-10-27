package config

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/viper"
	"log"
	"time"
)

type GrpcClientConfig struct {
	Port      string
	Host      string
	Title     string
	Text      string
	Id        string
	Owner     string
	StartTime string
	EndTime   string
	TsLayout  string
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

func GetGrpcClientConfig() *GrpcClientConfig {
	log.Println("Configuring client...")
	viper.SetDefault("id", "")
	viper.SetDefault("title", "")
	viper.SetDefault("body", "")
	viper.SetDefault("start-time", "")
	viper.SetDefault("end-time", "")
	viper.SetDefault("owner", "user")
	viper.SetDefault("grpc-cli-host", "localhost")
	viper.SetDefault("grpc-cli-port", "8080")
	return newGrpcClientConfig()
}

func (c *GrpcClientConfig) GetStartTime() (*timestamp.Timestamp, error) {
	return parseTs(c.StartTime, c.TsLayout)
}

func (c *GrpcClientConfig) GetEndTime() (*timestamp.Timestamp, error) {
	return parseTs(c.EndTime, c.TsLayout)
}

func newGrpcClientConfig() *GrpcClientConfig {
	return &GrpcClientConfig{
		Port:      viper.GetString("grpc-cli-port"),
		Host:      viper.GetString("grpc-cli-host"),
		Title:     viper.GetString("title"),
		Text:      viper.GetString("body"),
		Id:        viper.GetString("id"),
		Owner:     viper.GetString("owner"),
		StartTime: viper.GetString("start-time"),
		EndTime:   viper.GetString("end-time"),
		TsLayout:  viper.GetString("ts-layout"),
	}
}
