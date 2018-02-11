package model

// Project ...
type Project struct {
	// Name ...
	Name string `json:"name" bson:"name"`

	// Description ...
	Description string `json:"description" bson:"description"`

	// Details ...
	Details ProjectDetails `json:"details" bson:"details"`
}

type ProjectDetails struct {
	// IsPublic ...
	IsPublic bool `json:"isPublic" bson:"isPublic"`
	// ProjectType ...
	ProjectType string `json:"projectType" bson:"projectType"`
}
