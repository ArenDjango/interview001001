package contracts

import "context"

type ConfigurationServiceInterface interface {
	GetConfigurationForApplication(ctx context.Context) (*ConfigurationForApplication, error)
}

type ConfigurationForApplication struct {
	PollingClientsFromQortIntervalInSeconds uint64
	ClientsPollingBatchSize                 uint64
}

const (
	PollingClientsFromQortIntervalInSeconds = "PollingClientsFromQortIntervalInSeconds"
	ClientsPollingBatchSize                 = "QuantityOfClientsPollingEachPolling"
)

var DefaultValues = map[string]string{
	PollingClientsFromQortIntervalInSeconds: "120",
	ClientsPollingBatchSize:                 "1000",
}
