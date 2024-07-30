package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VladislavLisovenko/antibruteforce/internal/app"
	"github.com/VladislavLisovenko/antibruteforce/internal/config"
	"github.com/VladislavLisovenko/antibruteforce/internal/keyvaluestorage"
	"github.com/VladislavLisovenko/antibruteforce/internal/logger"
	"github.com/VladislavLisovenko/antibruteforce/internal/ratelimit"
	"github.com/VladislavLisovenko/antibruteforce/internal/server"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg, err := logger.New(cfg.Redis.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1) //nolint:gocritic
	}

	rURL, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	rClient := redis.NewClient(rURL)

	whiteList, err := keyvaluestorage.New(ctx, cfg.Redis.WhiteListKey, *rClient)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	blackList, err := keyvaluestorage.New(ctx, cfg.Redis.BlackListKey, *rClient)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	lim := cfg.Limits
	rt := ratelimit.New(lim.LoginLimit, lim.PasswordLimit, lim.IPLimit, lim.BucketSize, lim.BlockInterval)

	a := app.NewApp(rt, whiteList, blackList)
	s := server.New(cfg.HostInfo.Addr, a, logg)

	go shutdown(ctx, s, logg)
	go clearRateLimit(ctx, rt, cfg.Limits.BlockInterval)

	logg.Info("antibruteforce is running")

	if err := s.Start(ctx); err != nil {
		logg.Error(err)
		cancel()
	}

	<-ctx.Done()
}

func shutdown(ctx context.Context, s server.Server, logg logger.Logger) {
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		logg.Error("failed to stop http-server: " + err.Error())
	}

	logg.Info("antibruteforce is shutdown")
}

func clearRateLimit(ctx context.Context, rt ratelimit.RateLimit, interval float64) {
	c := time.Tick(time.Duration(interval) * time.Second)
	for {
		select {
		case <-c:
			rt.Cleanup()
		case <-ctx.Done():
			return
		}
	}
}
