package ratelimit

import (
	"sync"
	"time"
)

type tokenBucket struct {
	sync.RWMutex
	remaining float64
	max float64
	lastCheckAt time.Time
	rate time.Duration
}

func NewTokenBucket(opts []BucketOption) *tokenBucket {
	b := &tokenBucket{}
	for _, o := range opts {
		o(b)
	}
	return b
}

func WithMax(max float64) BucketOption {
	return func(b Bucket) {
		bucket := b.(*tokenBucket)
		bucket.max = max
		bucket.remaining = max
	}
}

func WithRate(rate time.Duration) BucketOption {
	return func(b Bucket) {
		bucket := b.(*tokenBucket)
		bucket.rate = rate
	}
}

func (b *tokenBucket) Remaining() float64 {
	b.RLock()
	defer b.RUnlock()
	return b.remaining
}

func (b *tokenBucket) LastCheckAt() time.Time {
	b.RLock()
	defer b.RUnlock()
	return b.lastCheckAt
}

func (b *tokenBucket) Rate() time.Duration {
	b.RLock()
	defer b.RUnlock()
	return b.rate
}

func (b *tokenBucket) Max() float64 {
	b.RLock()
	defer b.RUnlock()
	return b.max
}

func (b *tokenBucket) Take(amount float64) bool {
	if amount <= 0 {
		return false
	}
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	if !b.lastCheckAt.IsZero() && b.rate != 0 {
		times := float64(now.Sub(b.lastCheckAt) / b.rate)
		b.remaining += times
	}
	b.lastCheckAt = now
	if b.remaining > b.max {
		b.remaining = b.max
	}
	if b.remaining < amount {
		return false
	}
	b.remaining -= amount
	return true
}


