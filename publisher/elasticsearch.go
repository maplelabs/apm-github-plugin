package publisher

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"strconv"

	retryhttp "github.com/hashicorp/go-retryablehttp"
)

// ElasticSearchClient holds config for Kafkarest target
type ElasticSearchClient struct {
	// Host for elasticsearch
	Host string `yaml:"host" json:"host"`

	// Port for elasticsearch
	Port string `yaml:"port" json:"port"`

	// Protocol to connect to elasticsearch
	Protocol string `yaml:"protocol" json:"protocol"`

	// Index or profile name to send data 
	Index string `yaml:"index" json:"index"`

	// Username if basic auth is used
	Username string `yaml:"username,omitempty" json:"username,omitempty"`

	// Password if basic auth is used
	Password string `yaml:"password,omitempty" json:"password,omitempty"`

	// ES_7x to know if target is before elastic version 7.x
	ES_7x string `yaml:"old_es" json:"old_es"`

	// Path to send github data to target project
	Path string `yaml:"path" json:"path"`

}

// BulkResponse ...
type BulkResponse struct {
	Took   int                `json:"took"`
	Errors bool               `json:"errors"`
	Items  []map[string]Index `json:"items"`
}

// BulkError ...
type BulkError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// Shards ...
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// Index ...
type Index struct {
	Index       string    `json:"_index"`
	Type        string    `json:"_type"`
	ID          string    `json:"_id"`
	Version     int       `json:"_version"`
	Result      string    `json:"result"`
	Shards      Shards    `json:"_shards"`
	SeqNo       int       `json:"_seq_no"`
	PrimaryTerm int       `json:"_primary_term"`
	Status      int       `json:"status"`
	Error       BulkError `json:"error"`
}

// ErrorResponse es error response struct
type ErrorResponse struct {
	Error  Error `json:"error"`
	Status int   `json:"status"`
}

// RootCause ...
type RootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// CausedBy ...
type CausedBy struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// Error ...
type Error struct {
	RootCause []RootCause `json:"root_cause"`
	Type      string      `json:"type"`
	Reason    string      `json:"reason"`
	CausedBy  CausedBy    `json:"caused_by"`
}

// Publish pushes the data to target
func (es *ElasticSearchClient) Publish(data []interface{}) error {
	log.Debugf("In Publish: data: %+v", data)
	log.Infof("In Publish: es: %+v", es)
	if len(data) == 0 {
		return nil
	}
	client := HTTPClientWithRetry()
	var (
		bulkdata []byte
		err      error
	)
	const bulkStart = "{\"index\":{}}\n"
	const bulkEnd = "\n"

	for _, doc := range data {
		byteData, err := json.Marshal(doc)
		if err != nil {
			log.Errorf("error[%v] unable to marshal data", err)
		} else {
			bulkdata = append(bulkdata, bulkStart...)
			bulkdata = append(bulkdata, byteData...)
			bulkdata = append(bulkdata, bulkEnd...)
		}
	}

	bulkdata = append(bulkdata, bulkEnd...)
    log.Debugf("Publish bulkdata: %v", len(bulkdata))
	reqURL := fmt.Sprintf("%s_bulk/", es.getURL())
	log.Infof("reqURL: %v", reqURL)
	request, err := retryhttp.NewRequest(
		http.MethodPost,
		reqURL,
		bytes.NewReader(bulkdata),
	)
	if err != nil {
		log.Errorf("error[%v] unable to create new request", err)
		return err
	}
	log.Debugf("request: %+v, body: %+v, ", request.Request, request.Body)

	request.Header.Set("Content-Type", "application/json")
	if es.Username != "" && es.Password != "" {
		request.SetBasicAuth(es.Username, es.Password)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqWithTimeout := request.WithContext(timeout)
	if es.Protocol == "https" {
		client.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		}
	}

	response, err := client.Do(reqWithTimeout)
	if err != nil {
		log.Errorf("error[%v] request failed with error", err)
		return err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("error[%v] failed to read response body", err)
	}
	log.Debugf("respBody after writing to elastisearch\n%s", string(respBody))

	var bulkResp BulkResponse
	err = json.Unmarshal(respBody, &bulkResp)
	if err != nil {
		log.Errorf("error[%v] failed to Unmarshal response body", err)
		return errors.New(response.Status)
	}
	log.Debugf("error[%v] bulkResp", bulkResp)

	bulkrespErr := []error{}
	if bulkResp.Errors {
		for _, item := range bulkResp.Items {
			idoc := item["index"]
			if idoc.Status > 299 {
				bulkrespErr = append(
					bulkrespErr,
					fmt.Errorf("doc error: %s", idoc.Error.Reason),
				)
			}
		}
		log.Errorf("error[%v] response errors", bulkrespErr)
		return nil
	}

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated {
		log.Infof("successfully sent doc to %s - %d", es.getURL(), response.StatusCode)
		return nil
	}
	log.Errorf("[elasticsearch] failed to send doc to %s - %d", es.getURL(), response.StatusCode)
	var esError ErrorResponse
	err = json.Unmarshal(respBody, &esError)
	if err != nil {
		log.Errorf("error[%v] failed to Unmarshal response body, %s", err)
		return errors.New(response.Status)
	}
	log.Errorf("error[%v]", esError)
	return fmt.Errorf("status code: %d, error: %s", esError.Status, esError.Error.Reason)
}

