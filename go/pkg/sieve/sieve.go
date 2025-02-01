package sieve

import (
	"sieve/internal/runner"
)

type Sieve interface {
	NthPrime(int64) int64
}

func NewSieve() Sieve {
	return runner.NewRunner()
}
