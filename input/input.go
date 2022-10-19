/*
Package Input reads config file and populates the data in Config struct.
*/
package input

import (
	"errors"
	"os"
	"strconv"

	"github.com/maplelabs/github-audit/logger"
	"gopkg.in/yaml.v3"
)

const (
	DefaultBranch = "master"
)

var (
	ErrMissingAuditJob        = errors.New("no audit job defined")
	ErrMissingAuditJobName    = errors.New("missing audit job name")
	ErrMissingTarget          = errors.New("no target defined")
	ErrMissingAccessToken     = errors.New("missing access token")
	ErrMissingUsernameEmail   = errors.New("missing username or email")
	ErrMissingRepositoryName  = errors.New("missing repository name")
	ErrMissingRepositoryURL   = errors.New("missing repository URL")
	ErrMissingPollingInterval = errors.New("missing polling interval")
	ErrMissingRepositoryType  = errors.New("missing repository type")
	ErrMissingTargetNameList  = errors.New("missing target name in audit job")
	ErrMissingTargetName      = errors.New("missing target name")
	ErrMissingTargetType      = errors.New("missing target type")
	ErrPollingIntervalFormat  = errors.New("polling interval format is incorrect")
)

var (
	log logger.Logger
)

func init() {
	log = logger.GetLogger()
}

// Config represents the data parsed from config.yaml.
type Config struct {
	// Loglevel defines the logging level for the log file. <OPTIONAL>
	Loglevel string `yaml:"loglevel,omitempty" json:"loglevel,omitempty"`
	// Logpath represents the path where log file will be generated. <OPTIONAL>
	Logpath string `yaml:"logpath,omitempty" json:"logpath,omitempty"`
	// Auditjobs is list of auditjobs. <REQUIRED>
	AuditJobs []AuditJob `yaml:"auditJobs" json:"auditJobs"`
	// Targets is list of targets. <REQUIRED>
	Targets []Target `yaml:"targets" json:"targets"`
}

// validate checks for mandatory fields in config struct and returns error in case of something is missing.
func (c *Config) validate() error {
	// checking if at least one job is present.
	if len(c.AuditJobs) == 0 {
		return ErrMissingAuditJob
	}
	// checking if at least one target is present.
	if len(c.Targets) == 0 {
		return ErrMissingTarget
	}
	for _, j := range c.AuditJobs {
		// checking if audit job name is not empty.
		if j.Name == "" {
			return ErrMissingAuditJobName
		}
		// checking if access token is not empty.
		if j.AccessToken == "" {
			return ErrMissingAccessToken
		}
		// checking if either username or email is not empty.
		if j.Username == "" && j.Email == "" {
			return ErrMissingUsernameEmail
		}
		// checking if repository name is not empty.
		if j.RepositoryName == "" {
			return ErrMissingRepositoryName
		}
		// checking if repository url is not empty.
		if j.RepositoryURL == "" {
			return ErrMissingRepositoryURL
		}
		// checking if repository type is not empty.
		if j.RepositoryType == "" {
			return ErrMissingRepositoryType
		}
		// checking if any target is defined in output.
		if len(j.TargetName) == 0 {
			return ErrMissingTargetName
		}
		// checking if polling interval is not empty.
		if j.PollingInterval == "" {
			return ErrMissingPollingInterval
		}
		if j.PollingInterval != "" {
			return checkPollingIntervalFormat(j.PollingInterval)
		}
	}
	for _, t := range c.Targets {
		// checking if target name is missing.
		if t.Name == "" {
			return ErrMissingTargetName
		}
		// checking if target type is missing.
		if t.Type == "" {
			return ErrMissingTargetType
		}
	}
	return nil
}

// populateDefaultValues puts default values to optional fields in config.
func (c *Config) populateDefaultValues() {
	for i := range c.AuditJobs {
		if len(c.AuditJobs[i].Branches) == 0 {
			c.AuditJobs[i].Branches = append(c.AuditJobs[i].Branches, DefaultBranch)
		}
		//populating access token from environment variable if present. environment variable name is same as auditjob name.
		accessTokenFromEnv := os.Getenv(c.AuditJobs[i].Name)
		if accessTokenFromEnv != "" {
			c.AuditJobs[i].AccessToken = accessTokenFromEnv
		}
	}
}

