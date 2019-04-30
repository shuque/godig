package main

import (
	"encoding/hex"
	"fmt"
	"github.com/miekg/dns"
	"strconv"
)

//
// opcodeString()
//
func opcodeString(c int) (s string) {
	var ok bool
	s, ok = dns.OpcodeToString[c]
	if !ok {
		s = "OPCODE" + strconv.Itoa(c)
	}
	return
}

//
// rcodeString()
//
func rcodeString(c int) (s string) {
	var ok bool
	s, ok = dns.RcodeToString[c]
	if !ok {
		s = "RCODE" + strconv.Itoa(c)
	}
	return
}

//
// printHeader() - print DNS header (and section RR counts)
//
func printHeader(m *dns.Msg) {
	fmt.Printf(";; ->>HEADER<<- ;; opcode: %s, status: %s, id: %d\n",
		opcodeString(m.MsgHdr.Opcode),
		rcodeString(m.MsgHdr.Rcode),
		m.MsgHdr.Id)
	fmt.Printf(";; flags:")
	if m.MsgHdr.Response {
		fmt.Printf(" qr")
	}
	if m.MsgHdr.Authoritative {
		fmt.Printf(" aa")
	}
	if m.MsgHdr.Truncated {
		fmt.Printf(" tc")
	}
	if m.MsgHdr.RecursionDesired {
		fmt.Printf(" rd")
	}
	if m.MsgHdr.RecursionAvailable {
		fmt.Printf(" ra")
	}
	if m.MsgHdr.AuthenticatedData {
		fmt.Printf(" ad")
	}
	if m.MsgHdr.CheckingDisabled {
		fmt.Printf(" cd")
	}
	fmt.Printf("; QUERY: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d\n",
		len(m.Question), len(m.Answer), len(m.Ns), len(m.Extra))
	return
}

//
// printOptSection()
//
func printOptSection(m *dns.Msg, opt *dns.OPT) {
	flags := ""
	if opt.Do() {
		flags = "do"
	}
	fmt.Printf("\n; EDNS: version %d; flags: %s; udp: %d; ercode=%d\n",
		opt.Version(), flags, opt.UDPSize(), opt.ExtendedRcode())
	rcode_hi_bits := (opt.Hdr.Ttl >> 24) & 0xff
	if rcode_hi_bits != 0 {
		m.MsgHdr.Rcode = opt.ExtendedRcode()
	}

	for _, o := range opt.Option {
		switch o.(type) {
		case *dns.EDNS0_NSID:
			h, err := hex.DecodeString(o.String())
			if err != nil {
				fmt.Printf("; NSID: %s\n", o.String())
			} else {
				fmt.Printf("; NSID: %s (%s)\n", o.String(), string(h))
			}
		case *dns.EDNS0_SUBNET:
			fmt.Printf("; SUBNET: %s\n", o.String())
		case *dns.EDNS0_COOKIE:
			fmt.Printf("; COOKIE: %s\n", o.String())
		case *dns.EDNS0_UL:
			fmt.Printf("; UPDATE LEASE: %s\n", o.String())
		case *dns.EDNS0_LLQ:
			fmt.Printf("; LLQ: %s\n", o.String())
		case *dns.EDNS0_DAU:
			fmt.Printf("; DAU: %s\n", o.String())
		case *dns.EDNS0_DHU:
			fmt.Printf("; DHU: %s\n", o.String())
		case *dns.EDNS0_N3U:
			fmt.Printf("; N3HU: %s\n", o.String())
		case *dns.EDNS0_LOCAL:
			fmt.Printf("; LOCAL: %s\n", o.String())
		case *dns.EDNS0_PADDING:
			fmt.Printf("; PADDING: %s\n", o.String())
		default:
			fmt.Printf("; Option%d: %s\n", o.Option(), o.String())
		}
	}
}

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

	printHeader(r.response)

	opt := r.response.IsEdns0()
	if opt != nil {
		printOptSection(r.response, opt)
	}

	if len(r.response.Question) > 0 {
		fmt.Printf("\n;; QUESTION SECTION:\n")
		fmt.Printf("%s\n", r.response.Question[0].String())
	}

	if len(r.response.Answer) > 0 {
		fmt.Printf("\n;; ANSWER SECTION:\n")
		for _, rr := range r.response.Answer {
			fmt.Printf("%s\n", rr.String())
		}
	}

	if len(r.response.Ns) > 0 {
		fmt.Printf("\n;; AUTHORITY SECTION:\n")
		for _, rr := range r.response.Ns {
			fmt.Printf("%s\n", rr.String())
		}
	}

	if ((opt != nil) && len(r.response.Extra) > 1) || ((opt == nil) && len(r.response.Extra) > 0) {
		fmt.Printf("\n;; ADDITIONAL SECTION:\n")
		for _, rr := range r.response.Extra {
			rh := rr.Header()
			if rh.Rrtype != dns.TypeOPT {
				fmt.Printf("%s\n", rr.String())
			}
		}
	}

	// The size computation below may be different than the size of
	// the actual response, which the go library doesn't retain. It
	// is the size computed by re-packing the response message data
	// structure.
	fmt.Printf("\n;; time: %.3fms, server: %s:%d, size: %d bytes\n",
		float64(r.rtt)/1000000.0, Options.server, Options.port, r.response.Len())
	return
}
