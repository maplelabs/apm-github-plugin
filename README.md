# github-audit
This service observes GitHub commits, operations and sends to an observability platform


# Build Instruction
- Project uses golang 1.18 . 
- Make sure golang environment is setup.
- Simple run the Makefile to build the binary
```bash
    make build
```
- This will generate the github-audit binary named as ***github-audit***

# Usage instruction
- To view help menu 
```bash
./github-audit help
```

- To start program with config.yaml at default location i.e same as binary
```bash
./github-audit start 
```

- To stop program
```bash
./github-audit stop
```

- To check version of the program
```bash
./github-audit version
```
Note: github-audit.log file will be generated in same location as binary for checking logs

# Complete sample config.yaml
```yaml
## provides logging level <OPTIONAL> , Default: info
# loglevel: debug
## folder path to log file <OPTIONAL> , Default: same as git-audit binary location 
# logpath: ./test.log 
auditJobs:
## audit job name <REQUIRED>
- name: auditjob1
## polling interval to fetch data format: 30s, 5m , 1h , 1d  etc<REQUIRED> , Default: 5m  
  polling_interval: 30s 
## metadata if any required like tags etc
  metadata:
  tags:
    ## instance
    _tag_Name: "localdomain"
    ## project name used in snappyflow
    _tag_projectName: "apm-github-plugin-test"
    ## project app name used in snappyflow
    _tag_appName: "apm-github-plugin"
  ## git saas provider like github,bitbucket etc <REQUIRED>
  repo_host: github
  ## git repository name  <REQUIRED>
  repo_name: github-audit
  ## git repository owner  <REQUIRED>
  repo_owner: nikhil-dot-kumar  
  repo_config:
  ## absolute url of repository <REQUIRED>
    repo_url: https://github.com/nikhil-dot-kumar/github-audit
    ## private or public repository <REQUIRED>
    repo_type: public      
    ## credentials to access repository data <REQUIRED for private repo>
    credentials:  
      ## username is required    <REQUIRED for private repo>
      username: Nikhil-dot-Kumar  
      ## API token in base64 encode format. <REQUIRED for private repo> , cannot be empty
      access_token: xxxxx
    ## (optional) by default all branches will be monitored
    branches:
    - test
  ## output contains target list
  output:   
    target_name:
    - kafka1
    - es1
## target list given as global configuration
targets:    
- name: kafka1
  type: kafka-rest
  config:
    host: localhost
    port: "443"
    protocol: https
    token: xxxxx
    path: kafkapath
    topic: test-topic
- name: es1
  type: elasticsearch
  config:
    ## snappyflow url and access credential 
    host: localhost
    port: "443"
    protocol: https
    index: test-index
    username: test-user
    password: xxxx
    ## before elasticsearch 7x
    old_es: "false"
    ## name of project in snappyflow
    path: apm-github-plugin-test
```
