package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"gopkg.in/mgo.v2"
)

func (r *Resolver) ProjectCreate(
	context context.Context,
	args struct {
		ProjectInput *model.ProjectCreateInput
	},
) (*projectResolver, error) {
	userId := extractUserIdContext(context)
	db := extractDBSession(context)

	if specyficError, globalErr := args.ProjectInput.Validate(); specyficError != nil && globalErr != nil {
		return &projectResolver{errors: specyficError}, globalErr
	}
	project := args.ProjectInput.CreateProject(userId)
	err := db.Project().Insert(project)
	if err != nil {
		return nil, errors.InternalServerError
	}

	return &projectResolver{project: project}, nil
}

func (r Resolver) ProjectUpdate(
	context context.Context,
	args struct {
		ProjectUpdate *model.ProjectUpdateInput
	},
) (*projectResolver, error) {
	userId := extractUserIdContext(context)
	db := extractDBSession(context)

	query, change := args.ProjectUpdate.UpdateProject(userId)
	err := db.Project().Update(query, change)
	if err == mgo.ErrNotFound {
		return nil, errors.NotFound
	}
	if err != nil {
		return nil, errors.InternalServerError
	}

	project := model.Project{}
	getErr := db.Project().FindID(args.ProjectUpdate.ID).One(&project)
	if getErr == mgo.ErrNotFound {
		return nil, errors.NotFound
	}
	if getErr != nil {
		return nil, errors.InternalServerError
	}

	return &projectResolver{project: &project}, nil

}
