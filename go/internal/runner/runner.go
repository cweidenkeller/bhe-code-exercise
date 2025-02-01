package runner

import (
	"encoding/binary"
	"iter"
	"log/slog"
	"math"

	"github.com/bits-and-blooms/bloom/v3"
)

type Runner interface {
	NthPrime(n int64) int64
}

func EstimatePrimeDensity(n int64) int64 {
	fn := float64(n)
	return int64(fn * math.Log(fn))
}

func NewRunner() Runner {
	b := bloom.New(10*2000, 5)
	g := GoRoutinesRunner{Bloom: b}
	return &g
}

type GoRoutinesRunner struct {
	Bloom *bloom.BloomFilter
}

func (j *GoRoutinesRunner) BloomInsert(n int64) {
	n1 := make([]byte, 8)
	binary.BigEndian.PutUint64(n1, uint64(n))
	j.Bloom.Add(n1)
}

func (j *GoRoutinesRunner) InBloom(n int64) bool {
	n1 := make([]byte, 8)
	binary.BigEndian.PutUint64(n1, uint64(n))
	return j.Bloom.Test(n1)
}

func (j *GoRoutinesRunner) NthPrime(n int64) int64 {
	if n <= 2 {
		return 2
	}
	pd := EstimatePrimeDensity(n)
	pd = pd + pd/3
	sr := Sieve(pd)
	slog.Info("asdf", "ddd", sr)
	return sr[n]
}

func Sieve(n int64) []int64 {
	primes := make([]bool, n+1)
	for i := int64(2); i <= n; i++ {
		primes[i] = true
	}

	for p := int64(2); p*p <= n; p++ {
		if primes[p] {
			for i := p * p; i <= n; i += p {
				primes[i] = false
			}
		}
	}

	var primeNumbers []int64
	for p := int64(2); p <= n; p++ {
		if primes[p] {
			primeNumbers = append(primeNumbers, p)
		}
	}

	return primeNumbers
}

func isPrime(n int64) bool {
	if n <= 2 {
		return false
	}

	sqrtN := int64(math.Sqrt(float64(n)))
	for i := int64(2); i <= sqrtN; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func generatePrimeNumbers() iter.Seq[int64] {
	return func(yield func(i int64) bool) {
		n := int64(0)

		for {
			if isPrime(n) {
				if !yield(n) {
					return
				}
			}

			n++
		}
	}
}
