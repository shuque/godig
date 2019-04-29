package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var Version string = "0.3"
var Progname string = path.Base(os.Args[0])

//
// Default parameters
//
var TimeoutInitial time.Duration = time.Second * 2
var TimeoutTCP time.Duration = time.Second * 5
var Retries int = 3
var ExponentialBackoff = true
var BufsizeDefault uint16 = 4096

// Generic EDNS option
type EdnsoptStruct struct {
	code uint16
	data string // hex-encoded data string
}

//
// For goroutine communications and synchronization:
//     wg: a sync counter to determine when last routine has ended.
//     numParallel: the default number of concurrent queries we allow.
//     tokens: a counting semapahore to bound the parallelism.
//     results: the channel over which query results are communicated.
//
var wg sync.WaitGroup
var numParallel uint16 = 20
var tokens chan struct{}
var results chan *ResponseInfo

//
// Response Information structure
//
type ResponseInfo struct {
	qname     string
	qtype     string
	qclass    string
	truncated bool
	retried   bool
	timeout   bool
	response  *dns.Msg
	rtt       time.Duration
	err       error
}

//
// makeMessage() - construct DNS message structure
//
func makeMessage(qname, qtype, qclass string) *dns.Msg {

	m := new(dns.Msg)
	m.Id = dns.Id()
	m.Opcode = Options.opcode

	if Options.norecurse {
		m.RecursionDesired = false
	} else {
		m.RecursionDesired = true
	}

	if Options.adflag {
		m.AuthenticatedData = true
	}

	if Options.cdflag {
		m.CheckingDisabled = true
	}

	if Options.edns {
		m.Extra = append(m.Extra, makeOptRR())
	}

	m.Question = make([]dns.Question, 1)
	qtype_int, ok := dns.StringToType[strings.ToUpper(qtype)]
	if !ok {
		fmt.Printf("%s: Unrecognized query type.\n", qtype)
		usage()
	}
	qclass_int, ok := dns.StringToClass[strings.ToUpper(qclass)]
	if !ok {
		fmt.Printf("%s: Unrecognized query class.\n", qclass)
		usage()
	}
	m.Question[0] = dns.Question{qname, qtype_int, qclass_int}

	return m
}

//
// doQuery() - perform DNS query with timeouts and retries as needed
//
func doQuery(qname, qtype, qclass string, use_tcp bool) (response *dns.Msg, rtt time.Duration, err error) {

	var retries = Options.retries
	var timeout = Options.itimeout

	m := makeMessage(qname, qtype, qclass)

	if use_tcp {
		return sendRequest(m, true, Options.tcptimeout)
	}

	for retries > 0 {

		response, rtt, err = sendRequest(m, false, timeout)
		if err == nil {
			break
		}
		if nerr, ok := err.(net.Error); ok && !nerr.Timeout() {
			break
		}
		retries--
		if retries > 0 && ExponentialBackoff {
			timeout = timeout * 2
		}
	}
	return response, rtt, err

}

//
// sendRequest() - send a DNS query
//
func sendRequest(m *dns.Msg, use_tcp bool, timeout time.Duration) (response *dns.Msg, rtt time.Duration, err error) {

	c := new(dns.Client)
	c.Timeout = timeout

	if Options.v6 {
		if use_tcp {
			c.Net = "tcp6"
		} else {
			c.Net = "udp6"
		}
	} else if Options.v4 {
		if use_tcp {
			c.Net = "tcp4"
		} else {
			c.Net = "udp4"
		}
	} else {
		if use_tcp {
			c.Net = "tcp"
		} else {
			c.Net = "udp"
		}
	}

	if Options.tsig != "" {
		alg, name, key := getTsigParams()
		c.TsigSecret = map[string]string{name: key}
		m.SetTsig(name, alg, 300, time.Now().Unix())
	}

	return c.Exchange(m, addressString(Options.server, Options.port))

}

//
// doit() - query dispatching goroutine
// Uses the "tokens" buffered channel as a counting semaphore to bound
// the number of concurrent queries in flight.
//
func doit(qname, qtype, qclass string) {

	if Options.batchfile != "" {
		defer wg.Done()
	}
	r := new(ResponseInfo)
	r.qname, r.qtype, r.qclass = qname, qtype, qclass

	response, rtt, err := doQuery(qname, qtype, qclass, Options.tcp)

	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			r.timeout = true
		}
	} else if response.MsgHdr.Truncated {
		r.truncated = true
		if !Options.ignore {
			r.retried = true
			response, rtt, err = doQuery(qname, qtype, qclass, true)
		}
	}

	<-tokens // Release token.

	r.response = response
	r.rtt = rtt
	r.err = err
	results <- r

}

//
// getQueryFromString()
//
func getQueryFromString(line string) (qname, qtype, qclass string) {

	fields := strings.Fields(line)
	numFields := len(fields)

	switch numFields {
	case 3:
		qname = dns.Fqdn(fields[0])
		qtype = strings.ToUpper(fields[1])
		qclass = strings.ToUpper(fields[2])
	case 2:
		qname = dns.Fqdn(fields[0])
		qtype = strings.ToUpper(fields[1])
		qclass = "IN"
	case 1:
		qname = dns.Fqdn(fields[0])
		qtype = "A"
		qclass = "IN"
	default:
		fmt.Printf("Batchfile line error: %s\n", line)
		return
	}

	return
}

//
// runBatchfile()
//
func runBatchfile(batchfile string) {

	var qname, qtype, qclass string

	t0 := time.Now()

	go func() {
		f, err := os.Open(batchfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()
			qname, qtype, qclass = getQueryFromString(line)
			if qname == "" {
				fmt.Printf("Batchfile line error: %s\n", line)
				continue
			}
			wg.Add(1)
			tokens <- struct{}{}          // Obtain token; blocks if channel is full
			go doit(qname, qtype, qclass) // releases token at the end
		}
		wg.Wait()
		close(results)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	for r := range results {
		fmt.Printf("### QUERY: %s %s %s\n", r.qname, r.qtype, r.qclass)
		printResponse(r)
		fmt.Println()
	}

	elapsed := time.Since(t0)
	fmt.Printf(";; Elapsed time for batch: %s\n", elapsed)

	return
}

//
// initialize()
//
func initialize() {

	log.SetFlags(0)
	// Per RFC 6891, use mnemonic "BADVERS" for rcode 16.
	dns.RcodeToString[16] = "BADVERS"
	return

}

//
// main()
//
func main() {

	var err error

	initialize()
	qname, qtype, qclass := parseArgs(os.Args[1:])

	tokens = make(chan struct{}, int(numParallel))
	results = make(chan *ResponseInfo)

	if Options.server == "" {
		Options.server, err = getSysResolver()
		if err != nil {
			fmt.Printf("failed to get resolver adddress.\n")
			os.Exit(1)
		}
	}

	if qtype == "AXFR" {
		zoneTransfer(qname)
		return
	}

	if Options.batchfile != "" {
		runBatchfile(Options.batchfile)
		return
	}

	tokens <- struct{}{} // Obtain token, which is released by doit()
	go doit(qname, qtype, qclass)
	r := <-results
	printResponse(r)

	return

}
