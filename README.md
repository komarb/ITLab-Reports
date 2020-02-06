# ITLab-Reports
Service for storing work reports

Status | master | develop
---|---|---
build | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=77&branchName=master) | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=develop)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=77&branchName=develop)
test | [![master tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/77/master?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=77&branchName=master) | [![develop tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/77/develop?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=77&branchName=develop)

## Configuration

File ```src/ITLabReports/api/config.json``` must contain next content:

```json
{
  "DbOptions": {
    "host": "host to mongodb server",
    "port": "port to mongodb server",
    "dbname" : "name of db in mongodb",
    "collectionName" : "name of collection in mongodb"
  },
  "AuthOptions": {
    "keyUrl" : "url to jwks.json",
    "audience" : "audience for JWT auth",
    "testKeyUrl" : "https://pastebin.com/raw/D7UL1cbH"
  },
  "AppOptions": {
    "testMode": true or false 
  }
}

```