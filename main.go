package main

import (
	"fmt"
	"os"
	"os/signal"
	"pois/app"
	"pois/config"
	"pois/prometheus"
	"pois/shared"
	"pois/store"
	"pois/version"
	"syscall"
	"time"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
	flags "github.com/jessevdk/go-flags"
	"github.com/labstack/echo"
	migrate "github.com/rubenv/sql-migrate"
)

type options struct {
	Environment     string `default:"development" env:"POIS_ENV" short:"e" long:"environment" description:"Environment name to use for lookup of config values"`
	DBMigrateUp     bool   `env:"MIGRATION_UP_ENV" short:"u" long:"migrateUp" description:"Execute pending database migrations"`
	DBMigrateDown   bool   `env:"MIGRATION_DOWN_ENV" short:"d" long:"migrateDown" description:"Undo the most recent database migration"`
	ConfigDir       string `default:"" env:"CONFIG_DIR" short:"c" long:"confDir" description:"Directory path to the configuration file"`
	Version         bool   `env:"VERSION" short:"v" long:"version" description:"POIS application version"`
	AppliactionName bool   `env:"APPLICATION_NAME" short:"a" long:"appName" description:"POIS application full name"`
}

var migrations *migrate.FileMigrationSource = &migrate.FileMigrationSource{Dir: "migrations"}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	var opts options
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println("Failed to parse environment variables", err.Error())
		return
	}

	if opts.Version {
		fmt.Println(version.Version.Version())
		return
	} else if opts.AppliactionName {
		fmt.Println(version.Version.ApplicationName())
		return
	}
	//Initialize to parse the config files. It returns a fatal error if initialization fails
	var cfgDir string
	if opts.ConfigDir != "" {
		cfgDir = opts.ConfigDir
	} else {
		current_dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to form current system directory ", err.Error())
			return
		}
		cfgDir = current_dir + "/config"
	}

	// initilializing config files
	if err := config.Initialize(cfgDir); err != nil {
		app.Logger.Fatalf("config Initialization failure: %s", err.Error())
	}

	//Initialize the application. It is a fatal error if initialization fails
	if err := app.Initialize(cfgDir); err != nil {
		app.Logger.Fatalf("Initialization failure: %s", err.Error())
	}
	//Intilalize the Esam Application
	go app.EsamIntialize()
	err = prometheus.ReadPrometheusStatsFromJson(app.Logger)
	if err != nil {
		app.Logger.Errorf("Failed to Read prometheus stat on init  %s", err)
	}
	go app.PromInitialize()
	if opts.DBMigrateUp {
		app.Logger.Infof("DB Migration Up")
		migrateUp(app.Logger, opts.Environment)
		return
	} else if opts.DBMigrateDown {
		app.Logger.Infof("DB Migration Down")
		migrateDown(app.Logger, opts.Environment)
		return
	}

	// Provide the store to the request handler.
	shared.Stores = store.New(app.Logger, opts.Environment)
	shared.LoadChannelScheduleAndAlias()

	shared.InitailizeCleanUp(time.Duration(30) * time.Second)

	e := echo.New()
	baseGroup := e.Group("")
	baseGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("store", shared.Stores)
			c.Set("environment", opts.Environment)
			return next(c)
		}
	})
	defer zaplogger.LogSync()
	app.Logger.Infof("%s is ready", "Pois Application")
	// We are up and running so wait until we are told to stop.
	<-stopChan
	err = prometheus.WritePrometheusToJsonFile(app.Logger)
	if err != nil {
		app.Logger.Errorf("Failed to write prometheus stat to file", err.Error())
	}
	app.Logger.Infof("Application exited")
}

func migrateUp(log *zaplogger.Logger, env string) {
	migrate.SetTable("migrations")
	dbx := store.MustOpenDB(log, env)
	n, err := migrate.Exec(dbx.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("could not apply migrations error: %v", err)
	}
	log.Infof("migrations applied count: %v", n)
}

func migrateDown(log *zaplogger.Logger, env string) {
	migrate.SetTable("migrations")
	dbx := store.MustOpenDB(log, env)
	n, err := migrate.ExecMax(dbx.DB, "postgres", migrations, migrate.Down, 1)
	if err != nil {
		log.Fatalf("could not undo migration error: %v", err)
	}
	log.Infof("migration undone count: %v", n)
}
