package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
	"gopkg.in/mgo.v2/bson"
)

type userResolver struct {
	user *mongo.User
}

func (ur *userResolver) ID() string {
	return ur.user.ID.Hex()
}

func (ur *userResolver) Username() string {
	return ur.user.Username
}

func (ur *userResolver) Email() string {
	return ur.user.Email
}

func (ur *userResolver) Projects(ctx context.Context) ([]*projectResolver, error) {
	db := extractDBSession(ctx)
	userID := extractUserID(ctx)

	if userID != ur.user.ID {
		return nil, errors.Unauthorized
	}

	projects := []mongo.Project{}
	projectErr := db.Project().Find(bson.M{
		mongo.UserForeignKey: ur.user.ID,
	}).All(&projects)
	if projectErr != nil {
		return nil, errors.InternalServerError
	}

	resolvers := make([]*projectResolver, len(projects))
	for i, _ := range projects {
		resolvers[i] = &projectResolver{project: &projects[i]}
	}
	return resolvers, nil
}
