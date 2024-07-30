package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	prefixEnv = "ANTI_BRUTEFORCE_"

	blackListEnv = "BLACK_LIST"
	whiteListEnv = "WHITE_LIST"

	addrEnv = "LISTEN_ADDR"

	loginLimitEnv    = "N"
	passwordLimitEnv = "M"
	ipLimitEnv       = "K"

	redisURLEnv = "REDIS_URL"

	logLevelEnv = "LOG_LEVEL"

	bucketSizeEnv    = "BUCKET_SIZE"
	blockIntervalEnv = "BLOCK_INTERVAL"
	hostEnv          = "HOST"

	defaultAddr          = ":8081"
	defaultLoginLimit    = "10"
	defaultPasswordLimit = "100"
	defaultIPLimit       = "1000"
	defaultBucketSize    = "10"
	defaultBlockInterval = "3600"
	defaultWhiteListKey  = "ab:white"
	defaultBlackListKey  = "ab:black"
	defaultLogLevel      = "info"
	defaultRedisURL      = "redis://redis:6379/0"
	defaultHost          = "http://localhost:8081"
)

var (
	errZeroBucketSize    = errors.New("bucket size cannot be zero")
	errZeroIPLimit       = errors.New("ip limit cannot be zero")
	errZeroLoginLimit    = errors.New("login limit cannot be zero")
	errZeroPasswordLimit = errors.New("password limit cannot be zero")
	errSameRedisKeys     = errors.New("whitelist and blacklist cannot have the same keys")
)

type HostInfo struct {
	Addr string
	Host string
}

type Limits struct {
	LoginLimit    int
	PasswordLimit int
	IPLimit       int
	BucketSize    int
	BlockInterval float64
}

type Redis struct {
	WhiteListKey string
	BlackListKey string
	LogLevel     string
	URL          string
}

type Config struct {
	HostInfo HostInfo
	Limits   Limits
	Redis    Redis
}

func New() (Config, error) {
	cfg := Config{}
	n, err := strconv.Atoi(Env(loginLimitEnv, defaultLoginLimit))
	if err != nil {
		return cfg, err
	}
	cfg.Limits.LoginLimit = n

	m, err := strconv.Atoi(Env(passwordLimitEnv, defaultPasswordLimit))
	if err != nil {
		return cfg, err
	}
	cfg.Limits.PasswordLimit = m

	k, err := strconv.Atoi(Env(ipLimitEnv, defaultIPLimit))
	if err != nil {
		return cfg, err
	}
	cfg.Limits.IPLimit = k

	bucketSize, err := strconv.Atoi(Env(bucketSizeEnv, defaultBucketSize))
	if err != nil {
		return cfg, err
	}
	cfg.Limits.BucketSize = bucketSize

	blockInterval, err := strconv.ParseFloat(Env(blockIntervalEnv, defaultBlockInterval), 64)
	if err != nil {
		return cfg, err
	}
	cfg.Limits.BlockInterval = blockInterval

	cfg.Redis.WhiteListKey = Env(whiteListEnv, defaultWhiteListKey)
	cfg.Redis.BlackListKey = Env(blackListEnv, defaultBlackListKey)
	cfg.Redis.LogLevel = Env(logLevelEnv, defaultLogLevel)
	cfg.Redis.URL = Env(redisURLEnv, defaultRedisURL)
	cfg.HostInfo.Addr = Env(addrEnv, defaultAddr)
	cfg.HostInfo.Host = Env(hostEnv, defaultHost)

	return cfg, cfg.validate()
}

func Env(key string, defaultValue string) string {
	v, ok := os.LookupEnv(prefixEnv + key)
	if !ok {
		return defaultValue
	}
	return v
}

func (c Config) validate() error {
	if c.Limits.PasswordLimit == 0 {
		return errZeroPasswordLimit
	}

	if c.Limits.LoginLimit == 0 {
		return errZeroLoginLimit
	}

	if c.Limits.IPLimit == 0 {
		return errZeroIPLimit
	}

	if c.Limits.BucketSize == 0 {
		return errZeroBucketSize
	}

	if c.Redis.WhiteListKey == c.Redis.BlackListKey {
		return errSameRedisKeys
	}

	return nil
}
