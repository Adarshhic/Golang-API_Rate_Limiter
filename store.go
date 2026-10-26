
package main

import "time"


type BucketStore interface {
	Get(ip string) (*TokenBucket, bool)
	Set(ip string, bucket *TokenBucket)
	DeleteInactive(timeout time.Duration)
}