// getURL returns url to elasticendpoint for writing data
func (es *ElasticSearchClient) getURL() string {
	var url, esType string
	log.Debugf("getURL es.ES_7x: %v", es.ES_7x)
	//es.ES_7x = es.getESVersion()
	boolFlag, _ := strconv.ParseBool(es.ES_7x)
	if boolFlag {
		esType = "doc"
	} else {
		esType = "_doc"
	}
	// create target path to push data, eg, metric-<profile>-<project>-$_write
	profilePath := "metric-" + es.Index + "-" + es.Path + "-$_write/" + esType
	url = fmt.Sprintf("%s://%s:%s/%s/", es.Protocol, es.Host, es.Port, profilePath)
	log.Infof("ES getURL url: %v", url)
	return url
}

// getESVersion returns if Es version is greater than 7.x or not
func (es *ElasticSearchClient) getESVersion() bool {
	type EsResponse struct {
		Name        string `json:"name"`
		ClusterName string `json:"cluster_name"`
		ClusterUUID string `json:"cluster_uuid"`
		Version     struct {
			Number                           string    `json:"number"`
			BuildFlavor                      string    `json:"build_flavor"`
			BuildType                        string    `json:"build_type"`
			BuildHash                        string    `json:"build_hash"`
			BuildDate                        time.Time `json:"build_date"`
			BuildSnapshot                    bool      `json:"build_snapshot"`
			Distribution                     string    `json:"distribution"`
			LuceneVersion                    string    `json:"lucene_version"`
			MinimumWireCompatibilityVersion  string    `json:"minimum_wire_compatibility_version"`
			MinimumIndexCompatibilityVersion string    `json:"minimum_index_compatibility_version"`
		} `json:"version"`
		Tagline string `json:"tagline"`
	}
	var url string
	switch es.Protocol {
	case "http":
		url = fmt.Sprintf("http://%s:%s/", es.Host, es.Port)
	case "https":
		url = fmt.Sprintf("https://%s:%s/", es.Host, es.Port)
	}
	log.Infof("getESVersion url: %v", url)
	client := HTTPClientWithRetry()
	request, err := retryhttp.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	request.Header.Set("Content-Type", "application/json")
	if es.Username != "" && es.Password != "" {
		request.SetBasicAuth(es.Username, es.Password)
	}
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqWithTimeout := request.WithContext(timeout)
	if es.Protocol == "https" {
		client.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		}
	}

	resp, err := client.Do(reqWithTimeout)
	log.Infof("getESVersion resp: %v", resp)
	if err != nil {
		log.Errorf("error[%v] request to get es version failed with error", err)
		return false
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error[%v] failed to read response body", err)
		return false
	}
	log.Debug("response to version check ", string(respBody))
	var esResp EsResponse
	err = json.Unmarshal(respBody, &esResp)
	if err != nil {
		log.Errorf("error[%v] failed to Unmarshal response body", err)
		return false
	}
	version := esResp.Version.Number
	dist := esResp.Version.Distribution
	// if opensearch is used it is es ver7 onwards always
	if strings.HasPrefix(version, "7") || strings.ToLower(dist) == "opensearch" {
		return true
	} else {
		return false
	}
}
