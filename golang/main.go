package main

import (
	"./core"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Printf("http_proxy_go started\n")

	port := "127.0.0.1:8081"

	server, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("listen on %s\n", port)
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
