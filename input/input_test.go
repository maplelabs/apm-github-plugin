/*
Package Input reads config file and populates the data in Config struct.
*/
package input

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig_validate(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantErr bool
	}{
		{
			name: "correct input with all validated fields",
			c: &Config{
				Loglevel: "info",
				Logpath:  "./test.yaml",
				AuditJobs: []AuditJob{
					{
						Name:            "auditjob1",
						PollingInterval: "30s",
						Output: Output{
							TargetName: []string{"testtarget1"},
						},
						RepositoryType: "github",
						RepositoryName: "testRepo",
						RepositoryConfig: RepositoryConfig{
							RepositoryURL: "https://testurl",
							RepositoryCredentials: RepositoryCredentials{
								Username:    "testUSer",
								AccessToken: "1234adsr",
							},
						},
					},
				},
				Targets: []Target{
					{
						Name: "testtarget1",
						Type: "elasticsearch",
						TargetConfig: TargetConfig{
							"host":     "test",
							"protocol": "http",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "incorrect input with empty auditjobs",
			c: &Config{
				Loglevel:  "info",
				Logpath:   "./test.yaml",
				AuditJobs: []AuditJob{},
				Targets: []Target{
					{
						Name: "testtarget1",
						Type: "elasticsearch",
						TargetConfig: TargetConfig{
							"host":     "test",
							"protocol": "http",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "incorrect input with empty targets",
			c: &Config{
				Loglevel: "info",
				Logpath:  "./test.yaml",
				AuditJobs: []AuditJob{
					{
						Name:            "auditjob1",
						PollingInterval: "30s",
						Output: Output{
							TargetName: []string{"testtarget1"},
						},
						RepositoryType: "github",
						RepositoryName: "testRepo",
						RepositoryConfig: RepositoryConfig{
							RepositoryURL: "https://testurl",
							RepositoryCredentials: RepositoryCredentials{
								Username:    "testUSer",
								AccessToken: "1234adsr",
							},
						},
					},
				},
				Targets: []Target{},
			},
			wantErr: true,
		},
		{
			name: "incorrect input with some missing fields",
			c: &Config{
				Loglevel: "info",
				Logpath:  "./test.yaml",
				AuditJobs: []AuditJob{
					{
						Name:            "",
						PollingInterval: "30s",
						Output: Output{
							TargetName: []string{"testtarget1"},
						},
						RepositoryType: "github",
						RepositoryName: "",
						RepositoryConfig: RepositoryConfig{
							RepositoryURL: "https://testurl",
							RepositoryCredentials: RepositoryCredentials{
								Username:    "testUSer",
								AccessToken: "",
							},
						},
					},
				},
				Targets: []Target{
					{
						Name: "testtarget1",
						Type: "elasticsearch",
						TargetConfig: TargetConfig{
							"host":     "test",
							"protocol": "http",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkPollingIntervalFormat(t *testing.T) {
	type args struct {
		pi string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "checking interval ending with s seconds correctly",
			args:    args{pi: "10s"},
			wantErr: false,
		},
		{
			name:    "checking interval ending with m minutes correctly",
			args:    args{pi: "10m"},
			wantErr: false,
		},
		{
			name:    "checking interval ending with h hours correctly",
			args:    args{pi: "10h"},
			wantErr: false,
		},
		{
			name:    "checking interval ending with d days correctly",
			args:    args{pi: "10d"},
			wantErr: false,
		},
		{
			name:    "checking interval ending with other than s,m,d,h",
			args:    args{pi: "10a"},
			wantErr: true,
		},
		{
			name:    "checking interval starting with non numeric",
			args:    args{pi: "a10a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPollingIntervalFormat(tt.args.pi); (err != nil) != tt.wantErr {
				t.Errorf("checkPollingIntervalFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isNumeric(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correct numeric string",
			args: args{s: "10"},
			want: true,
		},
		{
			name: "incorrect numeric string",
			args: args{s: "1012as"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumeric(tt.args.s); got != tt.want {
				t.Errorf("isNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	var correct Config = Config{
		Loglevel: "info",
		Logpath:  "./test.yaml",
		AuditJobs: []AuditJob{
			{
				Name:            "auditjob1",
				PollingInterval: "30s",
				Output: Output{
					TargetName: []string{"testtarget1"},
				},
				RepositoryType: "github",
				RepositoryName: "testRepo",
				RepositoryConfig: RepositoryConfig{
					RepositoryURL: "https://testurl",
					RepositoryCredentials: RepositoryCredentials{
						Username:    "testUSer",
						AccessToken: "1234adsr",
					},
					Branches: []string{"master"},
				},
			},
		},
		Targets: []Target{
			{
				Name: "testtarget1",
				Type: "elasticsearch",
				TargetConfig: TargetConfig{
					"host":     "test",
					"protocol": "http",
				},
			},
		},
	}
	var incorrect Config = Config{
		Loglevel: "info",
		Logpath:  "./test.yaml",
		AuditJobs: []AuditJob{
			{
				Name:            "",
				PollingInterval: "30s",
				Output: Output{
					TargetName: []string{"testtarget1"},
				},
				RepositoryType: "github",
				RepositoryName: "",
				RepositoryConfig: RepositoryConfig{
					RepositoryURL: "",
					RepositoryCredentials: RepositoryCredentials{
						Username:    "testUSer",
						AccessToken: "1234adsr",
					},
				},
			},
		},
		Targets: []Target{
			{
				Name: "testtarget1",
				Type: "elasticsearch",
				TargetConfig: TargetConfig{
					"host":     "test",
					"protocol": "http",
				},
			},
		},
	}
	//creating correctconfigtest.yaml file with dummy content
	correctBytes, _ := yaml.Marshal(correct)
	_ = os.WriteFile("correctconfigtest.yaml", correctBytes, 0777)
	//creating incorrectconfigtest.yaml file with dummy content
	incorrectBytes, _ := yaml.Marshal(incorrect)
	_ = os.WriteFile("incorrectconfigtest.yaml", incorrectBytes, 0777)
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name:    "correct config.yaml",
			args:    args{configPath: "correctconfigtest.yaml"},
			want:    correct,
			wantErr: false,
		},
		{
			name:    "incorrect config.yaml",
			args:    args{configPath: "incorrectconfigtest.yaml"},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitConfig(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
	//deleting correctconfigtest.yaml file
	os.Remove("correctconfigtest.yaml")
	//deleting incorrectconfigtest.yaml file
	os.Remove("incorrectconfigtest.yaml")
}

func TestReadConfigFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "config file not present",
			args:    args{filePath: "test123.yaml"},
			wantErr: true,
		},
		{
			name:    "config file is present with test string",
			args:    args{filePath: "configtest.yaml"},
			want:    []byte("test"),
			wantErr: false,
		},
	}
	//creating test.yaml file with dummy content
	_ = os.WriteFile("configtest.yaml", []byte("test"), 0777)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfigFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
	//deleting test.yaml file
	os.Remove("configtest.yaml")
}
