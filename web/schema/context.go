package schema

import (
	"context"
	"runtime/debug"

	"github.com/wkozyra95/go-graphql-starter/model/mongo"
	"gopkg.in/mgo.v2/bson"
)

type ContextKey string

var CurrentUserKey ContextKey = "CurrentUser"
var DBSessionKey ContextKey = "DBSessionKey"

func extractDBSession(ctx context.Context) mongo.DB {
	dbSessionObj := ctx.Value(DBSessionKey)
	if dbSessionObj == nil {
		log.Error("[ASSERT] Missing db session in context")
		debug.PrintStack()
	}
	dbSession, assertOk := dbSessionObj.(mongo.DB)
	if !assertOk {
		log.Error("[ASSERT] Wrong type for db session")
		debug.PrintStack()
	}
	return dbSession
}

func extractUserID(ctx context.Context) bson.ObjectId {
	userIDObj := ctx.Value(CurrentUserKey)
	if userIDObj == nil {
		return ""
	}
	userID, assertOk := userIDObj.(bson.ObjectId)
	if !assertOk {
		log.Error("[ASSERT] Wrong type for userID in contex [%+v]", userIDObj)
		debug.PrintStack()
	}
	return userID
}
