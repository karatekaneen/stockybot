package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Security holds the schema definition for the Security entity.
type Security struct {
	ent.Schema
}

// Fields of the Security.
func (Security) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive().Unique(),
		field.String("name").NotEmpty().Unique(),
		field.String("country").NotEmpty(),
		field.String("link_name").Optional(),
		field.String("list").NotEmpty(),
		field.Enum("type").Values("stock", "index"),
	}
}

func (Security) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Edges of the Security.
func (Security) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("watchers", Watch.Type),
	}
}
