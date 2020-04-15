# ITLab-Reports
Service for storing work reports

Status | master | develop
---|---|---
build | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab-Reports?branchName=develop)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)
test | [![master tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/master?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=master) | [![develop tests](https://img.shields.io/azure-devops/tests/RTUITLab/RTU%20IT%20Lab/86/develop?label=%20&style=plastic)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=86&branchName=develop)

## Configuration

File ```src/ITLabReports/api/config.json``` must contain next content:

```js
{
  "DbOptions": {
    "host": "mongo", //host to mongodb server
    "port": "27017", //port to mongodb server
    "dbname": "db", //name of db in mongodb
    "collectionName": "collection" //name of collection in mongodb
  },
  "AppOptions": {
    "testMode": true|false //bool option for enabling Tests mode
  }
}
```

File ```src/ITLabReports/api/auth_config.json``` must contain next content:

```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json
    "audience": "example_audience", //audince for JWT
    "issuer" : "https://exampleissuersite.com", //issuer for JWT
    "scope" : "my_scope" //required scope for JWT
  }
}

```
## jwks.json
Example for jwks.json file
```js
{
   "keys":[
      {
         "kty":"RSA", //The family of cryptographic algorithms used with the key.
         "use":"sig", //How the key was meant to be used; sig represents the signature.
         "x5c":[ //The x.509 certificate chain. The first entry in the array is the certificate to use for token verification; the other certificates can be used to verify this first certificate.
            "123"
         ],
         "n":"456", //The modulus for the RSA public key.
         "e":"AQAB", //The exponent for the RSA public key.
         "kid":"789" //The unique identifier for the key.
      }
   ]
}
```
