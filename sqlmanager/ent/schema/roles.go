package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Roles holds the schema definition for the Roles entity.
type Roles struct {
	ent.Schema
}

// Fields of the Roles.
func (Roles) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("name"),
		field.String("description"),
		// implement sub roles
	}
}

// Edges of the Roles.
func (Roles) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("policy", PolicyOptions.Type).
			Ref("roles").
			Unique(),
	}
}
