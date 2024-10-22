package service

import (
	"tickets_api/internal/models"
	"tickets_api/internal/repository"
)

type TicketService struct {
	repo *repository.TicketRepository
}

func NewTicketServiceDefault(dsn string) *TicketService {
	ticketRepo := repository.NewTicketRepository(dsn)
	return &TicketService{repo: ticketRepo}
}

func NewTicketService(ticketRepo *repository.TicketRepository) *TicketService {
	return &TicketService{repo: ticketRepo}
}

func (service *TicketService) CloseTicketService() {
	service.repo.CloseTicketRepository()
}

func (service *TicketService) CreateTicket(name, description string, allocation int) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Name:        name,
		Description: description,
		Allocation:  allocation,
		Remaining:   allocation,
		//CreatedAt:   time.Now(),
	}

	if err := service.repo.CreateTicket(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (service *TicketService) GetTicket(id int) (*models.Ticket, error) {
	return service.repo.GetTicket(id)
}

func (service *TicketService) PurchaseTicket(id int, purchaseRequest *models.PurchaseRequest) error {

	if err := service.repo.PurchaseTicket(id, purchaseRequest); err != nil {
		return err
	}

	return nil
}
