package main

import (
	"fmt"
	"github.com/miekg/dns"
	"os"
	"strconv"
	"strings"
	"time"
)

//
// Options
//
type OptionsStruct struct {
	port          int
	tcp           bool
	itimeout      time.Duration
	tcptimeout    time.Duration
	retries       int
	opcode        int
	ignore        bool
	adflag        bool
	cdflag        bool
	norecurse     bool
	edns          bool
	edns_version  uint8
	edns_flags    uint16
	edns_opt      []*EdnsoptStruct
	dnssec        bool
	bufsize       uint16
	v4            bool
	v6            bool
	randcase      bool
	reverseLookup bool
	server        string
	cookie        bool
	cookiedata    string
	subnet        string
	nsid          bool
	expire        bool
	tsig          string
	batchfile     string
}

var Options OptionsStruct = OptionsStruct{port: 53, tcp: false,
	itimeout: TimeoutInitial, tcptimeout: TimeoutTCP, retries: Retries,
	adflag: false, cdflag: false, norecurse: false,
	edns: false, edns_version: 0, dnssec: false}

//
// parseArgs() - parse command line arguments and set options
//
func parseArgs(args []string) (qname, qtype, qclass string) {

	var i int
	var arg string
	var err error
	var no_qname = true

	qtype = "A"
	qclass = "IN"

FORLOOP:
	for i, arg = range args {

		switch {
		case arg == "-h":
			usage()
		case arg == "-v":
			fmt.Println("godig version " + Version)
			fmt.Println("miekg/dns version " + dns.Version.String())
			os.Exit(1)
		case arg == "-x":
			Options.reverseLookup = true
		case arg == "+tcp":
			Options.tcp = true
		case arg == "+ignore":
			Options.ignore = true
		case strings.HasPrefix(arg, "+opcode="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+opcode="))
			if err != nil {
				fmt.Printf("Invalid opcode: %s\n", arg)
				usage()
			}
			Options.opcode = n
		case arg == "+adflag":
			Options.adflag = true
		case arg == "+cdflag":
			Options.cdflag = true
		case arg == "+norecurse":
			Options.norecurse = true
		case arg == "+dnssec":
			Options.dnssec = true
			Options.edns = true
		case arg == "-4":
			Options.v4 = true
		case arg == "-6":
			Options.v6 = true
		case arg == "+0x20":
			Options.randcase = true
		case strings.HasPrefix(arg, "@"):
			Options.server = arg[1:]
		case strings.HasPrefix(arg, "-p"):
			n, err := strconv.Atoi(arg[2:])
			if err != nil {
				fmt.Printf("Invalid port (-p): %s\n", arg[2:])
				usage()
			}
			Options.port = n
		case strings.HasPrefix(arg, "+bufsize="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+bufsize="))
			if err != nil {
				fmt.Printf("Invalid bufsize: %s\n", arg)
				usage()
			}
			Options.bufsize = uint16(n)
			Options.edns = true
		case strings.HasPrefix(arg, "+edns="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+edns="))
			if err != nil {
				fmt.Printf("Invalid edns: %s\n", arg)
				usage()
			}
			if n < 0 {
				fmt.Printf("Invalid edns '%d': not a valid number\n", n)
				usage()
			}
			if n > 255 {
				fmt.Printf("Invalid edns '%d': out of range\n", n)
				usage()
			}
			Options.edns_version = uint8(n)
			Options.edns = true
		case strings.HasPrefix(arg, "+ednsflags="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+ednsflags="))
			if err != nil {
				fmt.Printf("Invalid ednsflags: %s\n", arg)
				usage()
			}
			Options.edns = true
			Options.edns_flags = uint16(n)
		case strings.HasPrefix(arg, "+ednsopt="):
			s := strings.SplitN(strings.TrimPrefix(arg, "+ednsopt="), ":", 2)
			n, err := strconv.Atoi(s[0])
			if err != nil {
				fmt.Printf("Invalid ednsopt: %s\n", arg)
				usage()
			}
			o := new(EdnsoptStruct)
			o.code = uint16(n)
			if len(s) == 2 {
				o.data = s[1]
			}
			Options.edns = true
			Options.edns_opt = append(Options.edns_opt, o)
		case strings.HasPrefix(arg, "+retry="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+retry="))
			if err != nil {
				fmt.Printf("Invalid retry parameter: %s\n", arg)
				usage()
			}
			Options.retries = n
		case strings.HasPrefix(arg, "+time="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+time="))
			if err != nil {
				fmt.Printf("Invalid timeout parameter: %s\n", arg)
				usage()
			}
			Options.itimeout = time.Duration(n) * time.Second
			Options.tcptimeout = time.Duration(n) * time.Second
		case arg == "+cookie":
			Options.edns = true
			Options.cookie = true
		case strings.HasPrefix(arg, "+cookie="):
			Options.edns = true
			Options.cookie = true
			Options.cookiedata = strings.TrimPrefix(arg, "+cookie=")
		case arg == "+nsid":
			Options.edns = true
			Options.nsid = true
		case arg == "+expire":
			Options.edns = true
			Options.expire = true
		case strings.HasPrefix(arg, "+subnet="):
			Options.edns = true
			Options.subnet = strings.TrimPrefix(arg, "+subnet=")
		case strings.HasPrefix(arg, "-y"):
			Options.tsig = arg[2:]
		case strings.HasPrefix(arg, "+batch="):
			Options.batchfile = strings.TrimPrefix(arg, "+batch=")
		case strings.HasPrefix(arg, "+parallel="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "+parallel="))
			if err != nil {
				fmt.Printf("Invalid #parallel queries: %s\n", arg)
				usage()
			}
			numParallel = uint16(n)
		case strings.HasPrefix(arg, "-"):
			fmt.Printf("Invalid option: %s\n", arg)
			usage()
		case strings.HasPrefix(arg, "+"):
			fmt.Printf("Invalid option: %s\n", arg)
			usage()
		default:
			no_qname = false
			break FORLOOP
		}

	}

	if no_qname {
		qname, qtype = ".", "NS"
	} else {
		switch len(args) - i {
		case 3:
			qname, qtype, qclass = args[i], args[i+1], args[i+2]
		case 2:
			qname, qtype = args[i], args[i+1]
		case 1:
			qname = args[i]
		default:
			fmt.Printf("Too many arguments.\n")
			usage()
		}
	}

	if Options.randcase {
		qname, err = randomizeCase(qname)
		if err != nil {
			fmt.Printf("Error randomizing case: %s", err)
			os.Exit(1)
		}
	}

	if Options.reverseLookup {
		arpaname, err := dns.ReverseAddr(qname)
		if err != nil {
			fmt.Printf("Invalid IP address for -x: %s\n", qname)
			usage()
		}
		qname = arpaname
		qtype = "PTR"
	}
	qname = dns.Fqdn(qname)

	return
}
