# pushcsv
Push CSV files easily to database (supports psql)

#### Supported Backends

* postgres
* sqlite
* mongodb

## Getting Started

```
$ git clone git@github.com:eirannejad/pushcsv.git
$ cd pushcsv
$ go get -u ./...
$ go install .
$ pushcsv --help
```

#### Usage

Insert `users.csv` into table `users`:

```
$ pushcsv postgres://localhost:5432/mydb users ~/users.csv
```

