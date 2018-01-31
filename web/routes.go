// Package web ...
package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	conf "github.com/wkozyra95/go-graphql-starter/config"
	"github.com/wkozyra95/go-graphql-starter/model/db"
)

const (
	paramProjectID = "projectId"
)

type requestCtx struct {
	server *serverContext
	chi    *chi.Context
	db     db.DB
}

type wrappedHandlerFunc = func(w http.ResponseWriter, r *http.Request, ctx requestCtx) error

func setupRoutes(context serverContext, config conf.Config) (http.Handler, error) {
	router := chi.NewRouter()
	f := func(handler wrappedHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			db := context.db()
			defer db.Close()
			err := handler(w, r, requestCtx{
				server: &context,
				chi:    chi.RouteContext(r.Context()),
				db:     db,
			})
			handleRequestError(w, err)
		}
	}

	router.Route("/auth", func(router chi.Router) {
		router.Post("/login", f(loginHandler))
		router.Post("/register", f(registerHandler))
	})
	router.Route("/projects", func(router chi.Router) {
		routeRoot := "/"
		routeID := fmt.Sprintf("/{%s}", paramProjectID)
		router.Use(context.jwt.middleware)

		router.Get(routeRoot, f(getProjectsHandler))
		router.Get(routeID, f(getProjectHandler))
		router.Post(routeRoot, f(createProjectHandler))
		router.Put(routeID, f(updateProjectHandler))
		router.Delete(routeID, f(deleteProjectHandler))
	})
	return router, nil
}
