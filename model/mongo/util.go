package mongo

import (
	"github.com/wkozyra95/go-graphql-starter/errors"
	"gopkg.in/mgo.v2/bson"
)

// ConvertToObjectID ...
func ConvertToObjectID(id string) (binaryID bson.ObjectId, convertErr error) {
	defer func() {
		if err := recover(); err != nil {
			binaryID = ""
			convertErr = errors.ErrMalformed
		}
	}()
	return bson.ObjectIdHex(id), nil
}
