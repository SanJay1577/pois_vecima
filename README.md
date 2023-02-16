# POIS

Placement Opportunity Information Service is a component which used to validate the correctness of the Advertising Markers(SCETE-35), Which uses two web interfaces  

RESTful - Customer back-office for the CRUD of channel schedules(CCMS) and alias
SOAP - CableLabs OC-SP-ESAM API for the Live Transcoders to access

The transcoder interacts with the POIS through the Event Signalling and Management(ESAM) interface. The POIS system also interacts with the back-office automation system in order to receive regular programme scheduling information (CCMS).

Live transcoders will make a Signal Processing Event(SPE) SCC-ESAM call to the POIS server with a SCTE35 payload.  The POIS server will confirm that the local break in the SCTE35 payload matched the expected window for local placement in the CCMS schedule. Then, it will respond with a Signal Processing Notification (SPN). 

## Prerequisites
go version 1.19

PostgreSQL version 12 or higher 

## Setup
### Install Go
[Install the latest version of `Go`](https://go.dev/doc/install) (or at least v1.19 to avoid issues).

### Clone repo:
git clone git@git.eng.vecima.com:cloud/pois.git

Move the Cloned repo under `$GOPATH/src/`. 
For example:  `/Users/your.username/go/src/pois`

### Environment Variables
Set the following variables in your bash profile:
```bash
export GOPRIVATE="git.eng.vecima.com"
export GOINSECURE="git.eng.vecima.com"
```
`GOPRIVATE` will tell `Go`'s internal resolver not to look on the `Go` servers for this package / package hash.  
`GOINSECURE` tells `Go` to skip certificate validation when you install dependencies from "git.eng.vecima.com". Though the certificate _is_ valid, it still gets flagged for some reason. Alternatively, you could add the certificate to your key chain and trust it.

### Install Dependencies
Run
```
go mod vendor
```
to download dependencies and place them into the project's local vendor directory.

## Setup Database
Default config assumes your postgres instance is on <host>:5432 and will connect as user `postgres` with "password" as the password. See [dbconfig.yml](./config/dbconfig.yml).

### To create new databse and set password
```
sudo -u postgres psql
CREATE DATABASE pois;
ALTER USER postgres PASSWORD 'password'; 
```

## Running the Application
Run the application directly with
``
go run main.go
``
## Port and Endpoints

Default ports used 
```
ALIAS      : 8130
CCMS       : 8130
ESAM-HTTP  : 4056
```

The endpoint for the Channel Alias Interface and Schedule will support the  HTTP RESTful operations for GET, PUT and DELETE.

Endpoints 
```
ALIAS : http://<host>:8130/pois/v1/channels/alias/{channelname}
CCMS  : http://<host>:8130/pois/v1/channels/{channelname}/{date}
ESAM  : http://<host>:4056/esam/v1/{provider}/request
```

Sample Curl Get Requests
```
ALIAS : curl -v http://<host>:8130/pois/v1/channels/alias/cnn 
CCMS  : curl -v http://<host>:8130/pois/v1/channels/cnn/02062023
Esam  : curl -v http://<host>:4056/esam/v1/comcast/request
```

#### Build & Run executable

go build -o pois
`./pois`

## Module Cleanup
If you have tested adding new dependencies, running
```
go mod tidy
```
will cleanup and remove references to any unused modules in the go.mod file.

## Running Tests
Use `go test`
```
go test ./...
```

## Apply DB Migrations
To migrate up to the new version
```
go run main.go --migrateUp
```

To migrate down to the previous version
```
go run main.go --migrateDown
```

See migartion scripts under ./pois/migrations/

## Setup pg_cron
Install pg_cron with version support to postgres DB
`
https://github.com/citusdata/pg_cron
`

To start the pg_cron background worker when PostgreSQL starts, you need to add pg_cron to shared_preload_libraries in `postgresql.conf`.

```
# required to load pg_cron background worker on start-up
shared_preload_libraries = 'pg_cron'
```
By default, the pg_cron background worker expects its metadata tables to be created in the "postgres" database. However, you can configure this by setting the cron.database_name configuration parameter in `postgresql.conf`.
```
# optionally, specify the database in which the pg_cron background worker should run (defaults to postgres)
cron.database_name = 'pois'
```

After restarting PostgreSQL, you can create the pg_cron functions and metadata tables using CREATE EXTENSION pg_cron.
```
-- run as superuser:
CREATE EXTENSION pg_cron;

-- optionally, grant usage to regular users:
GRANT USAGE ON SCHEMA cron TO postgres;
```
## Generate API Docs
To generate swagger yaml run the make command

`make swagger`

View swagger api documentation using `http://IP Address:4056/docs` endpoint

## Prometheus
After running pois application we can get the prometheus metric in below url 
`http://IP Address:2244/metrics`

