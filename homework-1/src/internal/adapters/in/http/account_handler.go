package handler

import (
	"github.com/gofiber/fiber/v2"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
)

type AccountHandler struct {
	service portin.TransactionService
}

func NewAccountHandler(service portin.TransactionService) *AccountHandler {
	return &AccountHandler{service: service}
}

// GetBalance godoc
// @Summary Get account balance
// @Tags Accounts
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} domain.AccountBalance
// @Failure 404 {object} ErrorResponse
// @Router /accounts/{accountId}/balance [get]
func (h *AccountHandler) GetBalance(c *fiber.Ctx) error {
	balance, err := h.service.GetAccountBalance(c.UserContext(), c.Params("accountId"))
	if err != nil {
		return err
	}

	return c.JSON(balance)
}

// GetSummary godoc
// @Summary Get account summary
// @Tags Accounts
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} domain.AccountSummary
// @Failure 404 {object} ErrorResponse
// @Router /accounts/{accountId}/summary [get]
func (h *AccountHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.service.GetAccountSummary(c.UserContext(), c.Params("accountId"))
	if err != nil {
		return err
	}

	return c.JSON(summary)
}
