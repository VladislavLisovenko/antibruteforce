package app

import (
	"context"

	"github.com/VladislavLisovenko/antibruteforce/internal/list"
	"github.com/VladislavLisovenko/antibruteforce/internal/ratelimit"
)

type application struct {
	rateLimit ratelimit.RateLimit
	blackList list.List
	whiteList list.List
}

type Application interface {
	AddToWhiteList(ctx context.Context, ip string) error
	AddToBlackList(ctx context.Context, ip string) error
	DeleteFromWhiteList(ctx context.Context, ip string) error
	DeleteFromBlackList(ctx context.Context, ip string) error
	ResetAuth(login string, ip string)
	CheckAuth(login string, password string, ip string) bool
}

func NewApp(rt ratelimit.RateLimit, whiteList list.List, blackList list.List) Application {
	return &application{
		rateLimit: rt,
		blackList: blackList,
		whiteList: whiteList,
	}
}

func (a *application) CheckAuth(login string, password string, ip string) bool {
	if a.checkBlackList(ip) {
		return false
	}
	if a.checkWhiteList(ip) {
		return true
	}

	return a.rateLimit.Check(login, password, ip)
}

func (a *application) ResetAuth(login string, ip string) {
	a.rateLimit.Reset(login, ip)
}

func (a *application) AddToBlackList(ctx context.Context, ip string) error {
	return a.blackList.Add(ctx, ip)
}

func (a *application) DeleteFromBlackList(ctx context.Context, ip string) error {
	return a.blackList.Delete(ctx, ip)
}

func (a *application) AddToWhiteList(ctx context.Context, ip string) error {
	return a.whiteList.Add(ctx, ip)
}

func (a *application) DeleteFromWhiteList(ctx context.Context, ip string) error {
	return a.whiteList.Delete(ctx, ip)
}

func (a *application) checkWhiteList(ip string) bool {
	return a.whiteList.Check(ip)
}

func (a *application) checkBlackList(ip string) bool {
	return a.blackList.Check(ip)
}
