/* Package dataprocessor deals with processing data received from tasks and converts them to correct
output */
package dataprocessor

import (
	"encoding/json"

	"github.com/maplelabs/github-audit/logger"
)

var (
	log logger.Logger
)

func init() {
	log = logger.GetLogger()
}

// DataProcessor provides methods for process various data based on documents
type DataProcessor interface {
	// ProcessCommits process commit documents , takes data in bytes and tags as input
	ProcessCommits([]byte, map[string]string) ([]interface{}, error)

	// ProcessPullRequests process pull request documents , takes data in bytes and tags as input
	ProcessPullRequests([]byte, map[string]string) ([]interface{}, error)

	// ProcessIssues process issue documents , takes data in bytes and tags as input
	ProcessIssues([]byte, map[string]string) ([]interface{}, error)
}

// NewDataProcessor returns a new data processor based on host type
func NewDataProcessor(host string, repoName string, repoURL string) DataProcessor {
	if host == "github" {
		return NewGithubProcessor(repoName, repoURL)
	}
	return nil
}

// AddTags adds tags to data which were passed in config.yaml
func AddTags(data []byte, tags map[string]string) []interface{} {
	var docMap []map[string]interface{}
	var finalDocs []interface{}
	err := json.Unmarshal(data, &docMap)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling tags", err)
		return finalDocs
	}

	for i := range docMap {
		for k, v := range tags {
			docMap[i][k] = v
		}
		finalDocs = append(finalDocs, docMap[i])
	}
	return finalDocs
}
