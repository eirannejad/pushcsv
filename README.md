<img src="art/logo.svg" width="128"/>

# pushcsv
Push csv/tsv data to database

#### Supported Backends

|Database   |Driver                                 |
|-----------|---------------------------------------|
|postgresql |github.com/lib/pq                      |
|mongodb    |gopkg.in/mgo.v2                        |
|mysql      |github.com/go-sql-driver/mysql         |
|sqlserver  |github.com/denisenkom/go-mssqldb       |
|sqlite3    |github.com/mattn/go-sqlite3            |

### Getting pushcsv

```shell
go install github.com/eirannejad/pushcsv
```

Build process on Windows requires `gcc` to support compiling the sqlite3 libraries. Therefore binaries are provided for convenience under [Releases](https://github.com/eirannejad/pushcsv/releases). Compiled `pushcsv` binaries can be installed on Windows using [Chocolatey](https://chocolatey.org/) package manager as well.

```shell
choco install pushcsv
```

### Usage

Examples of pushing `users.csv` into table (collection in case of mongodb) `users` on supported databases:

```shell
$ pushcsv postgres://user:pass@data.mycompany.com/mydb users ~/users.csv --headers --purge
$ pushcsv mongodb://user:pass@localhost:27017/mydb users ~/users.csv --map=name:fullname --map=email:userid
$ pushcsv "mysql:user:pass@tcp(localhost:3306)/tests" users ~/users.csv --purge --map=name:fullname --map=email:userid
$ pushcsv sqlserver://user:pass@my-azure-db.database.windows.net?database=mydb users ~/users.csv --purge --map=name:fullname --map=email:userid
$ pushcsv sqlite3:data.db users ~/users.csv
```

### Building from source

```shell
$ git clone git@github.com:eirannejad/pushcsv.git
$ cd pushcsv
$ go get -u ./...
$ go install .
$ pushcsv --help
```

To reduce the size of pushcsv compiled binary, install with `-ldflags "-w"` option.

```shell
go install -ldflags "-w" .
```
