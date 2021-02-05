package main

import (
	"./core"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Printf("http_proxy_go started\n")

	proxy_addr := "127.0.0.1:2333"

	server, err := net.Listen("tcp", proxy_addr)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("listen on %s\n", proxy_addr)
	for {
		client, err := server.Accept()
		// 错误处理
		if err != nil {
			fmt.Errorf("fail to Accept local addr %s\n", client.LocalAddr().String())
			continue
		}

		// 处理连接
		fmt.Printf("conn received\n")
		proxy := core.NewProxyServer(client)
		go proxy.Handle()
	}

	return
}
