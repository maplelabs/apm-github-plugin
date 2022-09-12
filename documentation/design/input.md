# Github Audit Plugin Input Config File

Github audit plugin will be using config.yaml file for taking input from user

## Sample Input Config
```yaml
auditJobs:
- name: auditjob1   # audit job name
  polling_interval: 300   # polling interval to fetch data
  metadata:   # metadata if any required like tags etc
    tags:
      tag1: tag1value
  repo_type: github   # git saas provider like github,bitbucket etc
  repo_name: testRepo   # git repository name
  repo_config:
    repo_url: https://www.github.com/test/testRepo    # absolute url of repository
    credentials:      #credentials to access repository data
      username: testRepo  # either email or username is required
      email: restRepo@test.com
      access_token: adkslas123a1312kba
    branches:    # (optional) by default all branches will be monitored
    - master
  output:   # output contains target list
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
targets:    #target list given as global configuration
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