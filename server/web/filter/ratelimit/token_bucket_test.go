package ratelimit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
)

func TestRate(t *testing.T) {
	b := NewTokenBucket([]BucketOption{WithRate(1 * time.Second)})
	assert.Equal(t, b.Rate(), 1 * time.Second)
}

func TestRemainingAndMax(t *testing.T) {
	b := NewTokenBucket([]BucketOption{WithMax(10)})
	assert.Equal(t, b.Remaining(), float64(10))
	assert.Equal(t, b.Max(), float64(10))
}

func TestTake(t *testing.T) {
	var opts []BucketOption
	opts = append(opts, WithMax(10))
	b := NewTokenBucket(opts)
	for i := 0; i < 10; i++ {
		assert.True(t, b.Take(1))
	}
	assert.False(t, b.Take(1))
	assert.Equal(t, b.Remaining(), float64(0))
	opts = append(opts, WithRate(1 * time.Millisecond))
	opts = append(opts, WithMax(1))
	b = NewTokenBucket(opts)
	assert.True(t, b.Take(1))
	time.Sleep(2 * time.Millisecond)
	fmt.Println(b.Remaining())
	assert.True(t, b.Take(1))
}