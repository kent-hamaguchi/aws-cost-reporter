package cost

import "time"

// Repository cost repository
type Repository interface {
	GetMonthly(int, time.Month) (Cost, error)
}

// Cost コスト
type Cost struct {
	Amount float64
}
