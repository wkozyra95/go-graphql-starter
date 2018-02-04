// Package web ...
package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/neelance/graphql-go/relay"
	conf "github.com/wkozyra95/go-graphql-starter/config"
	"github.com/wkozyra95/go-graphql-starter/model/db"
	"github.com/wkozyra95/go-graphql-starter/web/schema"
)

var log = conf.NamedLogger("web")

// NewRouter ...
func NewRouter(config conf.Config) (http.Handler, error) {
	dbCreator, dbErr := db.SetupDB(config)
	if dbErr != nil {
		log.Error(dbErr.Error())
		return nil, dbErr
	}

	jwt, jwtErr := newJwtProvider(config)
	if jwtErr != nil {
		log.Errorf("Create jwtProvider failed with error [%s]", jwtErr.Error())
		return nil, dbErr
	}

	context := &schema.Resolver{
		Config:        &config,
		GenerateToken: jwt.GenerateToken,
	}

	schema, schemaErr := schema.SetupSchema(context)
	if schemaErr != nil {
		log.Errorf("Create schema failed with error [%s]", schemaErr.Error())
		return nil, schemaErr
	}

	router := chi.NewRouter()

	router.Use(dbProvider(dbCreator).middleware)
	router.Use(jwt.middleware)

	handler := &relay.Handler{Schema: schema}
	router.Post("/relay", handler.ServeHTTP)
	router.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		response := schema.Exec(
			r.Context(),
			string(body),
			"",
			nil,
		)

		responseJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	})

	return router, nil
}
