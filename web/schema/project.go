package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
)

type projectResolver struct {
	project *model.Project
	errors  *model.ProjectError
}

func (pr *projectResolver) ID() string {
	return pr.project.ID.Hex()
}

func (pr *projectResolver) Name() (string, error) {
	if pr.errors != nil && pr.errors.Name != nil {
		return "", pr.errors.Name
	}
	if pr.project == nil {
		return "", errors.InternalServerError
	}
	return pr.project.Name, nil
}

func (pr *projectResolver) Description() (string, error) {
	if pr.errors != nil && pr.errors.Description != nil {
		return "", pr.errors.Description
	}
	if pr.project == nil {
		return "", errors.InternalServerError
	}
	return pr.project.Description, nil
}

func (pr *projectResolver) User(ctx context.Context) (*userResolver, error) {
	db := extractDBSession(ctx)
	userId := extractUserIdContext(ctx)

	user := model.User{}
	dbErr := db.User().FindID(pr.project.UserID).One(&user)
	if dbErr != nil {
		return nil, errors.InternalServerError
	}
	if user.ID != userId {
		return nil, errors.Unauthorized
	}
	return &userResolver{user: &user}, nil
}
