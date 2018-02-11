package mongo

import (
	"github.com/wkozyra95/go-graphql-starter/errors"
	"gopkg.in/mgo.v2/bson"
)

func ConvertToObjectID(id string) (binaryID bson.ObjectId, convertErr error) {
	defer func() {
		if err := recover(); err != nil {
			binaryID = ""
			convertErr = errors.Malformed
		}
	}()
	return bson.ObjectIdHex(id), nil
}

func ValidateReadRights(
	id bson.ObjectId,
	userID bson.ObjectId,
	collection Collection,
) (bool, error) {
	document := Document{}
	documentErr := collection.FindID(id).One(&document)
	if documentErr != nil {
		return false, documentErr
	}

	if document.UserID == "" {
		return document.ID == userID, nil
	}
	return document.UserID == userID, nil
}

func ValidateWriteRights(
	id bson.ObjectId,
	userID bson.ObjectId,
	collection Collection,
) (bool, error) {
	document := Document{}
	documentErr := collection.FindID(id).One(&document)
	if documentErr != nil {
		return false, documentErr
	}
	if document.UserID == "" {
		return document.ID == userID, nil
	}
	return document.UserID == userID, nil
}
