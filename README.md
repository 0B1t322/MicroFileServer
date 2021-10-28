# MicroFileServer
Service for storing small files

## Documantation
Can be open directly by swagger. All swagger files located in ```src/MicroFileServer/docs```.
Or if the service is running in default settings http://localhost:8081/api/mfs/swagger/

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
// Database settings
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=root
MONGO_INITDB_DATABASE=MFS


MFS_MONGO_URI=mongodb://user:password@host:port/DBName?authSource=admin

// Roles int "itlab" claim
MFS_AUTH_ROLE_USER=user
MFS_AUTH_ROLE_ADMIN=mfs.admin
// url to jwks.json
MFS_AUTH_KEY_URL=https://example.com
// audience for jwt
MFS_AUTH_AUDIENCE=claim
// issuer fow jwt
MFS_AUTH_ISSUER=https://example.com
// app port
MFS_APP_PORT=8081
// testmode can be true or false
MFS_APP_TEST_MODE=true
// max file size that can be upload
MFS_APP_MAX_FILE_SIZE=100
```

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