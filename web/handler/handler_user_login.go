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

// UserLogin ...
func UserLogin(
	username string,
	password string,
	DB db.DB,
) (model.User, error) {
	log.Info("login resolver")
	validateErr := userLoginValidate(username, password)
	log.Info("login resolver")
	if validateErr != nil {
		return model.User{}, validateErr
	}
	log.Info("login resolver")

	formErr := errors.Empty()
	formErr.Code = http.StatusBadRequest
	formErr.JSON["form"] = errors.TextError("Unknown combiantion of username and password")
	formErr.JSON["reason"] = errors.TextError(errors.ErrFormError)
	log.Info("login resolver")

	user := model.User{}
	userErr := DB.User().Find(bson.M{db.UserIDKeyUsername: username}).One(&user)
	if userErr == mgo.ErrNotFound {
		formErr.Msg = fmt.Sprintf("user (%s) not found", username)
		return user, formErr
	}
	log.Info("login resolver")
	if userErr != nil {
		return user, internalServerErr(
			fmt.Sprintf("user (%s) find error [%s]", username, userErr.Error()),
		)
	}
	log.Info("login resolver")

	if !user.ValidatePassword(password) {
		formErr.Msg = "invalid password"
		return user, formErr
	}
	return user, nil
}

func userLoginValidate(username, password string) error {
	formErr := errors.New("form error", http.StatusBadRequest)

	if username == "" {
		formErr.JSON["username"] = errors.TextError("Username can't empty")
	}
	if password == "" {
		formErr.JSON["password"] = errors.TextError("Password can't empty")
	}
	if len(formErr.JSON) > 0 {
		return formErr
	}
	return nil
}