/*
checkPollingIntervalFormat checks the format for polling interval.
Format: Integer[s/m/h/d]  (s:seconds , m:minutes , h:hours , d:days).
Example: Accepted formats are 10s, 10m, 10h, 1d
*/
func checkPollingIntervalFormat(pi string) error {
	if pi == "" {
		return ErrMissingPollingInterval
	}
	lastChar := string(pi[len(pi)-1])
	if lastChar == "h" || lastChar == "s" || lastChar == "m" || lastChar == "d" {
		remChar := string(pi[0 : len(pi)-1])
		if !isNumeric(remChar) {
			return ErrPollingIntervalFormat
		}
	} else {
		return ErrPollingIntervalFormat
	}
	return nil
}

//isNumeric checks if string is numeric or not
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// AuditJob represents the auditing job for which data will be fetched from remote git repository.
type AuditJob struct {
	// Name of the auditjob. <REQUIRED>
	Name string `yaml:"name" json:"name"`
	// Polling interval of the audit job. <REQUIRED> Format: 10s , 10m , 10h , 1d
	PollingInterval string `yaml:"polling_interval" json:"polling_interval"`
	// Metadata if any <OPTIONAL>
	Metadata interface{} `yaml:"metadata" json:"metadata"`
	// Tags if any <OPTIONAL>
	Tags `yaml:"tags,omitempty" json:"tags"`
	// Output targets. <REQUIRED>
	Output `yaml:"output" json:"output"`
	// Repository type , options include github , gitlabs etc. <REQUIRED>
	RepositoryType string `yaml:"repo_type" json:"repo_type"`
	// RepositoryName represents the name of the repository. <REQUIRED>
	RepositoryName string `yaml:"repo_name" json:"repo_name"`
	// RepositoryConfig defines reposiroty config. <REQUIRED>
	RepositoryConfig `yaml:"repo_config" json:"repo_config"`
}

// RepositoryConfig represents repostory configurations.
type RepositoryConfig struct {
	// RepositoryURL is the url for the repository. <REQUIRED>
	RepositoryURL string `yaml:"repo_url" json:"repo_url"`
	// RepositoryCredentials represents repository credentials. <REQUIRED>
	RepositoryCredentials `yaml:"credentials" json:"credentials"`
	//Branches represents branches that needs to be monitored. <OPTIONAL>
	Branches []string `yaml:"branches" json:"branches"`
}

// RepositoryCredentials consists of credential needed to authenticate with repository.
type RepositoryCredentials struct {
	//Either username or email is <REQUIRED>
	//Username of the repository.
	Username string `yaml:"username" json:"username"`
	//Email of the repository.
	Email string `yaml:"email" json:"email"`
	//AccessToken for accessing the github's APIs. <REQUIRED>
	AccessToken string `yaml:"access_token" json:"access_token"`
}

// Output represents the target where data will be sent.
type Output struct {
	//TargetName consists of the target names to which auditjob data needs to be sent. <REQUIRED>
	TargetName []string `yaml:"target_name" json:"target_name"`
}

// Tags are additional information that can be added in tags fields of output documents.
type Tags map[string]string

// Target represents the target where data will be published.
type Target struct {
	//Name of the target. <REQUIRED>
	Name string `yaml:"name" json:"name"`
	//Type of the target , possible values (elasticsearch , kafka-rest etc). <REQUIRED>
	Type string `yaml:"type" json:"type"`
	//TargetConfig consist of the target configurations.
	TargetConfig `yaml:"config" json:"config"`
}

//TargetConfig is key value target configurations.
type TargetConfig map[string]string

// InitConfig reads the config file and populates Config struct with relevant data. It returns error in case of file read error based on ReadConfigFile().
func InitConfig(configPath string) (Config, error) {
	log.Infof("reading config file located at %s", configPath)
	var conf Config
	data, err := ReadConfigFile(configPath)
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Errorf("error[%v] in unmarshaling config file to Config struct", err)
		return Config{}, err
	}
	// validating config
	log.Infof("validating fields in config file located at %s", configPath)
	err = conf.validate()
	if err != nil {
		log.Errorf("error[%v] in validating config file", err)
		return Config{}, err
	}
	// adding default values for missing fields
	log.Infof("populating default field values in github-audit config")
	conf.populateDefaultValues()
	return conf, err
}

// ReadConfigFile reads the config file and returns the content. A successful call returns err == nil. It returns error based on underline os.Readfile().
func ReadConfigFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("error[%v] in reading config file located at %s", filePath, err)
		return nil, err
	}
	return data, nil
}
