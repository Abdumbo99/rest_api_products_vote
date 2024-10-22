package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"tickets_api/internal/service"

	"tickets_api/internal/models"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandlerDefault(dsn string) *TicketHandler {
	ticketService := service.NewTicketServiceDefault(dsn)
	return &TicketHandler{service: ticketService}
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{service: ticketService}
}

func (handler *TicketHandler) CloseTicketHandler() {
	handler.service.CloseTicketService()
}

// @Summary Create a new ticket
// @Description Create a new ticket with a specified name, description, and allocation.
// @Tags tickets
// @Accept  json
// @Produce  json
// @Param   ticket  body  Ticket  true  "Ticket to create"
// @Success 200 {object} Ticket
// @Failure 400 {object} error
// @Router /tickets [post]
func (handler *TicketHandler) CreateTicket(c *gin.Context) {

	input := &struct {
		Name        string `json:"name"`
		Description string `json:"desc"`
		Allocation  int    `json:"allocation"`
	}{}

	err := c.ShouldBindJSON(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTicket, err := handler.service.CreateTicket(input.Name, input.Description, input.Allocation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, newTicket)
}

// @Summary Get a ticket by ID
// @Description Get a ticket by its ID.
// @Tags tickets
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "Ticket ID"
// @Success 200 {object} Ticket
// @Failure 400 {object} error
// @Router /tickets/{id} [get]
func (handler *TicketHandler) GetTicket(c *gin.Context) {
	ticketID := c.Param("id")
	id, err := strconv.Atoi(ticketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := handler.service.GetTicket(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ticket"})
		return
	}
	c.IndentedJSON(http.StatusOK, ticket)
}

// @Summary Purchase a ticket
// @Description Purchase a specified quantity of tickets.
// @Tags tickets
// @Accept  json
// @Produce  plain
// @Param   id      path  int           true  "Ticket ID"
// @Param   purchase body  PurchaseRequest  true  "Purchase request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} error
// @Router /tickets/{id}/purchases [post]
func (handler *TicketHandler) PurchaseTicket(c *gin.Context) {

	purchaseRequest := &models.PurchaseRequest{}
	if err := c.ShouldBindJSON(&purchaseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID := c.Param("id")
	id, err := strconv.Atoi(ticketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	err = handler.service.PurchaseTicket(id, purchaseRequest)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
