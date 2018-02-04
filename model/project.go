package model

import (
	"gopkg.in/mgo.v2/bson"
)

// Project ...
type Project struct {
	// ID ...
	ID bson.ObjectId `json:"id" bson:"_id"`
	// UserID ...
	UserID bson.ObjectId `json:"userId" bson:"userId"`

	// Name ...
	Name string `json:"name" bson:"name"`

	// Description ...
	Description string `json:"description" bson:"description"`
}
