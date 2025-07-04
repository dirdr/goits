package domain

import "time"

type Account struct {
	AcccountID int64     `json:"account_id"`
	CreatedAt  time.Time `json:"created_at"`
}
