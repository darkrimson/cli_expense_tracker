package model

import (
	"time"
)

type Expense struct {
	ID          int
	Date        time.Time
	Description string
	Amount      float64
}
