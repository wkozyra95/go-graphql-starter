package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
)

// ProjectResolver ...
type ProjectResolver struct {
	project *mongo.Project
}

// ID ...
func (pr *ProjectResolver) ID() string {
	return pr.project.ID.Hex()
}

// Name ...
func (pr *ProjectResolver) Name() string {
	return pr.project.Name
}

// Description ...
func (pr *ProjectResolver) Description() string {
	return pr.project.Description
}

// Details ...
func (pr *ProjectResolver) Details() *ProjectDetailsResolver {
	return &ProjectDetailsResolver{
		details: &pr.project.Details,
	}
}

// ProjectDetailsResolver ...
type ProjectDetailsResolver struct {
	details *model.ProjectDetails
}

// IsPublic ...
func (pr *ProjectDetailsResolver) IsPublic() bool {
	return pr.details.IsPublic
}

// ProjectType ...
func (pr *ProjectDetailsResolver) ProjectType() string {
	return pr.details.ProjectType
}

// User ...
func (pr *ProjectResolver) User(ctx context.Context) (*UserResolver, error) {
	db := extractDBSession(ctx)
	userID := extractUserID(ctx)

	user := mongo.User{}
	dbErr := db.User().FindID(pr.project.UserID).One(&user)
	if dbErr != nil {
		return nil, errors.ErrInternalServerError
	}
	if user.ID != userID {
		return nil, errors.ErrUnauthorized
	}
	return &UserResolver{user: &user}, nil
}
