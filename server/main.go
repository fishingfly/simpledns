package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var Level2Records = map[string]string{
	"cis-hub-dongguan-1.cmecloud.cn.": "192.168.0.4",
}
var Level3Records = map[string]string{
	"*.ecis-suzhou-1.cmecloud.cn.": "192.168.0.5",
	"*.ecis-hangzhou-1.cmecloud.cn.": "192.168.0.6",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			pointsCount := strings.Count(q.Name, ".")
			if pointsCount == 3 { // 二级域名，全域名匹配
				if _, ok := Level2Records[q.Name]; ok {
					log.Printf("ip is %s\n", Level2Records[q.Name])
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, Level2Records[q.Name]))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			}
			if pointsCount == 4 {// 三级域名
				arr := strings.SplitN(q.Name, ".", 2)
				log.Printf("三级域名分割 for %s\n", arr[1])
				// 从三级域名开始匹配
				if _, ok := Level3Records["*." + arr[1]]; ok {
					log.Printf("ip is %s\n", Level3Records["*." + arr[1]])
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, Level3Records["*." + arr[1]]))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc("service.", handleDnsRequest)
	dns.HandleFunc("cis-hub-dongguan-1.cmecloud.cn.", handleDnsRequest)
	dns.HandleFunc("ecis-suzhou-1.cmecloud.cn.", handleDnsRequest)
	dns.HandleFunc("ecis-hangzhou-1.cmecloud.cn.", handleDnsRequest)
	// start server
	port := 5354
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
