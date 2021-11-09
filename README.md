# MicroFileServer
[![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/MicroFileServer?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=104&branchName=master) [![Azure DevOps tests (compact)](https://img.shields.io/azure-devops/tests/rtuitlab/RTU%2520IT%2520Lab/104/master?compact_message)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=104&branchName=master)

Service for storing small files

## Documantation
Can be open directly by swagger. All swagger files located in ```src/MicroFileServer/docs```.
Or if the service is running in default settings http://localhost:8081/api/mfs/swagger/

## Tests
Project contains e2e tests, using TestMace

### E2E
**Requirements**:
- Node.JS
- Docker

1. Install testmace cli
    ```bash
    npm install --global @testmace/cli@1.3.1
    ```
1. Make sure that the local development server (and db) is turned off
1. Run tests
    ```bash
    ./runTests.sh
    ```
1. Results can be found in `tests-out`

## Configuration

File ```config.json``` must contain next content:

```js
{
  "DbOptions": {
    "uri": "mongodb://user:password@localhost:27017/MicroFileServer", //uri connection string | env: MFS_MONGO_URI
  },
  "AppOptions": {
    "testMode": true|false, //bool option for enabling Tests mode | env: MFS_APP_TEST_MODE
    "appPort": "8081", //app port | env: MFS_APP_PORT
    "maxFileSize": 100, //maximum file size for upload in MB | env: MFS_APP_MAX_FILE_SIZE
  }
}
```

File ```auth_config.json``` must contain next content:

```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json | env: MFS_AUTH_KEY_URL
    "audience": "example_audience", //audince for JWT | env: MFS_AUTH_AUDIENCE
    "issuer" : "https://exampleissuersite.com", //issuer for JWT | env: MFS_AUTH_ISSUER
    "roles": {
        "user": "user", // user role that will be check in itlab claim | env: MFS_AUTH_ROLE_USER
        "admin": "mfs.admin" // admin role | env: MFS_AUTH_ROLE_ADMIN
    }
  }
}

```

or you can configure by file in ```src/.env``` that should contains:
```.env
// url to jwks.json
MFS_AUTH_KEY_URL=https://example.com
// issuer fow jwt
MFS_AUTH_ISSUER=https://example.com
```
<!-- MFS_MONGO_URI=mongodb://user:password@host:port/DBName?authSource=admin -->
<!-- // audience for jwt
MFS_AUTH_AUDIENCE=claim -->
<!-- // app port
MFS_APP_PORT=8081
// testmode can be true or false
MFS_APP_TEST_MODE=true
// max file size that can be upload in MB
MFS_APP_MAX_FILE_SIZE=100 -->

## Build
### Requirements
- Go 1.16+

after creating .env file in ```src/MicroFileServer``` write:
```
go build -o main
./main
```

## Build using docker
after creating .env file in ```src``` write command in this directory:
```
docker-compose -f docker-compose.override.yml up --build
```

server will run in ```http://localhost:8081```

