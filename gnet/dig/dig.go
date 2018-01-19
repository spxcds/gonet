package dig

import (
	"net"
	"strings"
	"time"

	"errors"
	"github.com/miekg/dns"
)

var DNSTypeMap = map[string]uint16{
	"A":     1,
	"NS":    2,
	"CNAME": 5,
	"SOA":   6,
	"PTR":   12,
	"MX":    15,
	"TXT":   16,
	"SRV":   33,
}

func Dig(nameserver, port, subnet, domain string, args ...string) (delay time.Duration, cname []string, ip []string, DIGerr error) {

	client := new(dns.Client)
	msg := new(dns.Msg)
	if len(args) != 0 {
		msg.SetQuestion(dns.Fqdn(domain), DNSTypeMap[args[0]])
	} else {
		msg.SetQuestion(dns.Fqdn(domain), dns.TypeDNSKEY)
	}

	if subnet != "" {
		opt := new(dns.OPT)
		opt.Hdr.Name = "."
		opt.Hdr.Rrtype = dns.TypeOPT
		opt.SetUDPSize(4096)
		e := new(dns.EDNS0_SUBNET)
		e.Code = dns.EDNS0SUBNET
		e.Family = 1
		e.SourceNetmask = 32
		e.SourceScope = 0
		e.Address = net.ParseIP(subnet).To4()
		opt.Option = append(opt.Option, e)
		msg.Extra = append(msg.Extra, opt)
	}

	nameserver = nameserver + ":" + port

	result, delay, err := client.Exchange(msg, nameserver)
	if err != nil {
		DIGerr = err
		return
	}

	/*
		if result.MsgHdr.Rcode != 0 || len(result.Answer) == 0 {
			DIGerr = errors.New("return msg error")
			return
		}
	*/
	if result.MsgHdr.Rcode != 0 {
		DIGerr = errors.New("return msg error")
		return
	}

	var res []string = make([]string, len(result.Answer))
	for i, x := range result.Answer {
		s := strings.Split(x.String(), "\t")
		res[i] = s[len(s)-1]
	}
	cname = res
	ip = res
	return
	/*
		switch  args[0] {
		case "CNAME":
			return rtt, res, []string{}, nil
		case "A":
			return rtt, []string{}, res, nil
		}
		return rtt, []string{}, []string{}, fmt.Errorf("return error\n")
	*/
}
