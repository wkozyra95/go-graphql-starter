package handler

import (
	"fmt"
	"net/http"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/model/db"
	"gopkg.in/mgo.v2/bson"
)

// UserRegister ...
func UserRegister(user model.User, password string, DB db.DB) error {
	validateErr := userRegisterValidate(user, password)
	if validateErr != nil {
		return validateErr
	}
	existingUsers := []model.User{}
	getErr := DB.User().Find(bson.M{
		db.OR: []bson.M{
			bson.M{db.UserIDKeyUsername: user.Username},
			bson.M{db.UserIDKeyEmail: user.Email},
		},
	}).All(&existingUsers)
	if getErr != nil {
		return internalServerErr(
			fmt.Sprintf("users find error [%s]", getErr.Error()),
		)
	}
	for _, existingUser := range existingUsers {
		if existingUser.Email == user.Email {
			matchErr := errors.New("email already used", http.StatusBadRequest)
			matchErr.JSON["email"] = errors.TextError("This email is already in use")
			return matchErr
		}
		if existingUser.Username == user.Username {
			matchErr := errors.New("username already used", http.StatusBadRequest)
			matchErr.JSON["username"] = errors.TextError("This username is already in use")
			return matchErr
		}
	}

	user.ID = bson.NewObjectId()
	user.GeneratePasswordHash(password)
	insertErr := DB.User().Insert(user)
	if insertErr != nil {
		return internalServerErr(
			fmt.Sprintf("user %s insert error [%s]", user.ID.Hex(), insertErr.Error()),
		)
	}
	return nil
}

func userRegisterValidate(user model.User, password string) error {
	formErr := errors.New("form error", http.StatusBadRequest)

	if user.Email == "" {
		formErr.JSON["email"] = errors.TextError("Email can't empty")
	}

	if user.Username == "" {
		formErr.JSON["username"] = errors.TextError("Username can't empty")
	}

	if len(password) < 8 {
		formErr.JSON["password"] = errors.TextError("Password is to short, you need at least 8 characters.")
	}

	if password == "" {
		formErr.JSON["password"] = errors.TextError("Password can't empty")
	}

	if len(formErr.JSON) > 0 {
		return formErr
	}
	return nil
}
