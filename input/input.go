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
	DefaultMasterBranch = "master"
	DefaultMainBranch   = "main"
	PRIVATE             = "private"
)

var (
	ErrMissingAuditJob        = errors.New("no audit job defined")
	ErrMissingAuditJobName    = errors.New("missing audit job name")
	ErrMissingTarget          = errors.New("no target defined")
	ErrMissingAccessToken     = errors.New("missing access token")
	ErrMissingUsername        = errors.New("missing username")
	ErrMissingRepositoryName  = errors.New("missing repository name")
	ErrMissingPollingInterval = errors.New("missing polling interval")
	ErrMissingRepositoryHost  = errors.New("missing repository host")
	ErrMissingRepositoryOwner = errors.New("missing repository owner")
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
	// Loglevel defines the logging level for the log file.
	Loglevel string `yaml:"loglevel,omitempty" json:"loglevel,omitempty"`

	// Logpath represents the path where log file will be generated.
	Logpath string `yaml:"logpath,omitempty" json:"logpath,omitempty"`

	// Auditjobs is list of auditjobs.
	AuditJobs []AuditJob `yaml:"auditJobs" json:"auditJobs"`

	// Targets is list of targets.
	Targets []Target `yaml:"targets" json:"targets"`
}

// AuditJob represents the auditing job for which data will be fetched from remote git repository.
type AuditJob struct {
	// Name of the auditjob.
	Name string `yaml:"name" json:"name"`

	// Polling interval of the audit job.  Format: 10s , 10m , 10h , 1d
	PollingInterval string `yaml:"polling_interval" json:"polling_interval"`

	// Metadata if any
	Metadata interface{} `yaml:"metadata" json:"metadata"`

	// Tags if any
	Tags `yaml:"tags,omitempty" json:"tags"`

	// Output targets.
	Output `yaml:"output" json:"output"`

	// Repository type , options include github , gitlabs etc.
	RepositoryHost string `yaml:"repo_host" json:"repo_host"`

	// RepositoryName represents the name of the repository.
	RepositoryName string `yaml:"repo_name" json:"repo_name"`

	// RepositoryName represents the name of the repository.
	RepositoryOwner string `yaml:"repo_owner" json:"repo_owner"`

	// RepositoryConfig defines repository config.
	RepositoryConfig `yaml:"repo_config" json:"repo_config"`
}

// RepositoryConfig represents repostory configurations.
type RepositoryConfig struct {
	// RepositoryType is whether public or private repo.
	RepositoryType string

	// RepositoryURL is the url for the repository.
	RepositoryURL string `yaml:"repo_url" json:"repo_url"`

	// RepositoryCredentials represents repository credentials.
	RepositoryCredentials `yaml:"credentials" json:"credentials"`

	//Branches represents branches that needs to be monitored.
	Branches []string `yaml:"branches" json:"branches"`
}

// RepositoryCredentials consists of credential needed to authenticate with repository.
type RepositoryCredentials struct {
	//Username of the repository.
	Username string `yaml:"username" json:"username"`

	//AccessToken for accessing the github's APIs.
	AccessToken string `yaml:"access_token" json:"access_token"`
}

// Output represents the target where data will be sent.
type Output struct {
	//TargetName consists of the target names to which auditjob data needs to be sent.
	TargetName []string `yaml:"target_name" json:"target_name"`
}

// Tags are additional information that can be added in tags fields of output documents.
type Tags map[string]string

// Target represents the target where data will be published.
type Target struct {
	//Name of the target.
	Name string `yaml:"name" json:"name"`

	//Type of the target , possible values (elasticsearch , kafka-rest etc).
	Type string `yaml:"type" json:"type"`

	//TargetConfig consist of the target configurations.
	TargetConfig map[string]string `yaml:"config" json:"config"`
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
		if j.AccessToken == "" && j.RepositoryType == PRIVATE {
			return ErrMissingAccessToken
		}
		// checking if either username or email is not empty.
		if j.Username == "" {
			return ErrMissingUsername
		}
		// checking if repository name is not empty.
		if j.RepositoryName == "" {
			return ErrMissingRepositoryName
		}
		// checking if repository host is not empty.
		if j.RepositoryHost == "" {
			return ErrMissingRepositoryHost
		}
		// checking if repository owner is not empty.
		if j.RepositoryOwner == "" {
			return ErrMissingRepositoryOwner
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
			log.Debugf("branch is not configurd for auditjob %v adding default branch", c.AuditJobs[i].Name)
			c.AuditJobs[i].Branches = append(c.AuditJobs[i].Branches, DefaultMasterBranch, DefaultMainBranch)
		}
		//populating access token from environment variable if present. environment variable name is same as auditjob name.
		accessTokenFromEnv := os.Getenv(c.AuditJobs[i].Name)
		if accessTokenFromEnv != "" {
			log.Debugf("adding access token from environment variable for auditjob %v", c.AuditJobs[i].Name)
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
	// adding default values for missing fields
	log.Infof("populating default field values in github-audit config")
	conf.populateDefaultValues()
	// validating config
	log.Infof("validating fields in config file located at %s", configPath)
	err = conf.validate()
	if err != nil {
		log.Errorf("error[%v] in validating config file", err)
		return Config{}, err
	}
	return conf, nil
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
