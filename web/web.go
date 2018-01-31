// Package web ...
package web

import (
	"net/http"

	conf "github.com/wkozyra95/go-graphql-starter/config"
	"github.com/wkozyra95/go-graphql-starter/model/db"
)

var log = conf.NamedLogger("web")

type serverContext struct {
	config conf.Config
	db     func() db.DB
	jwt    jwtProvider
}

// NewRouter ...
func NewRouter(config conf.Config) (http.Handler, error) {
	dbCreator, dbErr := db.SetupDB(config)
	if dbErr != nil {
		log.Error(dbErr.Error())
		return nil, dbErr
	}

	jwt, jwtErr := newJwtProvider(config)
	if jwtErr != nil {
		log.Error(jwtErr.Error())
		return nil, dbErr
	}

	context := serverContext{
		config: config,
		db:     dbCreator,
		jwt:    jwt,
	}

	router, setupRoutesErr := setupRoutes(context, config)
	if setupRoutesErr != nil {
		log.Error(dbErr.Error())
		return nil, setupRoutesErr
	}

	return router, nil
}
