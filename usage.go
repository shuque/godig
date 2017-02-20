package main

import (
	"fmt"
	"os"
)

//
// usage() - Print usage string and exit.
//
func usage() {
	fmt.Println("godig version " + Version)
	fmt.Printf("Usage: %s [<options>] qname [qtype] [qclass]\n", Progname)
	fmt.Println(`
Supported Options:
  -h                   Print this usage string and exit
  -v                   Print program version and exit
  @server              Use specified server name or address as resolver
  -pNNN                Use NNN as the port number (default is 53)
  +tcp                 Use TCP as transport (default is UDP)
  +ignore              Ignore truncation, i.e. don't retry with TCP
  +retry=N             Set number of tries for UDP queries (default 3)
  +time=N              Set timeout (default TCP 5s, UDP 2s + exp backoff)
  +opcode=N            Set opcode to N (default is 0: Query)
  +adflag              Set AD (Authenticated Data) flag
  +cdflag              Set CD (Checking Disabled) flag
  +norecurse           Unset RD (Recursion Desired) bit
  +dnssec              Set DNSSEC-OK bit
  +bufsize=N           Use EDNS0 UDP payload size of N
  -4                   Use IPv4 transport
  -6                   Use IPv6 transport
  +0x20                Randomize case of query name (bit 0x20 hack)
  -x                   Do reverse DNS lookup on qname IP address
  +cookie[=x]          Use EDNS cookie option [with specified cookie]
  +subnet=X            Use specified address/mask as EDNS client subnet option
  +nsid                Send NSID option
  +edns=NNN            Set EDNS version (default 0)
  +ednsflags=N         Set EDNS flags field to N
  +ednsopt=###[:value] Set generic EDNS option
  +expire              Send and EDNS Expire option
  -yalg:name:key       Use TSIG with specified algorithm, name, and key
  +batch=filename      Run queries in specified batchfile, one per line
  +parallel=N          Use N concurrent queries at a time in batchfile mode
`)
	os.Exit(1)
}
