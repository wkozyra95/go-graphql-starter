package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
)

type projectResolver struct {
	project *mongo.Project
}

func (pr *projectResolver) ID() string {
	return pr.project.ID.Hex()
}

func (pr *projectResolver) Name() string {
	return pr.project.Name
}

func (pr *projectResolver) Description() string {
	return pr.project.Description
}

func (pr *projectResolver) Details() *projectDetailsResolver {
	return &projectDetailsResolver{
		details: &pr.project.Details,
	}
}

type projectDetailsResolver struct {
	details *model.ProjectDetails
}

func (pr *projectDetailsResolver) IsPublic() bool {
	return pr.details.IsPublic
}

func (pr *projectDetailsResolver) ProjectType() string {
	return pr.details.ProjectType
}

func (pr *projectResolver) User(ctx context.Context) (*userResolver, error) {
	db := extractDBSession(ctx)
	userID := extractUserID(ctx)

	user := mongo.User{}
	dbErr := db.User().FindID(pr.project.UserID).One(&user)
	if dbErr != nil {
		return nil, errors.InternalServerError
	}
	if user.ID != userID {
		return nil, errors.Unauthorized
	}
	return &userResolver{user: &user}, nil
}
