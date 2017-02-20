package main

import (
	"fmt"
	"github.com/miekg/dns"
)

//
// printResponse() - print info about a DNS response
//
func printResponse(r *ResponseInfo) {

	if r.err != nil && !r.truncated {
		fmt.Printf("DNS query failed: %s\n", r.err)
		return
	}

	if r.truncated {
		fmt.Printf("UDP response was truncated.")
		if r.retried {
			fmt.Printf(" Retried over TCP.\n")
		} else {
			fmt.Printf(" Not retried over TCP.\n")
		}
	}

	opt := r.response.IsEdns0()
	if opt != nil {
		rcode_hi_bits := (opt.Hdr.Ttl >> 24) & 0xff
		if rcode_hi_bits != 0 {
			r.response.MsgHdr.Rcode = opt.ExtendedRcode()
		}
	}

	fmt.Printf(";; ->>HEADER<<- %s QUERY: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d\n",
		r.response.MsgHdr.String(),
		len(r.response.Question),
		len(r.response.Answer),
		len(r.response.Ns),
		len(r.response.Extra))

	if opt != nil {
		println(opt.String())
	}

	if len(r.response.Question) > 0 {
		println("\n;; QUESTION SECTION:")
		fmt.Printf("%s\n", r.response.Question[0].String())
	}

	if len(r.response.Answer) > 0 {
		println("\n;; ANSWER SECTION:")
		for _, rr := range r.response.Answer {
			fmt.Printf("%s\n", rr.String())
		}
	}

	if len(r.response.Ns) > 0 {
		println("\n;; AUTHORITY SECTION:")
		for _, rr := range r.response.Ns {
			fmt.Printf("%s\n", rr.String())
		}
	}

	if ((opt != nil) && len(r.response.Extra) > 1) || ((opt == nil) && len(r.response.Extra) > 0) {
		println("\n;; ADDITIONAL SECTION:")
		for _, rr := range r.response.Extra {
			rh := rr.Header()
			if rh.Rrtype != dns.TypeOPT {
				fmt.Printf("%s\n", rr.String())
			}
		}
	}

	fmt.Printf("\n;; ResponseTime: %.3fms\n", float64(r.rtt)/1000000.0)
	return
}
