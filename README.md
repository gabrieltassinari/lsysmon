# lsysmon

Lsysmon is a monitoring tool for linux that provides charts and useful information from metrics of proc file system using [lpfs](https://github.com/rprobaina/lpfs) go library.

# Usage

1. Clone the repository
```bash
$ git clone https://github.com/gabrieltassinari/lsysmon.git
$ cd lsysmon
```
2. Create postgres database and setup environment variables
```bash
$ psql -f schema.sql

$ export PGUSER="postgres"
$ export PGPASSWORD="your-password"
$ export PGHOST="localhost"
$ export PGPORT=5432
$ export PGDATABASE="lsysmon"
```
3. Build and run
```bash
$ go build
$ ./lsysmon
```
4. Open in your browser [http://localhost:8080](http://localhost:8080)
