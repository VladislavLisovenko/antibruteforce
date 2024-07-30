package ratelimit

import (
	"sync"

	"github.com/VladislavLisovenko/antibruteforce/pkg/bucket"
)

type RateLimit interface {
	Check(login string, password string, ip string) bool
	Reset(login string, ip string)
	Cleanup()
}

type rateLimit struct {
	loginBucket    bucket.Bucket
	passwordBucket bucket.Bucket
	ipBucket       bucket.Bucket
}

func New(loginLimit int, passwordLimit int, ipLimit int, bucketSize int, interval float64) RateLimit {
	return &rateLimit{
		loginBucket:    bucket.New(bucketSize, loginLimit, interval),
		passwordBucket: bucket.New(bucketSize, passwordLimit, interval),
		ipBucket:       bucket.New(bucketSize, ipLimit, interval),
	}
}

func (r *rateLimit) Check(login string, password string, ip string) bool {
	okLogin := r.loginBucket.Check(login)
	okPassword := r.passwordBucket.Check(password)
	okIP := r.ipBucket.Check(ip)
	return okLogin && okPassword && okIP
}

func (r *rateLimit) Reset(login string, ip string) {
	r.loginBucket.ResetKey(login)
	r.ipBucket.ResetKey(ip)
}

func (r *rateLimit) Cleanup() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		r.loginBucket.Cleanup()
	}()
	go func() {
		defer wg.Done()
		r.passwordBucket.Cleanup()
	}()
	go func() {
		wg.Done()
		r.ipBucket.Cleanup()
	}()

	wg.Wait()
}
