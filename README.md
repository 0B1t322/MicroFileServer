# MicroFileServer
[![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/MicroFileServer?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=104&branchName=master) [![Azure DevOps tests (compact)](https://img.shields.io/azure-devops/tests/rtuitlab/RTU%2520IT%2520Lab/104/master?compact_message)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=104&branchName=master)

Service for storing small files

## Documantation
Can be open directly by swagger. All swagger files located in ```src/MicroFileServer/docs```.
Or if the service is running in default settings http://localhost:8080/api/mfs/swagger/

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

### Required parameters

If you want run application in production mode you need to add next paramets:

in enviroment:
```.env
// url to jwks.json
MFS_AUTH_KEY_URL=https://example.com
// issuer fow jwt
MFS_AUTH_ISSUER=https://example.com
MFS_APP_TEST_MODE=false
```

or in json configs:
file ```config.json```:
```js
{
    // other params....
    // ...
    "AppOptions": {
    "testMode": false,
  }
}
```

file ```auth_config.json```:
```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json
    "issuer" : "https://exampleissuersite.com", //issuer for JWT
    // other params...
    // ...
  }
}
```

### Standart params that can be override
You can override standart params of this applicaton

in ```config.json```:

```js
{
  "DbOptions": {
    "uri": "mongodb://user:password@localhost:27017/MicroFileServer", //uri connection to mongodb
  },
  "AppOptions": {
    "appPort": "8080", //app port
    "maxFileSize": 100, //maximum file size for upload in MB
  }
}
```

in ```auth_config.json```:

```js
{
  "AuthOptions": {
    "audience": "example_audience", //audince for JWT, claim where roles will search
    "roles": {
        "user": "user", // user role that will be check if admin role not found
        "admin": "mfs.admin" // admin role
    }
  }
}

```

Or it can be override by enviroment:
```.env
# user role that search if not find admin role
MFS_AUTH_ROLE_USER=user
MFS_AUTH_ROLE_ADMIN=mfs.admin

# claim where the roles will be searched
MFS_AUTH_AUDIENCE=itlab

MFS_APP_PORT=8080

# max uploading file size in MB
MFS_APP_MAX_FILE_SIZE=100

# mongodb uri
MFS_MONGO_URI=mongodb://root:root@mfs.db:27017/MFS?authSource=admin
```

## Build

### Requirements
- Go 1.16+


if you want to run in nativaly by golang, you should have enviroments params in launching directory

in ```src/MicroFileServer``` write:
```
go build -o main
```
that will produce main binary file that can be launch by
```
./main
```

## Build using docker
after creating .env file in root directory write command in this directory:
```
docker-compose -f docker-compose.yaml up --build
```

server will run in ```http://localhost:8080```

## Standart launching with docker
To run application in testmode with standart params write in root directory
```
docker-compose -f docker-compose.override.yaml up --build
```

