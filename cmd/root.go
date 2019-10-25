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
	cobra.OnInitialize(config.SetLoggerConfig)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
}
