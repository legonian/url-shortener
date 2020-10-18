# url-shortener

Example application written in Go. App using Echo web framework package and
PostgreSQL driver for SQL.

## Usage

Before running set ``$PORT`` to set server port and ``$DATABASE_URL`` for
PostgreSQL connection string. Also ``create_pg.sql`` script is presented to
create required table in PostgreSQL database.

## Installation

```
$ git clone https://github.com/legonian/url-shortener
$ cd url-shortener
$ go get -d
$ go build
```

## Author

Oleh Ihnatushenko