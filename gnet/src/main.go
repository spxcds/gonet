package main

import (
	"./dns"
	"fmt"
	"os"
)

func main() {

	// t, cname, ip := dns.LookUpHost("114.114.114.114:53", "www.didichuxing.com")
	_, cname, ip := dns.LookUpHost("172.20.1.1:53", os.Args[1])

	// fmt.Println("time = ", t)
	for _, v := range cname {
		fmt.Println(v)
	}

	for _, v := range ip {
		fmt.Println(v)
	}
}
