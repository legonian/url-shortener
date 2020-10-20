# url-shortener

Example application written in Go. App using Echo web framework package and
PostgreSQL driver for SQL.

## Usage

Before running set ``$PORT`` to set server port and ``$DATABASE_URL`` for
PostgreSQL connection string. Also ``create_pg.sql`` script is presented to
create required table in PostgreSQL database.

App using memory cache for redirected URLs. Cache limited by size and time
duration, that defined in CACHE_LIMIT and CACHE_DURATION constants in
``/database/cache.go``.

## Installation

To install:

```
$ git clone https://github.com/legonian/url-shortener
$ cd url-shortener
$ go get -d
$ go build
```

## Test

App include testing of initialization, database, cache, route requests.

To run all tests:

```
go test -v ./...
```

## Demo

Demo app is running in Heroku Cloud:
https://rocky-bayou-89648.herokuapp.com/

## Author

Oleh Ihnatushenko