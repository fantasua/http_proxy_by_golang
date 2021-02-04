package core

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

type ProxyServer struct {
	// 缓冲区
	buf [4096]byte
	// 客户端
	client net.Conn
	// 服务端
	server net.Conn
	// 是否使用https
	useHTTPS bool
}

func NewProxyServer(client net.Conn) *ProxyServer {
	return &ProxyServer{
		buf:      [4096]byte{},
		client:   client,
		server:   nil,
		useHTTPS: false,
	}
}

func (ps *ProxyServer) Handle() {
	if ps.client == nil {
		fmt.Errorf("connection is nil\n")
		return
	}
	defer ps.client.Close()

	fmt.Printf("remote addr: %v\n", ps.client.RemoteAddr())

	// 从请求中读取数据
	_, err := ps.client.Read(ps.buf[:])
	if err != nil {
		log.Println(err)
		return
	}

	var method, URL string
	fmt.Sscanf(string(ps.buf[:bytes.IndexByte(ps.buf[:], '\n')]), "%s%s", &method, &URL)
	targetURL, err := url.Parse(URL)
	if err != nil {
		log.Println(err)
		return
	}

	// 判断连接类型
	// "CONNECT"为https协议
	if method == "CONNECT" {
		ps.useHTTPS = true
	} else {
		// 非"CONNECT"则为http协议
		ps.useHTTPS = false
	}

	// 处理URL 并进行连接
	if ps.useHTTPS {
		// HTTPS
		ps.server, err = ps.processHTTPS(targetURL)
		if err != nil {
			return
		}
		// HTTPS需要向客户端告知连接建立完毕
		fmt.Fprint(ps.client, "HTTP/1.1 200 Connection established\\r\\n\\r\\n")
	} else {
		// HTTP
		ps.server, err = ps.processHTTPS(targetURL)
		if err != nil {
			return
		}
		// 将请求转发给目标server
		ps.server.Write(ps.buf[:])
	}
}

func (ps *ProxyServer) processHTTPS(targetURL *url.URL) (net.Conn, error) {
	address := targetURL.Scheme + ":" + targetURL.Opaque
	return ps.dial("tcp", address)
}

func (ps *ProxyServer) processHTTP(targetURL *url.URL) (net.Conn, error) {
	address := targetURL.Host
	// Host部分若不带端口 则为默认的80
	if strings.Index(targetURL.Host, ":") == -1 {
		address = targetURL.Host + ":80"
	}
	return ps.dial("tcp", address)
}

func (ps *ProxyServer) dial(network, address string) (net.Conn, error) {
	server, err := net.Dial(network, address)
	if err != nil {
		log.Println(err)
	}
	return server, err
}
