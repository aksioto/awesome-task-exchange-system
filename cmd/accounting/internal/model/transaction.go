package model

import (
	"time"
)

type Transaction struct {
	Id          int       `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	TaskID      string    `json:"task_id,omitempty" db:"task_id"`
	Description string    `json:"description" db:"description"`
	Increase    int       `json:"increase,omitempty" db:"increase"`
	Decrease    int       `json:"decrease,omitempty" db:"decrease"`
	Balance     int       `json:"balance,omitempty" db:"balance"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}

type TransactionType struct {
}

type UserDailyStatistics struct {
	Balance      int           `json:"balance,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

type ProfitStatistics struct {
	Profit int `json:"profit,omitempty"`
}
