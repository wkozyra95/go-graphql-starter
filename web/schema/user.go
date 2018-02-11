package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
	"gopkg.in/mgo.v2/bson"
)

// UserResolver ...
type UserResolver struct {
	user *mongo.User
}

// ID ...
func (ur *UserResolver) ID() string {
	return ur.user.ID.Hex()
}

// Username ...
func (ur *UserResolver) Username() string {
	return ur.user.Username
}

// Email ...
func (ur *UserResolver) Email() string {
	return ur.user.Email
}

// Projects ...
func (ur *UserResolver) Projects(ctx context.Context) ([]*ProjectResolver, error) {
	db := extractDBSession(ctx)
	userID := extractUserID(ctx)

	if userID != ur.user.ID {
		return nil, errors.ErrUnauthorized
	}

	projects := []mongo.Project{}
	projectErr := db.Project().Find(bson.M{
		mongo.UserForeignKey: ur.user.ID,
	}).All(&projects)
	if projectErr != nil {
		return nil, errors.ErrInternalServerError
	}

	resolvers := make([]*ProjectResolver, len(projects))
	for i := range projects {
		resolvers[i] = &ProjectResolver{project: &projects[i]}
	}
	return resolvers, nil
}
