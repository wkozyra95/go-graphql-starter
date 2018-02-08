package schema

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	mongo "github.com/wkozyra95/go-graphql-starter/model/db"
)

func (r Resolver) Project(
	ctx context.Context,
	args struct {
		ID string
	},
) (*projectResolver, error) {
	db := extractDBSession(ctx)
	userId := extractUserIdContext(ctx)

	projectId, projectIdErr := mongo.ConvertToObjectId(args.ID)
	if projectIdErr != nil {
		return nil, projectIdErr
	}

	project := model.Project{}
	dbErr := db.Project().FindID(projectId).One(&project)
	if dbErr == mgo.ErrNotFound {
		return nil, errors.NotFound
	}
	if dbErr != nil {
		return nil, errors.InternalServerError
	}
	if project.UserID != userId {
		return nil, errors.Unauthorized
	}
	return &projectResolver{project: &project}, nil
}

func (r Resolver) Projects(
	context context.Context,
) ([]*projectResolver, error) {
	db := extractDBSession(context)
	userId := extractUserIdContext(context)

	projects := []model.Project{}
	projectErr := db.Project().Find(bson.M{
		mongo.UserForeignKey: userId,
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
