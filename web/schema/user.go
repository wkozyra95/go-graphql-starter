package schema

import "github.com/wkozyra95/go-graphql-starter/model"

type userResolver struct {
	Resolver
	user model.User
}

func (ur *userResolver) ID() string {
	return ur.user.ID.Hex()
}

func (ur *userResolver) Username() string {
	return ur.user.Username
}

func (ur *userResolver) Email() string {
	return ur.user.Email
}

func (ur *userResolver) Projects() []*projectResolver {
	return []*projectResolver{
		&projectResolver{Resolver: ur.Resolver, project: model.Project{}},
		&projectResolver{Resolver: ur.Resolver, project: model.Project{}},
		&projectResolver{Resolver: ur.Resolver, project: model.Project{}},
		&projectResolver{Resolver: ur.Resolver, project: model.Project{}},
	}
}
