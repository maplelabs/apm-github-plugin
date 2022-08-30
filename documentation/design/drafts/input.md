# Github Audit Plugin Input Config File

Github audit plugin will be using config.yaml file for taking input from user

## Sample Input Config
```yaml
providers:
  github:       # type of provider , can be butbucket or gitlabs also or normal private self hosted git 
  - account_name: maplelabs     #github account username
    periodicity: 300        #preiodicity for pulling changes using github APIs
    credentials:
      username:     #either username or email is needed (mandatory)
      email:
      token: XXX        #github access token to be used as it is current github auth standard 
    repo_list:      #list of repository names
    - github_audit
    - log_generator
    publisher:          
    - es
    - kafka1
    metadata:       #metadata information required if any 
      tags:
        project: p1
  - account_name: snappyflow
    periodicity: 200
    credentials:
      username:
      email:
      token: XXX
    repo_list:
    - repo1
    - repo2
    publisher:
    - es
    - kafka1
    metadata:
      tags:
        project: p2

targets:        #list of targets
- name: es      #target name should be unique
  type: elasticsearch       #target type
  config:       #target config based on targets (can be different for each target) 
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