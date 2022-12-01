/*
Package publisher provides access to various target clients and APIs to push data.
*/
package publisher

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/maplelabs/github-audit/logger"

	retryhttp "github.com/hashicorp/go-retryablehttp"
)

var (
	log               logger.Logger
	ErrUnknownPubType = errors.New("unknown publisher type")
)

func init() {
	log = logger.GetLogger()
}

//Add various client constant here
const (
	KAFKAREST     = "kafka-rest"
	ELASTICSEARCH = "elasticsearch"
)

// Publisher is implemented by any client that has Publish method.
type Publisher interface {
	Publish([]interface{}) error
}

// NewPublisher returns an instance of publisher for pushing data based on publisher type
func NewPublisher(pubType string, config map[string]string) (Publisher, error) {
	switch strings.ToLower(pubType) {
	case KAFKAREST:
		return loadConfigKafkaRest(config)
	case ELASTICSEARCH:
		return loadConfigElasticSearch(config)
	default:
		return nil, ErrUnknownPubType
	}
}

// loadConfigKafkaRest loads config to KafkaRest and return pointer to KafkaRest
func loadConfigKafkaRest(config map[string]string) (*KafkaRestClient, error) {
	var (
		kafkaRest KafkaRestClient
		err       error
	)
	cfgByte, err := json.Marshal(config)
	if err != nil {
		log.Errorf("error[%v] in marshalling kafka config", err)
		return &kafkaRest, err
	}
	err = json.Unmarshal(cfgByte, &kafkaRest)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling kafka config to kafkarest struct", err)
		return &kafkaRest, err
	}

	return &kafkaRest, nil
}

// loadConfigElasticSearch loads config to ElasticSearch and return pointer to ElasticSearch
func loadConfigElasticSearch(config map[string]string) (*ElasticSearchClient, error) {
	var (
		es  ElasticSearchClient
		err error
	)
	cfgByte, err := json.Marshal(config)
	if err != nil {
		log.Errorf("error[%v] in marshalling elasticsearch config", err)
		return &es, err
	}
	err = json.Unmarshal(cfgByte, &es)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling elasticsearch config", err)
		return &es, err
	}
	return &es, nil
}

// HTTPClientWithRetry creates a HTTP client
func HTTPClientWithRetry() *retryhttp.Client {
	client := retryhttp.NewClient()
	client.RetryWaitMin = 500 * time.Millisecond
	client.RetryMax = 1
	client.Logger = log
	return client
}
