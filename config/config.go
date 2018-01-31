// Package config implements configuration of application.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config contains basic values controlling app.
type Config struct {
	DbURL string
	Port  int64
}

var log = NamedLogger("web")

// SetupConfig reads config and creates Confid struct.
func SetupConfig() (Config, error) {
	config := Config{}
	config.DbURL = os.Getenv("DB_URL")
	if config.DbURL == "" {
		return config, fmt.Errorf("Missing DB_URL environment variable")
	}

	port := os.Getenv("BACKEND_PORT")
	parsedPort, parsePortErr := strconv.ParseInt(port, 10, 64)
	config.Port = parsedPort
	if parsePortErr != nil {
		log.Warning(
			"Missing or invalid BACKEND_PORT environment variable, using 8080 instead [Error: %s]",
			parsePortErr.Error(),
		)
		config.Port = 8080
	}
	return config, nil

}
