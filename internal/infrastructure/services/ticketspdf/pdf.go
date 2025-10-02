package ticketspdf

import "github.com/allexborysov/aircraft/internal/domain/flight"

type TicketsPDFPresetner struct {
}

func New() *TicketsPDFPresetner {
	return &TicketsPDFPresetner{}
}

func (p *TicketsPDFPresetner) GeneratePDF(ticket *flight.Ticket) (string, error) {
	return "not_implemented", nil
}
