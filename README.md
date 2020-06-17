# M-Lab Measure Saver

[![GoDoc](https://godoc.org/github.com/m-lab/measure-saver?status.svg)](https://godoc.org/github.com/m-lab/measure-saver) [![Build Status](https://travis-ci.com/m-lab/measure-saver.svg?branch=master)](https://travis-ci.org/m-lab/measure-saver) [![Coverage Status](https://coveralls.io/repos/github/m-lab/measure-saver/badge.svg?branch=master)](https://coveralls.io/github/m-lab/measure-saver?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/measure-saver)](https://goreportcard.com/report/github.com/m-lab/measure-saver)

This repository contains the source code for the Measure Saver service that
ingests data from the M-Lab Measure Chrome Extension and stores it into a
PostgreSQL database.

## Testing the service

### Running a PostgreSQL instance

Firstly, make sure you're running a PostgreSQL database locally and you have
created a user and a database for this service to use.

For the purpose of testing this application, you can just run it in a Docker
container:

```bash
docker run --name postgres-dev -d postgres:12.3-alpine
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
  "Results": {
    [the NDT results object]
  }
}
```

The NDT Results object is defined [here](internal/model/measurement.go).
