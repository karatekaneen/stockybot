package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Watch holds the schema definition for the Security entity.
type Watch struct {
	ent.Schema
}

// Fields of the Watch.
func (Watch) Fields() []ent.Field {
	return []ent.Field{
		field.Time("watched_since").Default(time.Now),
		field.String("user_id").NotEmpty(),
	}
}

func (Watch) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}

// Edges of the Watch.
func (Watch) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("security", Security.Type).Ref("watchers").Unique(),
	}
}
