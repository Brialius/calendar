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
	Id        string
	Op        string
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
	_ = viper.BindPFlag("id", cmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("list", cmd.Flags().Lookup("list"))
	_ = viper.BindPFlag("add", cmd.Flags().Lookup("add"))
	_ = viper.BindPFlag("delete", cmd.Flags().Lookup("delete"))
	_ = viper.BindPFlag("update", cmd.Flags().Lookup("update"))
	grpcClientValidation()
	return newGrpcClientConfig(tsLayout)
}

func newGrpcClientConfig(tsLayout string) *GrpcClientConfig {
	port, err := strconv.Atoi(viper.GetString("grpc-cli-port"))
	if err != nil {
		log.Fatalf("port parameter should be a digit")
	}
	var st, et *timestamp.Timestamp
	if viper.GetBool("add") || viper.GetBool("update") || viper.GetBool("list") {
		st, err = parseTs(viper.GetString("start-time"), tsLayout)
		if err != nil {
			log.Fatal(err)
		}
	}
	if viper.GetBool("add") || viper.GetBool("update") {
		et, err = parseTs(viper.GetString("end-time"), tsLayout)
		if err != nil {
			log.Fatal(err)
		}
	}

	var op string
	switch {
	case viper.GetBool("add"):
		op = "add"
	case viper.GetBool("delete"):
		op = "delete"
	case viper.GetBool("update"):
		op = "update"
	case viper.GetBool("list"):
		op = "list"
	default:
		log.Fatalf("Unknown command")
	}

	return &GrpcClientConfig{
		Port:      port,
		Host:      viper.GetString("grpc-cli-host"),
		Title:     viper.GetString("title"),
		Text:      viper.GetString("body"),
		Id:        viper.GetString("id"),
		Op:        op,
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
	var errMsg string
	if viper.GetBool("delete") || viper.GetBool("update") {
		if len(viper.GetString("id")) == 0 {
			errMsg += "id parameter is required\n"
		}
	}
	if viper.GetBool("add") || viper.GetBool("update") {
		if len(viper.GetString("title")) == 0 {
			errMsg += "title parameter is required\n"
		}
		if len(viper.GetString("body")) == 0 {
			errMsg += "body parameter is required\n"
		}
		if len(viper.GetString("start-time")) == 0 {
			errMsg += "start-time parameter is required\n"
		}
		if len(viper.GetString("end-time")) == 0 {
			errMsg += "end-time parameter is required\n"
		}
	}
	if viper.GetBool("list") {
		if len(viper.GetString("start-time")) == 0 {
			errMsg += "start-time parameter is required\n"
		}
	}
	if len(errMsg) > 0 {
		log.Fatalf(errMsg)
	}
}
