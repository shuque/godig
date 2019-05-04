package main

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

//
// getSysResolver() - obtain (1st) system default resolver address
//
func getSysResolver() (resolver string, err error) {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err == nil {
		resolver = config.Servers[0]
	} else {
		fmt.Println("Error processing /etc/resolv.conf: " + err.Error())
	}
	return
}

//
// isAlpha() - is given character (byte) alphabetic?
//
func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

//
// randomizeCase() - randomize case of qname string (bit 0x20 hack)
//
func randomizeCase(qname string) (string, error) {

	var r uint8
	var err error

	bitset := make([]byte, 32) // 256 bits
	_, err = rand.Read(bitset)
	if err != nil {
		fmt.Println("rand.Read() error:", err)
		return qname, err
	}

	result := ""
	for i := 0; i < len(qname); i++ {
		r = (bitset[i/8] >> uint(i%8)) & 0x1
		if r == 1 && isAlpha(qname[i]) {
			result += string(qname[i] ^ 0x20)
		} else {
			result += string(qname[i])
		}
	}

	return result, err
}

//
// addressString() - return address:port string
//
func addressString(addr string, port int) string {
	if strings.Index(addr, ":") == -1 {
		return addr + ":" + strconv.Itoa(port)
	}
	return "[" + addr + "]" + ":" + strconv.Itoa(port)
}
