package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/model"
	"github.com/wkozyra95/go-graphql-starter/web/handler"
)

func (r Resolver) AuthLogin(
	context context.Context,
	args struct {
		LoginForm *loginForm
	},
) (*loginResponseResolver, error) {
	db := extractDBSession(context)
	user, userErr := handler.UserLogin(
		args.LoginForm.Username,
		args.LoginForm.Password,
		db,
	)
	if userErr != nil {
		return &loginResponseResolver{}, userErr
	}
	token := r.GenerateToken(user.ID)
	return &loginResponseResolver{
		Resolver: r,
		token:    token,
		user:     user,
	}, nil
}

func (r Resolver) AuthRegister(
	context context.Context,
	args struct {
		RegisterForm *registerForm
	},
) (bool, error) {
	db := extractDBSession(context)
	err := handler.UserRegister(
		model.User{
			Email:    args.RegisterForm.Email,
			Username: args.RegisterForm.Username,
		},
		args.RegisterForm.Password,
		db,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

type loginResponseResolver struct {
	Resolver
	token string
	user  model.User
}

func (lr *loginResponseResolver) Token() string {
	return lr.token
}

func (lr *loginResponseResolver) User() *userResolver {
	return &userResolver{Resolver: lr.Resolver, user: lr.user}
}

type loginForm struct {
	Username string
	Password string
}

type registerForm struct {
	Email    string
	Username string
	Password string
}
