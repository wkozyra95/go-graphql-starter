package schema

import (
	"context"
	"runtime/debug"

	"github.com/wkozyra95/go-graphql-starter/model/db"
	"gopkg.in/mgo.v2/bson"
)

type ContextKey string

var CurrentUserKey ContextKey = "CurrentUser"
var DBSessionKey ContextKey = "DBSessionKey"

func extractDBSession(ctx context.Context) db.DB {
	dbSessionObj := ctx.Value(DBSessionKey)
	if dbSessionObj == nil {
		log.Error("[ASSERT] Missing db session in context")
		debug.PrintStack()
	}
	dbSession, assertOk := dbSessionObj.(db.DB)
	if !assertOk {
		log.Error("[ASSERT] Wrong type for db session")
		debug.PrintStack()
	}
	return dbSession
}

func extractUserIdContext(ctx context.Context) bson.ObjectId {
	userIdObj := ctx.Value(CurrentUserKey)
	if userIdObj == nil {
		return ""
	}
	userId, assertOk := userIdObj.(bson.ObjectId)
	if assertOk {
		log.Error("[ASSERT] Wrong type for userId in context")
		debug.PrintStack()
	}
	return userId
}
