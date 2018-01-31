package handler

import (
	"net/http"
	"runtime/debug"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model/db"
	"gopkg.in/mgo.v2/bson"
)

func validateReadRights(id bson.ObjectId, collection db.Collection, context ActionContext) (bool, error) {
	document := db.Document{}
	documentErr := collection.FindID(id).One(&document)
	if documentErr != nil {
		return false, documentErr
	}

	return document.UserID == context.UserID, nil
}

func validateWriteRights(id bson.ObjectId, collection db.Collection, context ActionContext) (bool, error) {
	document := db.Document{}
	documentErr := collection.FindID(id).One(&document)
	if documentErr != nil {
		return false, documentErr
	}
	return document.UserID == context.UserID, nil
}

func internalServerErr(msg string) error {
	debug.PrintStack()
	return errors.NewSimple(
		msg,
		http.StatusInternalServerError,
		errors.ErrInternalServerError,
	)
}
