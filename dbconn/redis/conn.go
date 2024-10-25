package cache

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/hsuanshao/go-tools/ctx"
	"github.com/sirupsen/logrus"
)

var (
	ErrWithoutConnConfig = errors.New("connection config is required")
	ErrEmptyURI          = errors.New("empty uri been input")
	ErrPingFailed        = errors.New("redis ping check get error")
	ErrHostIsEmptyStr    = errors.New("config host should not as empty string")
	ErrIncorrectPort     = errors.New("config port value is incorrect")
	ErrGetConnectClient  = errors.New("establish connection client failed")
)

// RedisConfig describe connection information
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func connect(ctx ctx.CTX, uri, password string) (client *redis.Client, err error) {
	if strings.TrimSpace(uri) == "" {
		return nil, ErrEmptyURI
	}
	avalibleCPU := runtime.NumCPU()
	poolSize := 4
	if avalibleCPU > poolSize {
		poolSize = avalibleCPU / 2
	}

	clientOpt := &redis.Options{
		Addr:     uri,
		PoolFIFO: false,
		PoolSize: poolSize,
		Password: password,
	}

	client = redis.NewClient(clientOpt)

	pong, err := client.Ping(ctx.Context).Result()
	if err != nil {
		ctx.WithField("err", err).Error("ping redis server get error")
		return nil, ErrPingFailed
	}
	ctx.WithField("pong", pong).Info("ping response message")

	return client, nil
}

func GetRedisConn(ctx ctx.CTX, redisConf *RedisConfig) (conn *redis.Client, err error) {
	if redisConf == nil {
		ctx.Error("nil redis config had been input")
		return nil, ErrWithoutConnConfig
	}

	uri, password := "", ""
	if redisConf != nil {
		if strings.TrimSpace(redisConf.Host) == "" {
			ctx.WithField("redisConf", *redisConf).Error("host should not as empty string input")
			return nil, ErrHostIsEmptyStr
		}

		if redisConf.Port <= 0 {
			ctx.WithField("redisConf", *redisConf).Error("port value should not smaller than 1")
			return nil, ErrIncorrectPort
		}

		uri = fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
		password = redisConf.Password
	}

	client, err := connect(ctx, uri, password)
	if err != nil {
		ctx.WithFields(logrus.Fields{"config": *redisConf, "err": err}).Error("try connect to redis server and get connect client failed")
		return nil, ErrGetConnectClient
	}

	return client, nil
}
