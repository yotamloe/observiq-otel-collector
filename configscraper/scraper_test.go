package configscraper

import (
	"fmt"
	"testing"

	"github.com/observiq/observiq-otel-collector/factories"
	"github.com/stretchr/testify/require"
)

func TestScraper(t *testing.T) {
	baseFactories, err := factories.DefaultFactories()
	require.NoError(t, err)

	for k := range baseFactories.Receivers {
		fmt.Println("Receiver:", k)
		out, err := GetConfigMeta(k, &baseFactories)
		require.NoError(t, err, fmt.Sprintf("Failed for receiver: %s", k))
		require.NotNil(t, out)
	}

	for k := range baseFactories.Processors {
		fmt.Println("Processors:", k)
		out, err := GetConfigMeta(k, &baseFactories)
		require.NoError(t, err, fmt.Sprintf("Failed for processor: %s", k))
		require.NotNil(t, out)
	}

	for k := range baseFactories.Exporters {
		fmt.Println("Exporters:", k)
		out, err := GetConfigMeta(k, &baseFactories)
		require.NoError(t, err, fmt.Sprintf("Failed for exporter: %s", k))
		require.NotNil(t, out)
	}

	for k := range baseFactories.Extensions {
		fmt.Println("Extensions:", k)
		out, err := GetConfigMeta(k, &baseFactories)
		require.NoError(t, err, fmt.Sprintf("Failed for extension: %s", k))
		require.NotNil(t, out)
	}

}
