/* Package utils contains utility functions common to all packages */
package utils

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/maplelabs/github-audit/logger"
)

const (
	// Default Scheduling interval for each audit job.
	DEFAULTDURATION = time.Duration(5 * time.Minute)
)

var (
	log logger.Logger
)

func init() {
	log = logger.GetLogger()
}

// DecodeAccessKey decodes the base64 encoded access key as provided in config.yaml or through environment variable.
func DecodeAccessKey(key string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Errorf("error[%v] in decoding accessToken", err)
		return "", err
	}
	return string(decodedKey), nil
}

// ConvertIntervalToDuration converts the scheduling interval as provided in config.yaml to golang's duration
func ConvertIntervalToDuration(interval string) time.Duration {
	interval = strings.ToLower(interval)
	duration, err := time.ParseDuration(interval)
	if err != nil {
		log.Errorf("error[%v] in converting scheduling duration to golang's duration , setting default duration of 5 minutes", err)
		return DEFAULTDURATION
	}
	return duration
}
