# ITLab-Reports
Service for storing work reports

Status | master | develop
---|---|---
build | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=develop)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)
test | [![master tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/master?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![develop tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/develop?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)

## Configuration

File ```src/ITLabReports/api/config.json``` must contains next content:

```json
{
  "host": "host to mongodb server",
  "port": "port to mongodb server",
  "dbname" : "database name on mongodb server",
  "collectionName" : "collection name in database"
}
```
