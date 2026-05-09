package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type FlightSeat struct {
	ent.Schema
}

func (FlightSeat) Fields() []ent.Field {
	return []ent.Field{
		field.String("flight_id"),
		field.String("seat_number").NotEmpty(),
		field.Float("price"),
		field.String("passenger_id").Optional(),
	}
}

func (FlightSeat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("flight", Flight.Type).
			Ref("seats").
			Field("flight_id").
			Unique().
			Required(),
	}
}

func (FlightSeat) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("flight_id", "seat_number").Unique(),
	}
}
