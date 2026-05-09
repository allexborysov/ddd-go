package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Ticket struct {
	ent.Schema
}

func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("passenger_id").NotEmpty(),
		field.String("seat").NotEmpty(),
		field.Float("price"),
	}
}

func (Ticket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("flight", Flight.Type).
			Ref("tickets").
			Unique().
			Required(),
	}
}
