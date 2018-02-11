package schema

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
)

// Project ...
func (r Resolver) Project(
	ctx context.Context,
	args struct {
		ProjectID string
	},
) (*ProjectResolver, error) {
	db := extractDBSession(ctx)
	userID := extractUserID(ctx)

	projectID, projectIDErr := mongo.ConvertToObjectID(args.ProjectID)
	if projectIDErr != nil {
		return nil, projectIDErr
	}

	project := mongo.Project{}
	dbErr := db.Project().FindID(projectID).One(&project)
	if dbErr == mgo.ErrNotFound {
		return nil, errors.ErrNotFound
	}
	if dbErr != nil {
		return nil, errors.ErrInternalServerError
	}
	if project.UserID != userID {
		return nil, errors.ErrUnauthorized
	}
	return &ProjectResolver{project: &project}, nil
}

// Projects ...
func (r Resolver) Projects(
	context context.Context,
) ([]*ProjectResolver, error) {
	db := extractDBSession(context)
	userID := extractUserID(context)

	projects := []mongo.Project{}
	projectErr := db.Project().Find(bson.M{
		mongo.UserForeignKey: userID,
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
