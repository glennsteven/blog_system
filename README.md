# Blog System

### Database Migration
All database migration is in `db/migrations` folder.

### Configuration
All configuration is in `.env` file.

### instructions On How To Create The Data Are in The Route File
All route is in `internal/route/route.go`

### Download Dependency
```shell
go mod download
go mod tidy
```

### How to run application on local environment
- copy & paste file `.env.example` -> `.env`
- setup variable `.env` with your configuration local environment
- `mkdir` folder `bin` in `blog-system` folder
- run command `make build` for make a binary file golang
- for running this application, please run this command `./bin/blog-system`

### Create Migrations
```shell
migrate create -ext sql -dir db/migrations {table_name}
```


### Database Migration
migration up
```shell
migrate -database "{driver}://{username}:{password}@{host}:{port}/{database_name}?sslmode=disable" -path db/migrations up

```

migration down
```shell
migrate -database "{driver}://{username}:{password}@{host}:{port}/{database_name}?sslmode=disable" -path db/migrations down
```