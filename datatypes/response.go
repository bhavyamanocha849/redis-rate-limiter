package datatypes

import "time"

type Response struct {
	IsValid   bool
	Count     uint64
	ExpiresAt time.Duration
}
