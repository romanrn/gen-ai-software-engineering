package domain

import "time"

type AccountBalance struct {
	AccountID string  `json:"accountId"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
}

type AccountSummary struct {
	AccountID           string    `json:"accountId"`
	TotalDeposits       float64   `json:"totalDeposits"`
	TotalWithdrawals    float64   `json:"totalWithdrawals"`
	TransactionCount    int       `json:"transactionCount"`
	LastTransactionDate time.Time `json:"lastTransactionDate"`
}
