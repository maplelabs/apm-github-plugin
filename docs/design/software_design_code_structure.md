# Git-audit Code Stucture
Git-audit is using Golang as primary language for development. 
Below is the code structure for git-audit plugin
```
├── cmd
│   └──main.go
├── commands
│   ├── rootCmd.go
│   ├── start.go
│   └── stop.go
├── publisher
│   ├── publisher.go
│   ├── stdout.go
│   ├── kafka.go
│   └── elasticsearch.go
├── logger
│   ├── logger.go
│   └── zap.go
├── metricFormator
├── input
├── internal
│   ├── configurator
│   ├── taskManager
│   ├── dataProcessor
│   └── task
├── inMemDatastore
└── gitProvider
    ├── gitProvider.go
    ├── github.go
    ├── bitbucket.go
    └── gitlabs.go
```

# 1. ***cmd/main.go***
1. It is the starting file for the git-autdit software. 
2. It will contain the logic to CLI interface and initiating the git-audit

# 2. ***commands***
1. Package commands contains CLI commands as supported by git-audit plugin
2. Currently , ***start*** and ***stop*** commands are supported
3. It can be extended in future to support more commands

# 3. ***publisher***
1. Package publisher contains logic related to pushing data to targets
2. It can also be separately added to different repository, this is self contained package

# 4. ***metricFormator***
1. Package metricFormator contains config file to interprest APIs response

# 5. ***input***
1. Package input contains logic relted to input interface

# 6. ***inMemDatastore***
1. Package inMemDatastore contains golang struct for storing tasks related data
2. It is also be responsible for checkpointing data to ***checkpoint.json*** file

# 7. ***internal***
This directory will contain git-audit core engine logic 
1. ***configurator***
 - Package configurator contains logic to configurator

2. ***taskManager***
 - Package taskManager contains logic related to taskManager scheduling

3. ***task***
 - Package task contains task related logic

4. ***dataProcessor***
 - Package dataProcessor contains data processing related logic 

# 8. ***gitProvider***
1. Package gitProvider contains logic related to APIs calls to various VCS cloud softwares

# 9. ***logger***
1. Package logger contains logic related to logging 