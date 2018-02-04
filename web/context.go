package web

import (
	"context"
	"net/http"

	database "github.com/wkozyra95/go-graphql-starter/model/db"
	"github.com/wkozyra95/go-graphql-starter/web/schema"
)

type dbProvider func() database.DB

func (p dbProvider) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := p()
		defer db.Close()

		newCtx := context.WithValue(
			r.Context(),
			schema.DBSessionKey,
			db,
		)
		log.Info("middleware")
		updatedRequest := r.WithContext(newCtx)
		next.ServeHTTP(w, updatedRequest)
	})
}
