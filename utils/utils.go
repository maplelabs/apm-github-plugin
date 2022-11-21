/* Package utils contains utility functions common to all packages */
package utils

import (
	"encoding/base64"

	"github.com/maplelabs/github-audit/logger"
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
