package configscraper

import (
	"go.opentelemetry.io/collector/component"
)

func scrapeReceiverConfig(factory component.ReceiverFactory) ([]*Paramaeter, error) {
	// Get default values as a map we can lookup
	defaults, err := getDefaultValues(factory.CreateDefaultConfig())
	if err != nil {
		return nil, err
	}

	parameters, err := extractParameters(defaults)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}
