package schema

import (
	"context"
	"fmt"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ProjectCreate ...
func (r *Resolver) ProjectCreate(
	context context.Context,
	args struct {
		Project *projectCreateInput
	},
) (*ProjectResolver, error) {
	userID := extractUserID(context)
	db := extractDBSession(context)

	if err := args.Project.validate(); err != nil {
		return nil, err
	}
	project := args.Project.createProject(userID)
	err := db.Project().Insert(project)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	return &ProjectResolver{project: &project}, nil
}

type projectCreateInput model.Project

func (pc projectCreateInput) validate() error {
	return nil
}

func (pc projectCreateInput) createProject(userID bson.ObjectId) mongo.Project {
	return mongo.Project{
		ID:      bson.NewObjectId(),
		UserID:  userID,
		Project: model.Project(pc),
	}
}

// ProjectUpdate ...
func (r *Resolver) ProjectUpdate(
	ctx context.Context,
	args struct {
		Project *projectUpdateInput
	},
) (*ProjectResolver, error) {
	userID := extractUserID(ctx)
	db := extractDBSession(ctx)
	projectID, projectIDErr := mongo.ConvertToObjectID(args.Project.ID)
	if projectIDErr != nil {
		return nil, errors.ErrInternalServerError
	}

	project := mongo.Project{}
	getErr := db.Project().FindID(projectID).One(&project)
	if getErr == mgo.ErrNotFound {
		return nil, errors.ErrNotFound
	}
	if getErr != nil {
		return nil, errors.ErrInternalServerError
	}
	if project.UserID != userID {
		return nil, errors.ErrUnauthorized
	}
	args.Project.updateProject(&project)
	updateErr := db.Project().UpdateID(projectID, project)
	if updateErr != nil {
		return nil, updateErr
	}

	return &ProjectResolver{project: &project}, nil
}

type projectUpdateInput struct {
	ID          string
	Name        string
	Description string
	Details     model.ProjectDetails
}

func (pu projectUpdateInput) updateProject(project *mongo.Project) {
	project.Name = pu.Name
	project.Description = pu.Description
	project.Details = pu.Details
}

// ProjectPatch ...
func (r *Resolver) ProjectPatch(
	ctx context.Context,
	args struct {
		Project *projectPatchInput
	},
) (*ProjectResolver, error) {
	userID := extractUserID(ctx)
	db := extractDBSession(ctx)
	projectID, projectIDErr := mongo.ConvertToObjectID(args.Project.ID)
	if projectIDErr != nil {
		return nil, errors.ErrInternalServerError
	}

	project := mongo.Project{}
	getErr := db.Project().FindID(projectID).One(&project)
	if getErr == mgo.ErrNotFound {
		return nil, errors.ErrNotFound
	}
	if getErr != nil {
		return nil, errors.ErrInternalServerError
	}
	if project.UserID != userID {
		return nil, errors.ErrUnauthorized
	}

	applyErr := args.Project.apply(&project.Project)
	if applyErr != nil {
		return nil, applyErr
	}
	updateErr := db.Project().UpdateID(projectID, project)
	if updateErr != nil {
		return nil, updateErr
	}

	return &ProjectResolver{project: &project}, nil
}

type projectPatchInput struct {
	ID    string
	Patch []updatePatch
}

func (pi projectPatchInput) apply(project *model.Project) error {
	for _, patch := range pi.Patch {
		if value, ok := patch.Value.string(); patch.Field == "name" && ok {
			project.Name = value
		} else if value, ok := patch.Value.string(); patch.Field == "description" && ok {
			project.Description = value
		} else if value, ok := patch.Value.bool(); patch.Field == "details.isPublic" && ok {
			project.Details.IsPublic = value
		} else if value, ok := patch.Value.string(); patch.Field == "details.projectType" && ok {
			project.Details.ProjectType = value
		} else {
			return fmt.Errorf("invalid update %s: %+v", patch.Field, patch.Value.value)
		}
	}
	return nil
}

// ProjectDelete ...
func (r *Resolver) ProjectDelete(
	ctx context.Context,
	args struct {
		ProjectID string
	},
) (bool, error) {
	userID := extractUserID(ctx)
	db := extractDBSession(ctx)
	projectID, projectIDErr := mongo.ConvertToObjectID(args.ProjectID)
	if projectIDErr != nil {
		return false, errors.ErrInternalServerError
	}

	project := mongo.Project{}
	getErr := db.Project().FindID(projectID).One(&project)
	if getErr == mgo.ErrNotFound {
		return false, errors.ErrNotFound
	}
	if getErr != nil {
		return false, errors.ErrInternalServerError
	}
	if project.UserID != userID {
		return false, errors.ErrUnauthorized
	}

	if err := db.Project().RemoveID(projectID); err != nil {
		return false, errors.ErrInternalServerError
	}

	return true, nil
}
