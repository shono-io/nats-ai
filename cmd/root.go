/*
Copyright © 2024 Shono <code@shono.io>
*/
package cmd

import (
  "fmt"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "natsai/app"
  "os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "nats-ai",
  Short: "nats assistance service",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) {
    opts := app.FromViper(viper.GetViper())
    if err := app.Launch(opts...); err != nil {
      log.Panic().Err(err).Msg("failed to launch service")
    }
  },
}

func Execute() {
  err := rootCmd.Execute()
  if err != nil {
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shono.yaml)")

  app.ViperFlags(rootCmd)

  if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
    log.Panic().Err(err).Msg("failed to bind flags")
  }
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := os.UserHomeDir()
    cobra.CheckErr(err)

    // Search config in home directory with name ".nassist" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigType("yaml")
    viper.SetConfigName(".nats-ai")
  }

  viper.SetEnvPrefix("NATSAI")
  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
  }
}
