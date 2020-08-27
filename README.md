# M-Lab Measure Saver

[![GoDoc](https://godoc.org/github.com/m-lab/measure-saver?status.svg)](https://godoc.org/github.com/m-lab/measure-saver) [![Build Status](https://travis-ci.com/m-lab/measure-saver.svg?branch=master)](https://travis-ci.org/m-lab/measure-saver) [![Coverage Status](https://coveralls.io/repos/github/m-lab/measure-saver/badge.svg?branch=master)](https://coveralls.io/github/m-lab/measure-saver?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/measure-saver)](https://goreportcard.com/report/github.com/m-lab/measure-saver)

This repository contains the source code for the Measure Saver service that
ingests data from the M-Lab Measure Chrome Extension and stores it into a
PostgreSQL database.

## Building the service

### Using Docker

```bash
docker build -t measurementlab/measure-saver:latest .
```

This will build a minimal Alpine Linux image containing a statically-linked measure-saver executable. This is the recommended way of building and running the service.

For further details about how the Docker build works, please read the [Dockerfile](Dockerfile)

### From the source code

If you are making changes and just want to test them:

```bash
go build ./cmd/measure-saver
```

Then run it with

```bash
./measure-saver
```

Or, you can install it in `$GOPATH/bin` with:

```bash
go install ./cmd/measure-saver
```

## Testing the service

### Running a PostgreSQL instance

Firstly, make sure you're running a PostgreSQL database locally and you have
created a user and a database for this service to use.

For the purpose of testing this application, you can just run it in a Docker
container:

```bash
docker run --name postgres-dev -e POSTGRES_PASSWORD=secret -d postgres:12.3-alpine
```

Then, spawn a psql instance into the running container:

```bash
docker exec -it postgres-dev psql -U postgres
```

...and create a database:

```text
postgres=# create database "measure-saver";
```

### Running the service

When you run the service for the first time, the needed tables are
automatically created for you:

```bash
docker run --network=host measurementlab/measure-saver:latest
```

A more complete example of how to run measure-saver with a remote PostgreSQL database, an authorized keys file and a TLS certificate:

```bash
docker run --network=host measurementlab/measure-saver:latest \
  -db.addr "myhost:5432" \
  -db.name "database-name" \
  -db.user "user" \
  -db.pass "password" \
  -keys.file "authorized_api_keys.txt" \
  -tls.cert "certs/cert.pem" \
  -tls.key "certs/key.pem"
```

For a complete and up-to-date list of the available flags, please refer to the
output of `-help`.

### Using the service

The API only exposes one REST resource: `/v0/measurements`. To send a
measurement to this API, send a `POST` request to this endpoint containing the
JSON representing a measurement's result.

Example:

```json
{
  "BrowserID": "a-unique-browser-id",
  "DeviceType": "xyz",
  "Notes": "This is a note",
  "Download": 100,
  "Upload": 50,
  "Latency": 20,
  "ClientInfo": {
    [the ClientInfo object]
  },
  "ServerInfo": {
    [the ServerInfo object]
  },
  "Results": {
    [the NDT results object]
  }
}
```

The ClientInfo, ServerInfo and Results objects are defined
[here](internal/model/measurement.go).
