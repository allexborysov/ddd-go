package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Flight struct {
	ent.Schema
}

func (Flight) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("origin").NotEmpty(),
		field.String("destination").NotEmpty(),
		field.Time("scheduled_departure"),
		field.Time("scheduled_arrival"),
		field.Int64("close_booking_buffer_ns"),
	}
}

func (Flight) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("aircraft", Aircraft.Type).
			Ref("flights").
			Unique().
			Required(),
		edge.To("seats", FlightSeat.Type),
		edge.To("tickets", Ticket.Type),
	}
}
