package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
)

var gatewayConfigFile = flag.String("f", "etc/gateway.yaml", "tệp cấu hình gateway")

func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*gatewayConfigFile, &c)

	gw := gateway.MustNewServer(c)
	defer gw.Stop()

	fmt.Printf("Bắt đầu Gateway Server (REST API) tại %s:%d...\n", c.Host, c.Port)
	gw.Start()
}
