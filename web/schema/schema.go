package schema

import (
	"io/ioutil"

	"github.com/neelance/graphql-go"
	conf "github.com/wkozyra95/go-graphql-starter/config"
	"gopkg.in/mgo.v2/bson"
)

var log = conf.NamedLogger("schema")

// SetupSchema ...
func SetupSchema(resolver *Resolver) (*graphql.Schema, error) {
	schema, readErr := ioutil.ReadFile("schema.gql")
	if readErr != nil {
		return nil, readErr
	}
	return graphql.ParseSchema(
		string(schema),
		resolver,
	)
}

// Resolver ...
type Resolver struct {
	GenerateToken func(id bson.ObjectId) string
	Config        *conf.Config
}
