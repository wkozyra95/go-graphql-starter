package schema

import (
	"context"

	"github.com/wkozyra95/go-graphql-starter/model"
)

func (r Resolver) Project(
	context context.Context,
	args struct {
		ID string
	},
) *projectResolver {
	return &projectResolver{Resolver: r, project: model.Project{}}
}
