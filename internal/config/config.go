package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	R2 struct {
		AccountID      string
		AccessKey      string
		SecretKey      string
		BucketName     string
		LinkExpiration int
	}
}
