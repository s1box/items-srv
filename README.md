## Items microservice

This is a microservice written in Golang that manages `Item` entities.

The microservice exposes `Item` REST API:

| Method | Path  | Description  |
|--------|-------|--------------|
| GET    | /items        | Returns all `Items` |
| GET    | /items/<id>   | Returns `Item` by its ID |
| GET    | /items/random | Returns randomly picked existing `Item` |
| POST   | /items        | Creates a new `Item` |
| DELETE | /items/<id>   | Deletes `Item` by its ID |

Also, there are few additional endpoints:

| Method | Path      | Description  |
|--------|-----------|--------------|
| GET    | /status   | Shows stats of the service |
| GET    | /items/db | Creates required database structure (needs admin rights) |

### How to run tests

To run all unit tests, execute:

```sh
go test ./...
```

### How to build

To build the project, run:
```sh
go build -o items-srv
```

Resulting binary would be called `items-srv`

If you need fully static binary, set `CGO_ENABLED=0` environment variable before the build.

### How to run the microservice

1. Build the microservice as described in `How to build` section
2. Setup SQL database.
3. Configure microservice by setting the following environment variables:

| Environment variable name | Example value | Description |
|---------------------------|---------------|-------------|
| `HOSTNAME`                | `0.0.0.0`     | Host on which the service is listening |
| `PORT`                    | `8080`        | Port on which the servie is listening |
| `DB_USERNAME`             | `user`        | Database user login |
| `DB_PASSWORD`             | `password`    | Database user password |
| `DB_HOSTNAME`             | `sql-host`    | Host on which the database is running |
| `DB_PORT`                 | `3306`        | Port on which the database is running |

4. Run the compiled binary
5. Check the service status by quering `/status` endpoint (optional)

### License

GNU GPL v2 or any later version.
