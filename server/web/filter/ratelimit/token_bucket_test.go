package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRate(t *testing.T) {
	b := NewTokenBucket([]BucketOption{withRate(1 * time.Second)}).(*tokenBucket)
	assert.Equal(t, b.Rate(), 1*time.Second)
}

func TestRemainingAndCapacity(t *testing.T) {
	b := NewTokenBucket([]BucketOption{withCapacity(10)}).(*tokenBucket)
	assert.Equal(t, b.Remaining(), uint(10))
	assert.Equal(t, b.Capacity(), uint(10))
}

func TestTake(t *testing.T) {
	var opts []BucketOption
	opts = append(opts, withCapacity(10), withRate(10*time.Millisecond))
	b := NewTokenBucket(opts).(*tokenBucket)
	for i := 0; i < 10; i++ {
		assert.True(t, b.Take(1))
	}
	assert.False(t, b.Take(1))
	assert.Equal(t, b.Remaining(), uint(0))
	opts = append(opts, withRate(1*time.Millisecond))
	opts = append(opts, withCapacity(1))
	b = NewTokenBucket(opts).(*tokenBucket)
	assert.True(t, b.Take(1))
	time.Sleep(2 * time.Millisecond)
	assert.True(t, b.Take(1))
}
