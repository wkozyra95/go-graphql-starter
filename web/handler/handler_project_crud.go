package handler

import (
	"fmt"
	"net/http"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/model/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ProjectGet ...
func ProjectGet(id bson.ObjectId, context ActionContext) (model.Project, error) {
	accessOk, accessErr := validateReadRights(id, context.DB.Project(), context)
	if accessErr == mgo.ErrNotFound {
		return model.Project{}, errors.NewSimple(
			fmt.Sprintf("project %s notfound", id.Hex()),
			http.StatusNotFound,
			errors.ErrNotFound,
		)
	}
	if accessErr != nil {
		return model.Project{}, internalServerErr(
			fmt.Sprintf("project %s access error [%s]", id.Hex(), accessErr.Error()),
		)
	}
	if !accessOk {
		return model.Project{}, errors.NewSimple(
			fmt.Sprintf("project %s access unauthorized", id.Hex()),
			http.StatusBadRequest,
			errors.ErrUnauthorized,
		)
	}
	project := model.Project{}
	getErr := context.DB.Project().FindID(id).One(&project)
	if getErr == mgo.ErrNotFound {
		log.Error("[Assert] unreachable code")
		return model.Project{}, internalServerErr("unreachable code")
	}
	if getErr != nil {
		return model.Project{}, internalServerErr(
			fmt.Sprintf("project %s get error [%s]", id.Hex(), accessErr.Error()),
		)
	}

	return project, nil
}

// ProjectGetAll ...
func ProjectGetAll(context ActionContext) ([]model.Project, error) {
	projects := []model.Project{}
	getErr := context.DB.Project().Find(bson.M{db.UserForeignKey: context.UserID}).All(&projects)
	if getErr != nil {
		return nil, internalServerErr(
			fmt.Sprintf("projects all get error [%s]", getErr.Error()),
		)
	}
	return projects, nil
}

// ProjectCreate ...
func ProjectCreate(project model.Project, context ActionContext) (bson.ObjectId, error) {
	project.ID = bson.NewObjectId()
	project.UserID = context.UserID
	createErr := context.DB.Project().Insert(project)
	if createErr != nil {
		return "", internalServerErr(
			fmt.Sprintf("project %s create error [%s]", project.ID.Hex(), createErr.Error()),
		)
	}
	return project.ID, nil
}

// ProjectUpdate ...
func ProjectUpdate(project model.Project, context ActionContext) error {
	project.UserID = context.UserID
	accessOk, accessErr := validateWriteRights(project.ID, context.DB.Project(), context)
	if accessErr == mgo.ErrNotFound {
		return errors.NewSimple(
			fmt.Sprintf("project %s notfound", project.ID.Hex()),
			http.StatusNotFound,
			errors.ErrNotFound,
		)
	}
	if accessErr != nil {
		return internalServerErr(
			fmt.Sprintf("project %s access error [%s]", project.ID.Hex(), accessErr.Error()),
		)
	}
	if !accessOk {
		return errors.NewSimple(
			fmt.Sprintf("project %s access unauthorized", project.ID.Hex()),
			http.StatusBadRequest,
			errors.ErrUnauthorized,
		)
	}

	updateErr := context.DB.Project().UpdateID(project.ID, project)
	if updateErr == mgo.ErrNotFound {
		log.Error("[Assert] unreachable code")
		return internalServerErr("unreachable code")
	}
	if updateErr != nil {
		return internalServerErr(
			fmt.Sprintf("project %s update error [%s]", project.ID.Hex(), updateErr.Error()),
		)
	}

	return nil
}

// ProjectDelete ...
func ProjectDelete(id bson.ObjectId, context ActionContext) error {
	accessOk, accessErr := validateWriteRights(id, context.DB.Project(), context)
	if accessErr == mgo.ErrNotFound {
		return errors.NewSimple(
			fmt.Sprintf("project %s notfound", id.Hex()),
			http.StatusNotFound,
			errors.ErrNotFound,
		)
	}
	if accessErr != nil {
		return internalServerErr(
			fmt.Sprintf("project %s access error [%s]", id.Hex(), accessErr.Error()),
		)
	}
	if !accessOk {
		return errors.NewSimple(
			fmt.Sprintf("project %s access unauthorized", id.Hex()),
			http.StatusBadRequest,
			errors.ErrUnauthorized,
		)
	}
	removeErr := context.DB.Project().RemoveID(id)
	if removeErr == mgo.ErrNotFound {
		log.Error("[Assert] unreachable code")
		return internalServerErr("unreachable code")
	}
	if removeErr != nil {
		return internalServerErr(
			fmt.Sprintf("project %s update error [%s]", id.Hex(), removeErr.Error()),
		)
	}

	return nil
}
