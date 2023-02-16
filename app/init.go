package app

import (
	"fmt"
	"net/http"
	"strings"

	"pois/api"
	config "pois/config"
	serviceLog "pois/config/logging"
	pois "pois/pois"

	esam "pois/pois/esam"
	prom "pois/prometheus"
	"pois/version"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Logger *zaplogger.Logger

// Initialize sets up the application based on
// the configuration files provided in the configuration directory.
func Initialize(cfgDir string) error {
	// Read the configuration.  It is a fatal error if it fails.

	//Initial logger creation. It is a fatel error if it fails //pois application log initializion
	if err := zaplogger.Initialize(version.Version.Application(), config.ZapConfigData); err != nil {
		return fmt.Errorf(err.Error())
	}

	Logger = zaplogger.MainLog

	//Initialize the pois transactional level logs. It is a fatel error if it fails
	if err := pois.Initialize(); err != nil {
		return fmt.Errorf("%s - Pois Initialization Problem : %s", version.Version.ApplicationName(), err.Error())
	}

	//Initilaize Esam transactional Level Log
	if err := esam.Initialize(); err != nil {
		return fmt.Errorf("%s - Esam Intialization Problem : %s", version.Version.ApplicationName(), err.Error())
	}

	//Initialize API log level. It is a fatel error if it fails
	if err := api.Initialize(); err != nil {
		return fmt.Errorf("%s - API Initialization problem: %s", version.Version.ApplicationName(), err.Error())
	}

	return nil
}

// The below function will be initialied in the main.go function
// It will intialize the listen and serve for esam Request
func EsamIntialize() {
	serviceLog.EsamLog.Infof("API Initialization is processed ")
	serviceLog.EsamLog.Infof("Provider Length %v ", len(config.GetConfig().GetStringSlice("esam.providers")))
	for i, providerName := range config.GetConfig().GetStringSlice("esam.providers") {
		serviceLog.EsamLog.Infof("Provider %v : %v ", i, providerName)

		http.HandleFunc(strings.Replace(config.GetConfig().GetString("esam.rootpath"), "*", string(providerName), 1), esam.EsamRequest)
	}
	if config.GetConfig().GetString("esam.tlsport") != "" {
		http.ListenAndServeTLS(fmt.Sprintf(":%v", config.GetConfig().GetInt("esam.tlsport")), config.GetConfig().GetString("api.tlsCertPath"), config.GetConfig().GetString("api.tlsPrivateKeyPath"), nil)
		//.ListenAndServeTLS(config.GetEsamConfig().TlsPort, config.GetEsamConfig().Cert, config.GetEsamConfig().PrivateKey, nil)
		serviceLog.EsamLog.Infof("Server  is up and running")
		serviceLog.EsamLog.Debugf("API is Initialized properly and the port is up and running ")
		return
	}
	http.ListenAndServe(fmt.Sprintf(":%v", config.GetConfig().GetInt("esam.port")), nil)
	//http.ListenAndServe(config.GetEsamConfig().Port, nil)

	serviceLog.EsamLog.Infof("Server  is up and running")
	serviceLog.EsamLog.Debugf("API is Initialized properly and the port is up and running ")
}

// Initializing the Prometheus data in /metric port
func PromInitialize() {
	//Prometheus registries
	registry := prometheus.NewRegistry()
	registry.MustRegister(prom.CcmsRequest)
	registry.MustRegister(prom.CcmsDefaultRequest)
	registry.MustRegister(prom.AliasRequest)
	registry.MustRegister(prom.AliasDefaultRequest)
	registry.MustRegister(prom.EsamSucessRequest)
	registry.MustRegister(prom.EsamNoopRequest)
	registry.MustRegister(prom.EsamDeleteRequest)
	registry.MustRegister(prom.EsamDefaultRequest)
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	//swagger initilization
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	http.Handle(config.GetConfig().GetString("api.swaggerRootPath"), sh)
	http.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	http.Handle(config.GetConfig().GetString("api.prometheusRootPath"), promHandler)
	http.ListenAndServe(fmt.Sprintf(":%v", config.GetConfig().GetInt("api.missport")), nil)
}
