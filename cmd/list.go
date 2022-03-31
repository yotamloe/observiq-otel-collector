/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/observiq/observiq-otel-collector/factories"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		baseFactories, err := factories.DefaultFactories()
		if err != nil {
			log.Fatalf("Error registering factories: %s", err)
		}

		builder := strings.Builder{}
		builder.WriteString("Receivers:\n")

		for k := range baseFactories.Receivers {
			builder.WriteString(fmt.Sprintf("\t* %s\n", k))
		}

		builder.WriteString("\n\n")
		builder.WriteString("Processors:\n")

		for k := range baseFactories.Processors {
			builder.WriteString(fmt.Sprintf("\t* %s\n", k))
		}

		builder.WriteString("\n\n")
		builder.WriteString("Exporters:\n")

		for k := range baseFactories.Exporters {
			builder.WriteString(fmt.Sprintf("\t* %s\n", k))
		}

		builder.WriteString("\n\n")
		builder.WriteString("Extensions:\n")

		for k := range baseFactories.Extensions {
			builder.WriteString(fmt.Sprintf("\t* %s\n", k))
		}

		fmt.Println(builder.String())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
