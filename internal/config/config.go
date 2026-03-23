package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	CacheConf cache.CacheConf
	R2        struct {
		AccountID      string
		AccessKey      string
		SecretKey      string
		BucketName     string
		LinkExpiration int
	}
	DB struct {
		Host         string
		Port         int
		User         string
		Password     string
		DBName       string
		SSLMode      string
		MaxOpenConns int
		MaxIdleConns int
	}
}
