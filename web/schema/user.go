package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	mongo "github.com/wkozyra95/go-graphql-starter/model/db"
	"gopkg.in/mgo.v2/bson"
)

type userResolver struct {
	user *model.User
}

func (ur *userResolver) ID() (string, error) {
	if ur.user == nil {
		return "", errors.InternalServerError
	}
	return ur.user.ID.Hex(), nil
}

func (ur *userResolver) Username() (string, error) {
	if ur.user == nil {
		return "", errors.InternalServerError
	}
	return ur.user.Username, nil
}

func (ur *userResolver) Email() (string, error) {
	if ur.user == nil {
		return "", errors.InternalServerError
	}
	return ur.user.Email, nil
}

func (ur *userResolver) Projects(ctx context.Context) ([]*projectResolver, error) {
	db := extractDBSession(ctx)
	userId := extractUserIdContext(ctx)

	if userId != ur.user.ID {
		return nil, errors.Unauthorized
	}

	projects := []model.Project{}
	projectErr := db.Project().Find(bson.M{
		mongo.UserForeignKey: ur.user.ID,
	}).All(&projects)
	if projectErr != nil {
		return nil, errors.InternalServerError
	}

	resolvers := make([]*projectResolver, len(projects))
	for i, project := range projects {
		resolvers[i] = &projectResolver{project: &project}
	}
	return resolvers, nil
}
