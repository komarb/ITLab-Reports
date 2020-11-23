# ITLab Reports API
Created with Golang, Gorilla and MongoDB

## Configuration

File ```src/api/auth_config.json``` must contain next content:

```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json | env: ITLAB_REPORTS_AUTH_KEY_URL
    "audience": "example_audience", //audince for JWT | env: ITLAB_REPORTS_AUTH_AUDIENCE
    "issuer" : "https://exampleissuersite.com", //issuer for JWT | env: ITLAB_REPORTS_AUTH_ISSUER
    "scope" : "my_scope", //required scope for JWT | env: ITLAB_REPORTS_AUTH_SCOPE
  }
}
```

File ```src/api/config.json``` must contain next content:

```js
{
  "DbOptions": {
    "uri": "mongodb://user:password@localhost:27017",    // env: ITLAB_REPORTS_MONGO_URI
    "dbName" : "ITLabReports", // env: ITLAB_REPORTS_MONGO_DB_NAME
    "collectionName" : "reports", // env: ITLAB_REPORTS_MONGO_DB_COLLECTION_NAME
  },
  "AppOptions": {
    "appPort": "8080", // env: ITLAB_REPORTS_APP_PORT
    "testMode": false      //testMode=true disables jwt validation | env: ITLAB_REPORTS_APP_TEST_MODE
  }
}
```

## Build 
### Requirements
- Go 1.13.8+

In ```src``` directory:
```
go build main.go
./main
```
## Installation using Docker
Install Docker and in project root directory write this code:
```
docker-compose up -d --build
```
If you’re using Docker natively on Linux, Docker Desktop for Mac, or Docker Desktop for Windows, then the server will be running on
```http://localhost:8080```

If you’re using Docker Machine on a Mac or Windows, use ```docker-machine ip MACHINE_VM``` to get the IP address of your Docker host. Then, open ```http://MACHINE_VM_IP:8080``` in a browser

## Requests
You can get Postman requests collection [here](https://www.getpostman.com/collections/4085657bcce140031d0c)

## DB Backup and Restore
To make a backup of a DB, open root folder where MongoDB is installed, open a command promt and type the command mongodump
To restore the backup, open root folder where MongoDB is installed, open a command promt and type the command mongorestore
(All DB paths are default)

