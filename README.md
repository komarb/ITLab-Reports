# ITLab-Reports
Service for storing work reports

[![Actions Status](https://github.com/itlabrtumirea/itlab-reports/workflows/go/badge.svg)](https://github.com/itlabrtumirea/itlab-reports/actions)


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