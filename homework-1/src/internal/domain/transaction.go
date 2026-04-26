package domain

import "time"

type TransactionType string
type TransactionStatus string

const (
	TypeDeposit    TransactionType = "deposit"
	TypeWithdrawal TransactionType = "withdrawal"
	TypeTransfer   TransactionType = "transfer"

	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID          string            `json:"id"`
	FromAccount string            `json:"fromAccount"`
	ToAccount   string            `json:"toAccount"`
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	Type        TransactionType   `json:"type"`
	Timestamp   time.Time         `json:"timestamp"`
	Status      TransactionStatus `json:"status"`
}
