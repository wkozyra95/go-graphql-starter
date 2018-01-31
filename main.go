package main

import (
	"fmt"
	"net/http"
	"os"

	conf "github.com/wkozyra95/go-graphql-starter/config"
	"github.com/wkozyra95/go-graphql-starter/web"
)

var log = conf.NamedLogger("main")

func main() {
	config, configErr := conf.SetupConfig()
	if configErr != nil {
		log.Errorf("Config error [%s]", configErr.Error())
		os.Exit(-1)
	}

	router, routerErr := web.NewRouter(config)
	if routerErr != nil {
		log.Errorf("Setup router error [%s]", configErr.Error())
		os.Exit(-1)
	}

	portStr := fmt.Sprintf(":%d", config.Port)
	log.Infof("Serving content on port %d", config.Port)
	listenErr := http.ListenAndServe(portStr, router)
	if listenErr != nil {
		log.Errorf("Server crashed [%s]", listenErr.Error())
		os.Exit(-1)
	}
}
