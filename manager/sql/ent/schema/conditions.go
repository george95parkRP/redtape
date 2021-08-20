package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Conditions holds the schema definition for the Conditions entity.
type Conditions struct {
	ent.Schema
}

// Fields of the Conditions.
func (Conditions) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("type"),
		field.JSON("options", map[string]interface{}{}),
	}
}

// Edges of the Conditions.
func (Conditions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subjects", Subjects.Type).
			Ref("conditions").
			Unique(),
	}
}
