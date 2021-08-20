package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subjects holds the schema definition for the Subjects entity.
type Subjects struct {
	ent.Schema
}

// Fields of the Subjects.
func (Subjects) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
	}
}

// Edges of the Subjects.
func (Subjects) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("policy", PolicyOptions.Type).
			Ref("subjects").
			Unique(),
		edge.To("conditions", Conditions.Type),
	}
}
