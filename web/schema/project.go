package schema

import "github.com/wkozyra95/go-graphql-starter/model"

type projectResolver struct {
	Resolver
	project model.Project
}

func (pr *projectResolver) ID() string {
	return pr.project.ID.Hex()
}

func (pr *projectResolver) Name() string {
	return pr.project.Name
}

func (pr *projectResolver) Description() string {
	return pr.project.Description
}

func (pr *projectResolver) User() *userResolver {
	return &userResolver{Resolver: pr.Resolver, user: model.User{}}
}
