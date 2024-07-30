package keyvaluestorage

import (
	"context"
	"net/netip"
	"sync"

	"github.com/redis/go-redis/v9"
)

type keyValueStorage struct {
	l           map[netip.Prefix]struct{}
	redisClient redis.Client
	redisKey    string
	mu          sync.RWMutex
}

type KeyValueStorage interface {
	Add(ctx context.Context, el string) error
	Delete(ctx context.Context, el string) error
	Check(el string) bool
	Reset(ctx context.Context) error
}

func New(ctx context.Context, redisKey string, redisClient redis.Client) (KeyValueStorage, error) {
	elements, err := redisClient.SMembers(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	l := make(map[netip.Prefix]struct{}, len(elements))
	for _, v := range elements {
		net, err := netip.ParsePrefix(v)
		if err != nil {
			return nil, err
		}
		l[net] = struct{}{}
	}

	return &keyValueStorage{
		l:           l,
		redisKey:    redisKey,
		redisClient: redisClient,
		mu:          sync.RWMutex{},
	}, nil
}

func (l *keyValueStorage) Add(ctx context.Context, el string) error {
	net, err := netip.ParsePrefix(el)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.l[net] = struct{}{}

	r := l.redisClient.SAdd(ctx, l.redisKey, el)
	return r.Err()
}

func (l *keyValueStorage) Delete(ctx context.Context, el string) error {
	net, err := netip.ParsePrefix(el)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.l, net)

	r := l.redisClient.SRem(ctx, l.redisKey, el)
	return r.Err()
}

func (l *keyValueStorage) Check(el string) bool {
	addr, err := netip.ParseAddr(el)
	if err != nil {
		return false
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	for net := range l.l {
		if net.Contains(addr) {
			return true
		}
	}

	return false
}

func (l *keyValueStorage) Reset(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l = make(map[netip.Prefix]struct{})
	return l.redisClient.Del(ctx, l.redisKey).Err()
}
