package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/model"
)

func (r Resolver) User(
	context context.Context,
	args struct {
		Username string
	},
) (*userResolver, error) {
	return &userResolver{Resolver: r, user: model.User{}}, nil
}
