package schema

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
)

func (r Resolver) Project(
	ctx context.Context,
	args struct {
		ProjectID string
	},
) (*projectResolver, error) {
	db := extractDBSession(ctx)
	userID := extractUserID(ctx)

	projectID, projectIDErr := mongo.ConvertToObjectID(args.ProjectID)
	if projectIDErr != nil {
		return nil, projectIDErr
	}

	project := mongo.Project{}
	dbErr := db.Project().FindID(projectID).One(&project)
	if dbErr == mgo.ErrNotFound {
		return nil, errors.NotFound
	}
	if dbErr != nil {
		return nil, errors.InternalServerError
	}
	if project.UserID != userID {
		return nil, errors.Unauthorized
	}
	return &projectResolver{project: &project}, nil
}

func (r Resolver) Projects(
	context context.Context,
) ([]*projectResolver, error) {
	db := extractDBSession(context)
	userID := extractUserID(context)

	projects := []mongo.Project{}
	projectErr := db.Project().Find(bson.M{
		mongo.UserForeignKey: userID,
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
