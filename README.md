# pushcsv
Push csv/tsv data to database

#### Supported Backends

* postgres
* mongodb
* mysql
* sqlite

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
$ pushcsv postgres://user:pass@localhost:5432/mydb users ~/users.csv
```