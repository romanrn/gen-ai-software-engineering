package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
)

type TransactionHandler struct {
	service portin.TransactionService
}

func NewTransactionHandler(service portin.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

type CreateTransactionRequest struct {
	FromAccount string  `json:"fromAccount"`
	ToAccount   string  `json:"toAccount"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Type        string  `json:"type"`
}

// Create godoc
// @Summary Create a transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body CreateTransactionRequest true "Transaction payload"
// @Success 201 {object} domain.Transaction
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ValidationErrorResponse
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	tx, err := h.service.CreateTransaction(c.UserContext(), portin.CreateTransactionInput{
		FromAccount: req.FromAccount,
		ToAccount:   req.ToAccount,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        req.Type,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(tx)
}

// List godoc
// @Summary List transactions
// @Tags Transactions
// @Produce json
// @Param accountId query string false "Filter by account"
// @Param type query string false "Filter by type" Enums(deposit,withdrawal,transfer)
// @Param from query string false "Start date YYYY-MM-DD"
// @Param to query string false "End date YYYY-MM-DD"
// @Success 200 {array} domain.Transaction
// @Failure 400 {object} ErrorResponse
// @Router /transactions [get]
func (h *TransactionHandler) List(c *fiber.Ctx) error {
	filter := portin.TransactionFilter{
		AccountID: c.Query("accountId"),
		Type:      c.Query("type"),
	}

	if from := c.Query("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid 'from' date, expected YYYY-MM-DD")
		}
		filter.From = &t
	}

	if to := c.Query("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid 'to' date, expected YYYY-MM-DD")
		}
		// include the full end day
		t = t.Add(24*time.Hour - time.Millisecond)
		filter.To = &t
	}

	txs, err := h.service.ListTransactions(c.UserContext(), filter)
	if err != nil {
		return err
	}

	return c.JSON(txs)
}

// GetByID godoc
// @Summary Get transaction by ID
// @Tags Transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} domain.Transaction
// @Failure 404 {object} ErrorResponse
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	tx, err := h.service.GetTransaction(c.UserContext(), c.Params("id"))
	if err != nil {
		return err
	}

	return c.JSON(tx)
}
