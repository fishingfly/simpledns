package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func main() {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, "udp", "127.0.0.1:5354")
		},
	}

	//ips, _ := r.LookupHost(context.Background(), "test.service")
	//fmt.Println(ips[0])
	ips, _ := r.LookupHost(context.Background(), "cis-hub-dongguan-1.cmecloud.cn")
	fmt.Println(ips[0])
	ips, _ = r.LookupHost(context.Background(), "test1.ecis-suzhou-1.cmecloud.cn")
	fmt.Println(ips[0])
	ips, _ = r.LookupHost(context.Background(), "test2.ecis-suzhou-1.cmecloud.cn")
	fmt.Println(ips[0])
	ips, _ = r.LookupHost(context.Background(), "test2.ecis-hangzhou-1.cmecloud.cn")
	fmt.Println(ips[0])
}
