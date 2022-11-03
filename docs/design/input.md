# Github Audit Plugin Input Config File

Github audit plugin will be using config.yaml file for taking input from user

## Sample Input Config
```yaml
## provides logging level <OPTIONAL> , Default: info
loglevel: debug
## folder path to log file <OPTIONAL> , Default: same as git-audit binary location 
logpath: ./test.log 
auditJobs:
## audit job name <REQUIRED>
- name: auditjob1
## polling interval to fetch data to be defined  <REQUIRED> , Default: 5s, 10m etc  
  polling_interval: 5m 
## metadata if any required like tags etc
  metadata:
  tags:
    tag1: tag1value
  ## git saas provider like github,bitbucket etc <REQUIRED>
  repo_host: github
  ## git repository name  <REQUIRED>
  repo_name: testRepo
  ## git repository owner  <REQUIRED>
  repo_owner: testOwner   
  repo_config:
  ## absolute url of repository <REQUIRED>
    repo_url: https://www.github.com/test/testRepo
    ## private or public repository <REQUIRED>
    repo_type: public      
    ## credentials to access repository data <REQUIRED for private repo>
    credentials:  
      ## username is required    <REQUIRED for private repo>
      username: testRepo  
      ## API token in base64 encode format. <REQUIRED for private repo> , cannot be empty
      access_token: adkslas123a1312kba
    ## (optional) by default all branches will be monitored
    branches:
    - master
  ## output contains target list
  output:   
    target_name:
    - es1
- name: auditjob2
  polling_interval: 300s
  metadata:
  tags:
    tag1: tag1value
  repo_host: github
  repo_name: testRepo2
  repo_owner: testOwner2
  repo_config:
    repo_type: private
    repo_url: https://www.github.com/test2/testRepo2
    credentials:
      username: testRepo2
      access_token: test123
    branches:
    - master
    - release
  output:
    target_name:
    - kafka1
## target list given as global configuration
targets:    
- name: es1
  type: elasticsearch
  config:
    username: 123
    password: 123
    ip: ''
- name: kafka1
  type: kafka
  config:
    host: 12
    topic: 12
- name: webhook
  type: http
  config:
    url: https://somewebhookurl
```