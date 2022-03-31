package configscraper

import (
	"fmt"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver"
	"github.com/stretchr/testify/require"
)

func TestScraper(t *testing.T) {
	factory := rabbitmqreceiver.NewFactory()
	out, err := ScrapeReceiverConfig(factory)
	require.NoError(t, err)
	require.NotNil(t, out)
	fmt.Println(*out)
}
