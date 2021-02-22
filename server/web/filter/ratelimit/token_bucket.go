package ratelimit

import (
	"sync"
	"time"
)

type tokenBucket struct {
	sync.RWMutex
	remaining   uint
	capacity    uint
	lastCheckAt time.Time
	rate        time.Duration
}

// NewTokenBucket return an bucket that implements token bucket
func NewTokenBucket(opts []BucketOption) Bucket {
	b := &tokenBucket{lastCheckAt: time.Now()}
	for _, o := range opts {
		o(b)
	}
	return b
}

func withCapacity(capacity uint) BucketOption {
	return func(b Bucket) {
		bucket := b.(*tokenBucket)
		bucket.capacity = capacity
		bucket.remaining = capacity
	}
}

func withRate(rate time.Duration) BucketOption {
	return func(b Bucket) {
		bucket := b.(*tokenBucket)
		bucket.rate = rate
	}
}

func (b *tokenBucket) Remaining() uint {
	b.RLock()
	defer b.RUnlock()
	return b.remaining
}

func (b *tokenBucket) Rate() time.Duration {
	b.RLock()
	defer b.RUnlock()
	return b.rate
}

func (b *tokenBucket) Capacity() uint {
	b.RLock()
	defer b.RUnlock()
	return b.capacity
}

func (b *tokenBucket) Take(amount uint) bool {
	if b.rate <= 0 {
		return true
	}
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	times := uint(now.Sub(b.lastCheckAt) / b.rate)
	b.lastCheckAt = b.lastCheckAt.Add(time.Duration(times) * b.rate)
	b.remaining += times
	if b.remaining < amount {
		return false
	}
	b.remaining -= amount
	if b.remaining > b.capacity {
		b.remaining = b.capacity
	}
	return true
}
