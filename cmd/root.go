/*
Copyright © 2023 MAX PARTENFELDER <maxpartenfelder@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/mxcd/broke/internal/util"
	"github.com/mxcd/go-config/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "broke",
	Short: "broke - IAM broker",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	defineCliFlags()
	initConfig()
	applyCliFlags()
	util.InitLogger()
	config.Print()
}

func defineCliFlags() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display more verbose output in console output. (default: false)")
	rootCmd.PersistentFlags().BoolP("very-verbose", "", false, "Display even more verbose output in console output. (default: false)")
}

func applyCliFlags() {
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag != nil {
		if verboseFlag.Value.String() == "true" {
			config.Set().String("LOG_LEVEL", "debug")
		}
	}

	veryVerboseFlag := rootCmd.PersistentFlags().Lookup("very-verbose")
	if veryVerboseFlag != nil {
		if veryVerboseFlag.Value.String() == "true" {
			config.Set().String("LOG_LEVEL", "trace")
		}
	}
}

func initConfig() error {
	err := config.LoadConfig([]config.Value{
		config.String("LOG_LEVEL").NotEmpty().Default("info"),

		config.String("IAM_USERNAME").Default(""),
		config.String("IAM_PASSWORD").Default("").Sensitive(),

		config.Bool("DEV").Default(false),
	})
	return err
}
