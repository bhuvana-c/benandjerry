package main

import (
	"database/sql"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"github.com/zalora/benandjerry/config"
	"github.com/zalora/benandjerry/db"
	"github.com/zalora/benandjerry/handler"
	"github.com/zalora/benandjerry/model"
)

const (
	serviceName = "benandjerry"
)

var (
	buildTimestamp = "unknown"
	commitID       = "unknown"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}

func main() {
	config, err := config.Load()
	failOnError(err, "invalid config settings")
	setupLog(config)
	log.Infof("%s built on %s from commit %s", serviceName, buildTimestamp, commitID)

	err = db.RunMigrateScripts(&config.PostgresConn, config.MigrationScriptsPath)
	failOnError(err, "error while running migration scripts")

	router := httprouter.New()
	indexHandler := &handler.IndexHandler{}
	indexHandler.AddRoutes(router)

	db, err := sql.Open("postgres", config.PostgresConn.GetURL())

	failOnError(err, "error while connecting to database")
	defer db.Close()

	iceCreamStore, err := model.NewIceCreamStore(db)
	failOnError(err, "error while initializing store")

	recipesHandler := &handler.IceCreamHandler{iceCreamStore}
	recipesHandler.AddRoutes(router)

	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(router)
	n.Run(":" + strconv.Itoa(int(config.ListenPort)))
}

func setupLog(config *config.Config) {
	setLogLevel(config)
	setLogFormat(config)
}

func setLogFormat(config *config.Config) {
	switch config.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
}

func setLogLevel(config *config.Config) {
	switch config.LogLevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
