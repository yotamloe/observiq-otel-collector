/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/observiq/observiq-otel-collector/factories"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/collector/component"
	"gopkg.in/yaml.v3"
)

var (
	listYaml bool
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

		if listYaml {
			printYaml(&baseFactories)
		} else {
			printString(&baseFactories)
		}

	},
}

func printYaml(baseFactories *component.Factories) {
	configMap := make(map[string][]string)

	receivers := make([]string, 0, len(baseFactories.Receivers))
	for k := range baseFactories.Receivers {
		receivers = append(receivers, string(k))
	}
	sort.Strings(receivers)
	configMap["receivers"] = receivers

	processors := make([]string, 0, len(baseFactories.Processors))
	for k := range baseFactories.Processors {
		processors = append(processors, string(k))
	}
	sort.Strings(processors)
	configMap["processors"] = processors

	exporters := make([]string, 0, len(baseFactories.Exporters))
	for k := range baseFactories.Exporters {
		exporters = append(exporters, string(k))
	}
	sort.Strings(exporters)
	configMap["exporters"] = exporters

	extensions := make([]string, 0, len(baseFactories.Extensions))
	for k := range baseFactories.Extensions {
		extensions = append(extensions, string(k))
	}
	sort.Strings(extensions)
	configMap["extensions"] = extensions

	data, err := yaml.Marshal(configMap)
	if err != nil {
		log.Fatalf("Error during yaml format: %s\n", err)
	}

	fmt.Println(string(data))
}

func printString(baseFactories *component.Factories) {
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
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listYaml, "yaml", "y", false, "Prints list in yaml")
}
