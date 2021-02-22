package ratelimit

type Bucket interface {
	Take(amount float64) bool
}

type BucketOption func(Bucket)

