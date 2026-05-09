package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Aircraft struct {
	ent.Schema
}

func (Aircraft) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MinLen(5).
			MaxLen(20).
			Immutable(),
		field.Strings("seats"),
	}
}

func (Aircraft) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("flights", Flight.Type),
	}
}
