package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

//
// makeOptrr() - construct OPT Psudo RR structure
//
func makeOptRR() *dns.OPT {

	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	if Options.bufsize > 0 {
		opt.SetUDPSize(Options.bufsize)
	} else {
		opt.SetUDPSize(BufsizeDefault)
	}
	if Options.dnssec {
		opt.SetDo()
	}

	if Options.ednsFlags != 0 {
		opt.Hdr.Ttl |= uint32(Options.ednsFlags)
	}

	if Options.nsid {
		e := new(dns.EDNS0_NSID)
		e.Code = dns.EDNS0NSID
		e.Nsid = ""
		opt.Option = append(opt.Option, e)
	}

	if Options.expire {
		e := new(dns.EDNS0_LOCAL)
		e.Code = dns.EDNS0EXPIRE
		e.Data = nil
		opt.Option = append(opt.Option, e)
	}

	if Options.cookie {
		e := new(dns.EDNS0_COOKIE)
		e.Code = dns.EDNS0COOKIE
		if Options.cookiedata == "" {
			b := make([]byte, 8)
			_, err := rand.Read(b)
			if err != nil {
				fmt.Println("Cookie generation error:", err)
				os.Exit(1)
			}
			e.Cookie = hex.EncodeToString(b)
		} else {
			e.Cookie = Options.cookiedata
		}
		opt.Option = append(opt.Option, e)
	}

	if Options.subnet != "" {
		parts := strings.Split(Options.subnet, "/")
		addr := parts[0]
		mask, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("Error parsing client subnet: %s: %s\n",
				Options.subnet, err)
			os.Exit(1)
		}
		e := new(dns.EDNS0_SUBNET)
		e.Code = dns.EDNS0SUBNET
		e.SourceNetmask = uint8(mask)
		e.SourceScope = 0
		ip := net.ParseIP(addr)
		if ip.To4() != nil { // IPv4 address
			e.Family = 1
			e.Address = ip.To4()
		} else { // IPv6 address
			e.Family = 2
			e.Address = ip
		}
		opt.Option = append(opt.Option, e)
	}

	if Options.ednsOpt != nil {
		for _, o := range Options.ednsOpt {
			e := new(dns.EDNS0_LOCAL)
			e.Code = o.code
			h, err := hex.DecodeString(o.data)
			if err != nil {
				fmt.Printf("Error decoding generic edns option data.\n")
				os.Exit(1)
			}
			e.Data = h
			opt.Option = append(opt.Option, e)
		}
	}

	opt.SetVersion(Options.ednsVersion)
	return opt
}
