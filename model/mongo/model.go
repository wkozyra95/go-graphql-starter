package mongo

import (
	"github.com/wkozyra95/go-graphql-starter/model"
	"gopkg.in/mgo.v2/bson"
)

// User ...
type User struct {
	// ID ...
	ID bson.ObjectId `bson:"_id"`
	// User ...
	model.User `bson:",inline"`
}

// Project ...
type Project struct {
	// ID ...
	ID bson.ObjectId `bson:"_id"`
	// UserID ...
	UserID bson.ObjectId `bson:"userId"`
	// Project ...
	model.Project `bson:",inline"`
}
