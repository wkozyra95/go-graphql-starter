package db

import (
	"github.com/wkozyra95/go-graphql-starter/errors"
	"gopkg.in/mgo.v2/bson"
)

func ConvertToObjectId(id string) (binaryId bson.ObjectId, convertErr error) {
	defer func() {
		if err := recover(); err != nil {
			binaryId = ""
			convertErr = errors.NotFound
		}
	}()
	return bson.ObjectIdHex(id), nil
}

func ValidateReadRights(
	id bson.ObjectId,
	userId bson.ObjectId,
	collection Collection,
) (bool, error) {
	document := Document{}
	documentErr := collection.FindID(id).One(&document)
	if documentErr != nil {
		return false, documentErr
	}

	if document.UserID == "" {
		return document.ID == userId, nil
	}
	return document.UserID == userId, nil
}

func ValidateWriteRights(
	id bson.ObjectId,
	userId bson.ObjectId,
	collection Collection,
) (bool, error) {
	document := Document{}
	documentErr := collection.FindID(id).One(&document)
	if documentErr != nil {
		return false, documentErr
	}
	if document.UserID == "" {
		return document.ID == userId, nil
	}
	return document.UserID == userId, nil
}
