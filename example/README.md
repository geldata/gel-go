# Example

This example demonstrates basic usage of the `gel-go` client to interact with a Gel (formally EdgeDB) database. It shows how to create, read, update, and delete movies via a simple HTTP API.

## Prerequisites

* Go 1.20+ installed
* `gel` CLI installed ([installation guide](https://docs.geldata.com/reference/using/cli))

## Setup

Initialize the Gel project and run migrations to create the database schema:

```bash
gel init
```

Next, start the http server.

```bash
go run main.go
```

## API Usage

**Create Movie**

```bash
curl -X POST \
  --json '{"Title":"Pupl Fiction","Year":1994,"Description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption."}' \
  --location http://localhost:8080/movie
```

**List Movies**

```bash
curl --location http://localhost:8080/movies
```

**Get Movie**

```bash
curl --location http://localhost:8080/movie?id=<uuid>
```

**Update Movie**

```bash
curl -X PUT \
  --json '{"id": "<uuid>", "Title": "Pulp Fiction"}' \
  --location http://localhost:8080/movie
```

**Delete Movie**

```bash
curl -X DELETE \
  --location http://localhost:8080/movie?id=<uuid>
```