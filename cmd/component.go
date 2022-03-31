/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/observiq/observiq-otel-collector/configscraper"
	"github.com/observiq/observiq-otel-collector/factories"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/collector/config"
	"gopkg.in/yaml.v3"
)

var (
	compName   string
	genMeta    bool
	baseConfig bool
)

// componentCmd represents the component command
var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if compName == "" {
			log.Fatal("Must specify component name")
		}

		baseFactories, err := factories.DefaultFactories()
		if err != nil {
			log.Fatalf("Error registering factories: %s", err)
		}

		params, err := configscraper.GetConfigMeta(config.Type(compName), &baseFactories)
		if err != nil {
			log.Fatalln(err)
		}

		if genMeta && baseConfig {
			log.Fatalln("Can not specify both metadata and baseConfig")
		}

		switch {
		case genMeta:
			data, err := yaml.Marshal(params)
			if err != nil {
				log.Fatalf("Error processing component: %s\n", err)
			}

			fmt.Println(string(data))
		case baseConfig:
			baseConfigMap := make(map[string]interface{})
			for _, param := range params {
				if param.Required {
					defaultVal := param.DefaultValue
					if defaultVal == nil {
						defaultVal = configscraper.ParamTypeDefault(param.Type)
					}

					baseConfigMap[param.Name] = defaultVal
				}
			}

			data, err := yaml.Marshal(baseConfigMap)
			if err != nil {
				log.Fatalf("Error processing component: %s\n", err)
			}

			fmt.Println(string(data))
		}
	},
}

func init() {
	rootCmd.AddCommand(componentCmd)

	componentCmd.Flags().StringVarP(&compName, "component", "c", "", "Component to find config for")
	componentCmd.Flags().BoolVarP(&genMeta, "metadata", "m", false, "Prints metadata of component config")
	componentCmd.Flags().BoolVarP(&baseConfig, "baseConfig", "b", false, "Prints minimum base config")

}
