package handler

import (
	"github.com/wkozyra95/go-graphql-starter/model/db"
	"gopkg.in/mgo.v2/bson"
)

// ActionContext ...
type ActionContext struct {
	DB     db.DB
	UserID bson.ObjectId
}
