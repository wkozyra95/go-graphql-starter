package schema

import (
	"context"
	"fmt"

	"github.com/wkozyra95/go-graphql-starter/errors"
	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/model/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (r Resolver) AuthLogin(
	context context.Context,
	args struct {
		LoginForm *loginForm
	},
) (*loginResponseResolver, error) {
	db := extractDBSession(context)

	if err := args.LoginForm.validate(); err != nil {
		return nil, err
	}

	user := mongo.User{}
	userErr := db.User().Find(bson.M{
		mongo.UserIDKeyUsername: args.LoginForm.Username,
	}).One(&user)
	if userErr == mgo.ErrNotFound {
		return nil, fmt.Errorf("Unknown combination of username and password")
	}
	if userErr != nil {
		return nil, errors.InternalServerError
	}
	if err := args.LoginForm.validatePassword(user.PasswordHash); err != nil {
		return nil, err
	}

	token := r.GenerateToken(user.ID)
	return &loginResponseResolver{
		token: token,
		user:  &user,
	}, nil
}

type loginResponseResolver struct {
	token string
	user  *mongo.User
}

func (lr *loginResponseResolver) Token() string {
	return lr.token
}

func (lr *loginResponseResolver) User() *userResolver {
	return &userResolver{user: lr.user}
}

type loginForm struct {
	Username string
	Password string
}

func (lf loginForm) validate() error {
	if lf.Username == "" {
		return fmt.Errorf("Username can't be empty")
	}
	if lf.Password == "" {
		return fmt.Errorf("Password can't be empty")
	}
	if len(lf.Password) < 8 {
		return fmt.Errorf("Password is to short, you need at least 8 characters")
	}
	return nil
}

func (lf loginForm) validatePassword(hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(lf.Password))
	if err != nil {
		return fmt.Errorf("Unknown combination of user and password.")
	}
	return nil
}

func (r Resolver) AuthRegister(
	context context.Context,
	args struct {
		RegisterForm *registerForm
	},
) (*userResolver, error) {
	db := extractDBSession(context)

	if err := args.RegisterForm.validate(); err != nil {
		return nil, err
	}

	count, countErr := db.User().
		Find(bson.M{"username": args.RegisterForm.Username}).Count()
	if countErr != nil {
		return nil, errors.InternalServerError
	}
	if count > 0 {
		return nil, fmt.Errorf("User with that username already exists")
	}

	user := args.RegisterForm.createUser()
	insertErr := db.User().Insert(user)
	if insertErr != nil {
		return nil, errors.InternalServerError
	}
	return &userResolver{user: &user}, nil
}

type registerForm struct {
	Email    string
	Username string
	Password string
}

func (rf registerForm) validate() error {
	if rf.Username == "" {
		return fmt.Errorf("Username can't be empty")
	}
	if rf.Password == "" {
		return fmt.Errorf("Password can't be empty")
	}
	if rf.Email == "" {
		return fmt.Errorf("Email can't be empty")
	}
	if len(rf.Password) < 8 {
		return fmt.Errorf("Password is to short, you need at least 8 characters")
	}
	return nil
}

func (rf registerForm) createUser() mongo.User {
	return mongo.User{
		ID: bson.NewObjectId(),
		User: model.User{
			Username:     rf.Username,
			Email:        rf.Email,
			PasswordHash: rf.generatePasswordHash(),
		},
	}
}

func (rf registerForm) generatePasswordHash() string {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(rf.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error("[ASSERT] Failed to hash password")
		panic(err)
	}
	return string(hashedPassword)
}
