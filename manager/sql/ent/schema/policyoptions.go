package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PolicyOptions holds the schema definition for the PolicyOptions entity.
type PolicyOptions struct {
	ent.Schema
}

// Fields of the PolicyOptions.
func (PolicyOptions) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("name"),
		field.String("description"),
		field.Strings("resources"),
		field.Strings("actions"),
		field.Strings("scopes"),
		field.String("effect"),
	}
}

// Edges of the PolicyOptions.
func (PolicyOptions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subjects", Subjects.Type),
	}
}
