# Github Audit Plugin Input Config File

Github audit plugin will be using config.yaml file for taking input from user

## Sample Input Config
```yaml
## provides logging level <OPTIONAL> , Default: info
loglevel: debug
## folder path to log file <OPTIONAL> , Default: same as git-audit binary location 
logpath: ./test/ 
auditJobs:
## audit job name <REQUIRED>
- name: auditjob1
## polling interval to fetch data to be defined in cron job format  <REQUIRED> , Default: 5 * * * *   
  polling_interval: 5 * * * *   
## metadata if any required like tags etc
  metadata:
  tags:
    tag1: tag1value
  ## git saas provider like github,bitbucket etc <REQUIRED>
  repo_type: github
  ## git repository name  <REQUIRED>
  repo_name: testRepo   
  repo_config:
  ## absolute url of repository <REQUIRED>
    repo_url: https://www.github.com/test/testRepo   
    ## credentials to access repository data <REQUIRED>
    credentials:  
      ## either email or username is required    <REQUIRED>
      username: testRepo  
      email: restRepo@test.com
      ## API token in base64 encode format. <REQUIRED> , cannot be empty
      access_token: adkslas123a1312kba
    ## (optional) by default all branches will be monitored
    branches:
    - master
  ## output contains target list
  output:   
    target_name:
    - es1
- name: auditjob2
  polling_interval: 300
  metadata:
    tags:
      tag1: tag1value
  repo_type: github
  repo_name: testRepo2
  repo_config:
    repo_url: https://www.github.com/test2/testRepo2
    credentials:
      username: testRepo2
      email: restRepo2@test.com
      access_token: adkslas123a1231312kba
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
    username: ''
    password: ''
    ip: ''
- name: kafka1
  type: kafka
  config:
    host: ''
    topic: ''
- name: webhook
  type: http
  config:
    url: https://somewebhookurl
```