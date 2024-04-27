/*
Copyright © 2024 Shono <code@shono.io>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shono-io/mini"
	"nassist/api"
	"nassist/memory"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nassist",
	Short: "nats assistance service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := mini.FromViper(viper.GetViper())
		if err != nil {
			log.Panic().Err(err).Msg("failed to create service")
		}

		if err != nil {
			panic(err)
		}

		d, err := time.ParseDuration(viper.GetString("ttl"))
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse ttl, reversing to default value")
			d = 10 * time.Minute
		}

		ts := memory.NewThreadStore(viper.GetInt("size"), d)

		if err := api.Attach(srv, viper.GetString("endpoint"), ts); err != nil {
			log.Panic().Err(err).Msg("failed to attach api")
		}

		if err := srv.Run(context.Background(), mini.NewIdleWorker()); err != nil {
			log.Panic().Err(err).Msg("failed to run service")
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
	rootCmd.PersistentFlags().StringP("endpoint", "e", "http://brain:11434/api", "the endpoint to your local llm service")
	rootCmd.PersistentFlags().StringP("ttl", "t", "10m", "the amount of time to keep a thread in memory")
	rootCmd.PersistentFlags().IntP("size", "s", 10, "the number of threads to keep in memory")

	mini.ConfigureCommand(rootCmd)

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
		viper.SetConfigName(".shono")
	}

	viper.SetEnvPrefix("SHONO")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
