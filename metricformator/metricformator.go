/* MetricFormator provides functionality to customise output documents */
package metricformator

import (
	"encoding/json"
	"os"

	"github.com/maplelabs/github-audit/logger"
)

var (
	// logging utility
	log logger.Logger

	// file storing metric formatting related information
	metricFormator = "metricformator/metricformator.json"
)

func init() {
	log = logger.GetLogger()
}

// MetricFormator represents metric formator for customisation
type MetricFormator struct {
	ChangeDefaultKeys map[string]string `json:"changeDefaultKeys"`
	AddNewGlobalKeys  map[string]string `json:"addNewGlobalKeys"`
}

// NewMetricFormator return a new metric formator
func NewMetricFormator() *MetricFormator {
	mf := new(MetricFormator)
	fileByte, err := os.ReadFile(metricFormator)
	if err != nil {
		log.Errorf("error[%v] in reading metricFormator file", err)
	}
	if len(fileByte) > 0 {
		err = json.Unmarshal(fileByte, mf)
		if err != nil {
			log.Errorf("error[%v] in unmarshalling metricFormator file", err)
		}
	}
	return mf
}

// CustomizeMetrics modifies output documents based on customisation
func (m *MetricFormator) CustomizeMetrics(data []byte) []byte {
	var docMap []map[string]interface{}
	err := json.Unmarshal(data, &docMap)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling data for customisation", err)
		return data
	}
	for i := range docMap {
		for k, v := range m.ChangeDefaultKeys {
			if val, ok := docMap[i][k]; ok {
				docMap[i][v] = val
				delete(docMap[i], k)
			}
		}
		for k, v := range m.AddNewGlobalKeys {
			docMap[i][k] = v
		}
	}
	byteData, err := json.Marshal(docMap)
	if err != nil {
		log.Errorf("error[%v] in marshalling data for customisation", err)
		return data
	}
	return byteData
}
