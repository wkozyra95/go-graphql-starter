package model

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Project ...
type Project struct {
	// ID ...
	ID bson.ObjectId `json:"id" bson:"_id"`
	// UserID ...
	UserID bson.ObjectId `json:"userId" bson:"userId"`

	// Name ...
	Name string `json:"name" bson:"name"`

	// Description ...
	Description string `json:"description" bson:"description"`
}

type ProjectCreateInput struct {
	Name        string
	Description string
}

func (pi ProjectCreateInput) CreateProject(userId bson.ObjectId) *Project {
	return &Project{
		ID:          bson.NewObjectId(),
		UserID:      userId,
		Name:        pi.Name,
		Description: pi.Description,
	}
}

func (pi ProjectCreateInput) Validate() (*ProjectError, error) {
	detailedErr := &ProjectError{}
	if pi.Name == "" {
		detailedErr.Name = fmt.Errorf("Name can't be empty")
	}
	if pi.Description == "" {
		detailedErr.Description = fmt.Errorf("Description can't be empty")
	}
	if detailedErr.HasError() {
		return detailedErr, nil
	}
	return nil, nil
}

type ProjectError struct {
	Name        error
	Description error
}

func (pe ProjectError) HasError() bool {
	return pe.Name != nil || pe.Description != nil
}

type projectUpdateInputSetField struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type projectUpdateInputRemoveField struct {
	Name        bool `json:"name"`
	Description bool `json:"description"`
}

type ProjectUpdateInput struct {
	ID          bson.ObjectId                 `json:"id"`
	UpdateField projectUpdateInputSetField    `json:"set"`
	DeleteField projectUpdateInputRemoveField `json:"remove"`
}

func (pu ProjectUpdateInput) isUpdating() bool {
	return pu.UpdateField.Name != "" || pu.UpdateField.Description != ""
}

func (pu ProjectUpdateInput) isDeleting() bool {
	return pu.DeleteField.Name || pu.DeleteField.Description
}

func (pi ProjectUpdateInput) UpdateProject(userId bson.ObjectId) (query bson.M, change bson.M) {
	query["_id"] = pi.ID
	query["userId"] = userId
	if pi.isUpdating() {
		changeSet := bson.M{}
		if name := pi.UpdateField.Name; name != "" {
			changeSet["name"] = name
		}
		if description := pi.UpdateField.Description; description != "" {
			changeSet["description"] = description
		}
		change["$set"] = changeSet
	}

	if pi.isDeleting() {
		deleteSet := bson.M{}
		if pi.DeleteField.Name {
			deleteSet["name"] = ""
		}
		if pi.DeleteField.Description {
			deleteSet["description"] = ""
		}
		change["$unset"] = deleteSet
	}
	return
}
