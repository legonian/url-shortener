# url-shortener

URL Shortener web application written in Go. With this app you can
type some long link and get easy to share short one with a basic info about
number of click your short link has.

## Features

This app include following features:
+ Using net/http with Chi-router
+ PostgreSQL for database to store shortcuts
+ In-memory cache with length and durations limit
+ Clicks counter
+ Separate database code, so its easy to change database implementations
+ Tests to check different routes and functions

## Demo

Demo app example is running in Heroku Cloud (free dynos):

https://rocky-bayou-89648.herokuapp.com/

<img align="left" width="400" alt="Index Page" src="https://github.com/legonian/url-shortener/raw/master/public/img/example1.png">
<img width="400" alt="Info Page" src="https://github.com/legonian/url-shortener/raw/master/public/img/example2.png">

## Usage

Before running set ``$PORT`` to set server port and ``$DATABASE_URL`` for
PostgreSQL connection string. Also ``create_pg.sql`` script is presented to
create required table in PostgreSQL database.

App using memory cache for redirected URLs. Cache limited by size and time
duration, that defined in ``CACHE_LIMIT`` and ``CACHE_DURATION`` constants in
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

## Author

Oleh Ihnatushenko