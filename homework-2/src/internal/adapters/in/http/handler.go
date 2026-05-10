package handler

import (
	"io"
	"strings"
	"support-tickets/internal/domain"
	"support-tickets/internal/ports/in"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service in.TicketService
}

func New(service in.TicketService) *Handler {
	return &Handler{service: service}
}

type CreateTicketRequest struct {
	CustomerID    string                `json:"customer_id" form:"customer_id"`
	CustomerEmail string                `json:"customer_email" form:"customer_email"`
	CustomerName  string                `json:"customer_name" form:"customer_name"`
	Subject       string                `json:"subject" form:"subject"`
	Description   string                `json:"description" form:"description"`
	Category      string                `json:"category" form:"category"`
	Priority      string                `json:"priority" form:"priority"`
	Metadata      *TicketMetadataInput  `json:"metadata,omitempty" form:"metadata,omitempty"`
}

type TicketMetadataInput struct {
	Source     string `json:"source"`
	Browser    string `json:"browser"`
	DeviceType string `json:"device_type"`
}

// CreateTicket godoc
// @Summary Create a new support ticket
// @Description Create a new support ticket with optional auto-classification
// @Tags tickets
// @Accept json
// @Produce json
// @Param request body CreateTicketRequest true "Ticket request"
// @Param auto_classify query boolean false "Auto-classify ticket"
// @Success 201 {object} domain.Ticket
// @Failure 400 {object} errorhandler.ErrorResponse
// @Router /tickets [post]
func (h *Handler) CreateTicket(c *fiber.Ctx) error {
	var req CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	metadata := domain.TicketMetadata{}
	if req.Metadata != nil {
		metadata = domain.TicketMetadata{
			Source:     domain.Source(req.Metadata.Source),
			Browser:    req.Metadata.Browser,
			DeviceType: domain.DeviceType(req.Metadata.DeviceType),
		}
	}

	ticket, err := h.service.Create(
		req.CustomerID,
		req.CustomerEmail,
		req.CustomerName,
		req.Subject,
		req.Description,
		domain.Category(req.Category),
		domain.Priority(req.Priority),
		metadata,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(ticket)
}

// BulkImport godoc
// @Summary Bulk import tickets from file
// @Description Import tickets from CSV, JSON, or XML file
// @Tags tickets
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to import (CSV, JSON, or XML)"
// @Success 200 {object} in.BulkImportResult
// @Failure 400 {object} errorhandler.ErrorResponse
// @Router /tickets/import [post]
func (h *Handler) BulkImport(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return &domain.ErrInvalidInput{Message: "file is required"}
	}

	src, err := file.Open()
	if err != nil {
		return &domain.ErrInvalidInput{Message: "failed to open file"}
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return &domain.ErrInvalidInput{Message: "failed to read file"}
	}

	// Detect format
	format := "json"
	if strings.HasSuffix(file.Filename, ".csv") {
		format = "csv"
	} else if strings.HasSuffix(file.Filename, ".xml") {
		format = "xml"
	}

	contentType := c.FormValue("format")
	if contentType != "" {
		format = contentType
	}

	result, err := h.service.BulkImport(format, data)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// ListTickets godoc
// @Summary List all tickets
// @Description List tickets with optional filtering
// @Tags tickets
// @Produce json
// @Param status query string false "Filter by status"
// @Param category query string false "Filter by category"
// @Param priority query string false "Filter by priority"
// @Param customer_id query string false "Filter by customer_id"
// @Success 200 {array} domain.Ticket
// @Router /tickets [get]
func (h *Handler) ListTickets(c *fiber.Ctx) error {
	req := &in.TicketListRequest{
		Status:     getQueryParam(c, "status"),
		Category:   getQueryParam(c, "category"),
		Priority:   getQueryParam(c, "priority"),
		CustomerID: getQueryParam(c, "customer_id"),
	}

	tickets, err := h.service.List(req)
	if err != nil {
		return err
	}

	if tickets == nil {
		tickets = []*domain.Ticket{}
	}

	return c.JSON(tickets)
}

// GetTicket godoc
// @Summary Get a ticket by ID
// @Description Get a specific ticket
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} domain.Ticket
// @Failure 404 {object} errorhandler.ErrorResponse
// @Router /tickets/{id} [get]
func (h *Handler) GetTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	ticket, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return c.JSON(ticket)
}

type UpdateTicketRequest struct {
	Status     *string   `json:"status,omitempty"`
	Category   *string   `json:"category,omitempty"`
	Priority   *string   `json:"priority,omitempty"`
	AssignedTo *string   `json:"assigned_to,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
}

// UpdateTicket godoc
// @Summary Update a ticket
// @Description Update ticket fields
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param request body UpdateTicketRequest true "Update request"
// @Success 200 {object} domain.Ticket
// @Failure 400 {object} errorhandler.ErrorResponse
// @Failure 404 {object} errorhandler.ErrorResponse
// @Router /tickets/{id} [put]
func (h *Handler) UpdateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return &domain.ErrInvalidInput{Message: "Invalid request body"}
	}

	updates := make(map[string]interface{})
	if req.Status != nil {
		updates["status"] = domain.Status(*req.Status)
	}
	if req.Category != nil {
		updates["category"] = domain.Category(*req.Category)
	}
	if req.Priority != nil {
		updates["priority"] = domain.Priority(*req.Priority)
	}
	if req.AssignedTo != nil {
		updates["assigned_to"] = *req.AssignedTo
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}

	ticket, err := h.service.Update(id, updates)
	if err != nil {
		return err
	}

	return c.JSON(ticket)
}

// DeleteTicket godoc
// @Summary Delete a ticket
// @Description Delete a specific ticket
// @Tags tickets
// @Param id path string true "Ticket ID"
// @Success 204
// @Failure 404 {object} errorhandler.ErrorResponse
// @Router /tickets/{id} [delete]
func (h *Handler) DeleteTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AutoClassify godoc
// @Summary Auto-classify a ticket
// @Description Run auto-classification on a specific ticket
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} domain.Classification
// @Failure 404 {object} errorhandler.ErrorResponse
// @Router /tickets/{id}/auto-classify [post]
func (h *Handler) AutoClassify(c *fiber.Ctx) error {
	id := c.Params("id")
	classification, err := h.service.AutoClassify(id)
	if err != nil {
		return err
	}

	return c.JSON(classification)
}

func getQueryParam(c *fiber.Ctx, key string) *string {
	val := c.Query(key)
	if val == "" {
		return nil
	}
	return &val
}
