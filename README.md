# url-shortener

Example application written in Go. App using Echo web framework package and
PostgreSQL driver for SQL.

## Usage

Before running set ``$PORT`` to set server port and ``$DATABASE_URL`` for
PostgreSQL connection string. Also ``create_pg.sql`` script is presented to
create required table in PostgreSQL database.

App using memory cache to redirect URLs. Cache limited by time duration, to
see current value or set it to custom go to ``/handler/cache.go`` and
change CACHE_DURATION constant.

## Installation

```
$ git clone https://github.com/legonian/url-shortener
$ cd url-shortener
$ go get -d
$ go build
```

## Demo

Demo app is running in Heroku Cloud:
https://rocky-bayou-89648.herokuapp.com/

## Author

Oleh Ihnatushenko