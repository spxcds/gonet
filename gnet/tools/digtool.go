package main

import (
	"fmt"
	"gnet/dig"
)

var (
    server = "1.1.1.1"
    port = "53"
)

var ch chan record

type record struct {
    domain string
    subnet string
    qtype  string
}

func mydig() {
    var r record
    for {
        r = <-ch
        rtt, cname, ip, err := dig.Dig(server, port, r.subnet, r.domain, r.qtype)
        if err != nil {
            fmt.Println("rtt = ", rtt, " cname = ", cname, " ip = ", ip, " err = ", err)
        }
    }
}

func main() {
    ch = make(chan record, 200)
    for i := 0; i < 10; i++ {
        go mydig()
    }
    for {
        var r record
        if _, err := fmt.Scanf("%s%s%s", &r.domain, &r.subnet, &r.qtype); err != nil {
            break
        }
        ch <- r
    }
}
