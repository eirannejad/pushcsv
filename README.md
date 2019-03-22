# pushcsv
Push csv/tsv data to database

#### Supported Backends

* postgres
* mongodb
* mysql
* sqlite

### Getting pushcsv

Build process on Windows requires gcc to support compiling the sqlite3 libraries. Therefore binaries are provided for convenience under [Releases](https://github.com/eirannejad/pushcsv/releases). `pushcsv` can be installed on Windows using [Chocolatey](https://chocolatey.org/) package manager as well.

```
choco install pushcsv
```

#### Building from source

```
$ git clone git@github.com:eirannejad/pushcsv.git
$ cd pushcsv
$ go get -u ./...
$ go install -ldflags "-w" .
$ pushcsv --help
```

### Usage

Examples of pushing `users.csv` into table `users` on various databases:

```
$ pushcsv postgres://user:pass@data.mycompany.com/mydb users ~/users.csv --headers --purge
$ pushcsv mongodb://localhost:27017/mydb users ~/users.csv --map=name:fullname --map=email:userid
$ pushcsv "mysql:ein:test@tcp(localhost:3306)/tests" users ~/users.csv --purge --map=name:fullname --map=email:userid
$ pushcsv sqlite3:data.db users ~/users.csv
```