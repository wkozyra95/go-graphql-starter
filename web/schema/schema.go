package schema

import (
	"github.com/neelance/graphql-go"
	conf "github.com/wkozyra95/go-graphql-starter/config"
	"gopkg.in/mgo.v2/bson"
)

var log = conf.NamedLogger("schema")

func SetupSchema(resolver *Resolver) (*graphql.Schema, error) {
	return graphql.ParseSchema(
		schema,
		resolver,
	)
}

var schema = `
	schema {
		query: Query
		mutation: Mutation
	}
	# Query
	type Query {
		user(username: String!): User
		project(id: String!): Project
	}
	# Mutation
	type Mutation {
		authLogin(loginForm: LoginForm!): LoginResponse!
		authRegister(registerForm: RegisterForm!): Boolean!
	}
	# User
	type User {
		id: String!
		username: String!
		email: String!
		projects: [Project!]!
	}
	# Project
	type Project {
		id: String!
		name: String!
		description: String!
		user: User!
	}
	# LoginForm
	input LoginForm {
		username: String!
		password: String!
	}
	# LoginResponse
	type LoginResponse {
		token: String!
		user: User!
	}
	# RegisterForm
	input RegisterForm {
		email: String!,
		username: String!,
		password: String!,
	}
`

type Resolver struct {
	GenerateToken func(id bson.ObjectId) string
	Config        *conf.Config
}
