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
$ go install -ldflags "-w" .
$ pushcsv --help
```

#### Usage

Insert `users.csv` into table `users`:

```
$ pushcsv postgres://user:pass@data.mycompany.com/mydb users ~/users.csv --headers --purge
$ pushcsv mongodb://localhost:27017/mydb users ~/users.csv --map=name:fullname --map=email:userid
$ pushcsv "mysql:ein:test@tcp(localhost:3306)/tests" users ~/users.csv --purge --map=name:fullname --map=email:userid
$ pushcsv sqlite3:data.db users ~/users.csv
```