package publisher

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	retryhttp "github.com/hashicorp/go-retryablehttp"
)

// KafkaRestClient holds config for Kafkarest target
type KafkaRestClient struct {
	// Host for kafka rest
	Host string `yaml:"host" json:"host"`

	// Port for kafka rest
	Port string `yaml:"port" json:"port"`

	// Protocol to connect
	Protocol string `yaml:"protocol" json:"protocol"`

	// Topic name to send data to
	Topic string `yaml:"topic" json:"topic"`

	// Username if basic auth is used
	Username string `yaml:"user,omitempty" json:"user,omitempty"`

	// Password if basic auth is used
	Password string `yaml:"password,omitempty" json:"password,omitempty"`

	// Token for authentication
	Token string `json:"token"`

	// Path holds endpoint url
	Path string `json:"path"`
}

// BulkResponse ...
type KafkaBulkResponse struct {
	KafkaErrorResponse
	Offsets       []Offset `json:"offsets"`
	KeySchemaID   int      `json:"key_schema_id"`
	ValueSchemaID int      `json:"value_schema_id"`
}

// Offset ...
type Offset struct {
	Partition int         `json:"partition"`
	Offset    int         `json:"offset"`
	ErrorCode interface{} `json:"error_code"`
	Error     interface{} `json:"error"`
}

// ErrorResponse es error response struct
type KafkaErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"error_code"`
}

// Publish pushes the data to target
func (kc *KafkaRestClient) Publish(data []interface{}) error {
	var (
		bulkdata []byte
		err      error
	)
	if len(data) == 0 {
		return nil
	}
	client := HTTPClientWithRetry()
	const recordStart = `{"records":[`
	const recordEnd = `]}`
	const dataStart = `{"value":`
	const dataEnd = "}"

	bulkdata = append(bulkdata, recordStart...)

	for i, doc := range data {
		byteData, err := json.Marshal(doc)
		if err != nil {
			log.Errorf("error[%v] unable to marshal data: ")
		} else {
			bulkdata = append(bulkdata, dataStart...)
			bulkdata = append(bulkdata, byteData...)
			bulkdata = append(bulkdata, dataEnd...)
			if i != len(data)-1 {
				bulkdata = append(bulkdata, []byte(",")...)
			}
		}
	}

	bulkdata = append(bulkdata, recordEnd...)

	reqURL := kc.getURL(kc.Topic)
	request, err := retryhttp.NewRequest(
		http.MethodPost,
		reqURL,
		bytes.NewReader(bulkdata),
	)

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/vnd.kafka.json.v2+json")
	request.Header.Set("Accept", "application/vnd.kafka.v2+json")

	if kc.Token != "" {
		request.Header.Set("Authorization", kc.Token)
	}

	if kc.Username != "" && kc.Password != "" {
		request.SetBasicAuth(kc.Username, kc.Password)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	reqWithTimeout := request.WithContext(timeout)
	if kc.Protocol == "https" {
		client.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		}
	}

	response, err := client.Do(reqWithTimeout)
	if err != nil {
		log.Errorf("error[%v] request failed", err)
		return err
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("error[%v] failed to read response body,")
	}

	var bulkResp KafkaBulkResponse
	err = json.Unmarshal(respBody, &bulkResp)
	if err != nil {
		log.Errorf("error[%v] failed to Unmarshal response body")
		return errors.New(response.Status)
	}

	bulkrespErr := []error{}

	if bulkResp.KafkaErrorResponse.Status != 0 {
		bulkrespErr = append(
			bulkrespErr,
			fmt.Errorf("Doc error: %s", bulkResp.KafkaErrorResponse.Message),
		)
		log.Errorf("error[%v] kafkarest bulk error", bulkrespErr)
		return nil
	}

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated {
		log.Infof("successfully sent doc to %s", kc.getURL(kc.Topic))
		return nil
	}

	return fmt.Errorf("error in sending data to kafka rest ")
}

func (kc *KafkaRestClient) getURL(topic string) (url string) {
	if kc.Path == "" {
		return fmt.Sprintf("%s://%s:%s/topics/%s", kc.Protocol, kc.Host, kc.Port, topic)
	}
	return fmt.Sprintf("%s://%s:%s/%s/topics/%s", kc.Protocol, kc.Host, kc.Port, kc.Path, topic)
}
