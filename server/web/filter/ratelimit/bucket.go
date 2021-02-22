package ratelimit

import "time"

// Bucket is an interface store ratelimit info
type Bucket interface {
	Take(amount uint) bool
	Capacity() uint
	Remaining() uint
	Rate() time.Duration
}

// BucketOption is constructor option
type BucketOption func(Bucket)
