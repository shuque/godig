package main

import (
	"fmt"
	"github.com/miekg/dns"
	"os"
	"strings"
	"time"
)

//
// getTsigParams()
//
func getTsigParams() (alg, name, key string) {

	parts := strings.Split(Options.tsig, ":")
	alg, name, key = parts[0], parts[1], parts[2]
	return alg, name, key

}

//
// zoneTransfer()
//
func zoneTransfer(zonename string) {

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(zonename)

	if Options.tsig != "" {
		alg, name, key := getTsigParams()
		t.TsigSecret = map[string]string{name: key}
		m.SetTsig(name, alg, 300, time.Now().Unix())
	}

	t0 := time.Now()
	c, err := t.In(m, addressString(Options.server, Options.port))
	if err != nil {
		fmt.Printf("zone xfr failed: %s\n", err)
		os.Exit(1)
	}

	for e := range c {
		for _, rr := range e.RR {
			fmt.Println(rr)
		}
	}
	t1 := time.Now()

	// elapsed time here is for transfer plus printing
	fmt.Printf("\n;; ElapsedTime: %v\n", t1.Sub(t0))
	return

}
