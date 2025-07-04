package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Redis redis.RedisConf
	Mysql struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
}
