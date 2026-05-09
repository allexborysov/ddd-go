package ticketspdf

import "github.com/allexborysov/aircraft/internal/domain/flight"

type TicketsPDFPresenter struct {
}

func New() *TicketsPDFPresenter {
	return &TicketsPDFPresenter{}
}

func (p *TicketsPDFPresenter) GeneratePDF(ticket *flight.Ticket) (string, error) {
	return "not_implemented", nil
}
