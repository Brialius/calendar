package cmd

import (
	"github.com/Brialius/calendar/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar is a calendar micorservice",
}

func init() {
	cobra.OnInitialize(config.SetConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("config", "c", "", "Config file location")
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}
