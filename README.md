# godig
A DNS client written in Go.

A DNS client written in Go. Roughly similar to ISC BIND's "dig" and
supports many of the same options. The batchfile mode (+batch=)
uses goroutines to concurrently dispatch queries with a configurable
(+parallel=N) level of parallelism. 

Pre-requisite:
* Miek Gieben's dns package:
  * https://github.com/miekg/dns

#### Installation

[TODO ...]

#### Usage and sample runs

Usage:
```
$ godig -h
godig version 0.3
Usage: godig [<options>] qname [qtype] [qclass]

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
```

Output from sample runs testing some of the options:

```
$ godig
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 60826
;; flags: qr rd ra; QUERY: 1, ANSWER: 13, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;.      IN       NS

;; ANSWER SECTION:
.       298158  IN      NS      i.root-servers.net.
.       298158  IN      NS      f.root-servers.net.
.       298158  IN      NS      g.root-servers.net.
.       298158  IN      NS      m.root-servers.net.
.       298158  IN      NS      k.root-servers.net.
.       298158  IN      NS      d.root-servers.net.
.       298158  IN      NS      l.root-servers.net.
.       298158  IN      NS      j.root-servers.net.
.       298158  IN      NS      a.root-servers.net.
.       298158  IN      NS      e.root-servers.net.
.       298158  IN      NS      c.root-servers.net.
.       298158  IN      NS      b.root-servers.net.
.       298158  IN      NS      h.root-servers.net.

;; ResponseTime: 5.312ms
```

```
$ godig www.google.com
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 1691
;; flags: qr rd ra; QUERY: 1, ANSWER: 5, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;www.google.com.        IN       A

;; ANSWER SECTION:
www.google.com. 263     IN      A       173.194.32.148
www.google.com. 263     IN      A       173.194.32.147
www.google.com. 263     IN      A       173.194.32.146
www.google.com. 263     IN      A       173.194.32.145
www.google.com. 263     IN      A       173.194.32.144

;; ResponseTime: 1.269ms
```

```
$ godig www.huque.com. AAAA
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 42740
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;www.huque.com. IN       AAAA

;; ANSWER SECTION:
www.huque.com.  300     IN      CNAME   cheetara.huque.com.
cheetara.huque.com.     86400   IN      AAAA    2600:3c03:e000:81::a

;; ResponseTime: 0.271ms
```

```
$ godig upenn.edu MX IN
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 52740
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;upenn.edu.     IN       MX

;; ANSWER SECTION:
upenn.edu.      863     IN      MX      20 cluster5a.us.messagelabs.com.
upenn.edu.      863     IN      MX      10 cluster5.us.messagelabs.com.

;; ResponseTime: 0.171ms
```

```
$ godig _kerberos._udp.upenn.edu. SRV
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 24881
;; flags: qr rd ra; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;_kerberos._udp.upenn.edu.      IN       SRV

;; ANSWER SECTION:
_kerberos._udp.upenn.edu.       3564    IN      SRV     20 0 88 kdc2.net.isc.upenn.edu.
_kerberos._udp.upenn.edu.       3564    IN      SRV     20 0 88 kdc3.net.isc.upenn.edu.
_kerberos._udp.upenn.edu.       3564    IN      SRV     10 0 88 kdc1.net.isc.upenn.edu.

;; ResponseTime: 0.145ms
```

```
$ godig @adns1.upenn.edu -6 +norecurse www.upenn.edu A
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 8349
;; flags: qr aa; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;www.upenn.edu. IN       A

;; ANSWER SECTION:
www.upenn.edu.  300     IN      CNAME   www.upenn.edu-dscg.edgesuite.net.

;; ResponseTime: 20.780ms
```

```
$ godig berkeley.edu DNSKEY
UDP response was truncated. Retried over TCP.
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 26905
;; flags: qr rd ra; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;berkeley.edu.	IN	 DNSKEY

;; ANSWER SECTION:
berkeley.edu.	172800	IN	DNSKEY	256 3 10 AwEAAaDAoJ8dQMH9YUD7vvaASkLA+nCdj1lGF0QTwxFZsLk+uwa0MDGUf0g6x1fZRCtHTItqFdeMOc8GE49Q36akwoG2/N2SjIReyeI2yvtGz/v3hnJcD91+6Ub5b6FEJSnU/BTaaCPyx+434mChZ4vcRD89bLOP+qwOH2qtWUEOfIq3
berkeley.edu.	172800	IN	DNSKEY	256 3 10 AwEAAbMkqJ3iI5th1xiG5ptVZoCTNeTqR1fTHSVbTRPSuCXvNPvRNr4fEn9uH/Z2UASF5gTho5dEjabEOq9SpMigviSIHpovTExdNW43hu4/4FCKt9FXz//xu8i0gAsJBcSXBzDFniVgYfjyhxquYD3eKjV7IwTh2f5Sog9F7+5vNHz5
berkeley.edu.	172800	IN	DNSKEY	257 3 10 AwEAAbFmP5ygKmvhsxBnK4LcuMBXqABF8uxXCuxKFvHNNYPQaKG2KXlahnnN8194C8p2wp+T8vwidVt4Gx+O+wbhQ8zHyvoJou2+9yZeAtCr1rYga6l+IEuN6pCYIEBCoiAjMum9uB1o4OW63zDSTWFkSCLmVjVD5rai3TfvQY1M6gMdPxoPiCApD+C+1gd5hkavUvn0y0Y8aAUPShLKW7XgPB14L/Z0QbFfEWcSzrBlFkLS5BjDp1pXx+KgU6CPYL0pgQf4pdPyW6FtcaKA3G8O6xPj7X2T2Ngy0pAAepI3eI8AWpzRNqnmSbglVRCF+eS6QoMuN9teTfViTDrS1VVzMLc=

;; ResponseTime: 0.221ms
```

```
$ godig -x 2607:f470:1001::1:a
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 49357
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;a.0.0.0.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.1.0.0.1.0.7.4.f.7.0.6.2.ip6.arpa.      IN       PTR

;; ANSWER SECTION:
A.0.0.0.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.1.0.0.1.0.7.4.f.7.0.6.2.ip6.arpa.       86365   IN      PTR     adns1.upenn.edu.

;; ResponseTime: 0.513ms
```

```
$ godig +dnssec www.huque.com A
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 22715
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 4, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: do; udp: 4096

;; QUESTION SECTION:
;www.huque.com. IN       A

;; ANSWER SECTION:
www.huque.com.  300     IN      CNAME   cheetara.huque.com.
www.huque.com.  300     IN      RRSIG   CNAME 8 3 300 20160531205924 20160501201010 14703 huque.com. k0+jr+993zCf68FsXcSeyLubVXcMeEG2P4vnvOq+Dj5CoQJ5Ca/qqe3zfReMc3bGuCg9q9wF+VDuHOx78VM3ZeEvs7VvazYa7gFf7FlW9iO5VfPFaV/eeFgUSQw8VWfS8drYsaeeJW0Mn9XVrHi9qv8MkPiJne6XM/yDzWBUHAc=
cheetara.huque.com.     86400   IN      A       50.116.63.23
cheetara.huque.com.     86400   IN      RRSIG   A 8 3 86400 20160604014047 20160505012022 14703 huque.com. aUBP6/ywxtt+9SIYLwG1ZO86fCSmHvK93mZAHGN7qkvot9drfllnpjpvilg6geNyZD2wwNTEvnVWs/bcyXLaQMH8m2wtX1fx8d+x+6/hZrjD35M1KRfYyG8OHk4ad3ly4wZDxtiHXfD7RNX15uVbfxZ+1eMm1dXAib2Nt34Ru9A=

;; ResponseTime: 0.521ms
```

```
$ godig +dnssec _443._tcp.www.ietf.org TLSA
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 22021
;; flags: qr rd ra ad; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: do; udp: 4096

;; QUESTION SECTION:
;_443._tcp.www.ietf.org.        IN       TLSA

;; ANSWER SECTION:
_443._tcp.www.ietf.org. 1765    IN      TLSA    3 1 1 0c72ac70b745ac19998811b131d662c9ac69dbdbe7cb23e5b514b56664c5d3d6
_443._tcp.www.ietf.org. 1765    IN      RRSIG   TLSA 5 5 1800 20170308083405 20160308073501 40452 ietf.org. Gu2ytOVOM9Lx59qfcrBI7MGe3NCdloy+XDDTT0T1FXBJeeuCRMiQVPmrfK9/LbCrOidkEvXoeTs9k/i1IFmVWKD8WhSi7l5TeIE62dKu+wuqkkeZddJ9wZm5cEdy0yVah0hmtzyMTRXOfvABg18l7UgIxb6qZbj8pYWEU8dIO7RqymXTN+EISbwK9gP3G0ngpN3PDbUdjaEx83yPtlCQe5/EIo8yJCrlVD2ijuH4cmiIRye3SW7CVOAc8NF7Kis7LvfJ7xXRVj2Q4Nkkc20jVwt4OLQJPbubYS7erK4i8Xg4YHzSayOYKae49hrGEDCEXCpy3uudAhNw/LqPiXUQcA==

;; ResponseTime: 0.525ms
```

```
$ godig +dnssec www78.upenn.edu A
;; ->>HEADER<<- ;; opcode: QUERY, status: NXDOMAIN, id: 44448
;; flags: qr rd ra ad; QUERY: 1, ANSWER: 0, AUTHORITY: 6, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: do; udp: 4096

;; QUESTION SECTION:
;www78.upenn.edu.       IN       A

;; AUTHORITY SECTION:
upenn.edu.      3565    IN      SOA     assailants.net.isc.upenn.edu. hostmaster.upenn.edu. 1005239355 10800 3600 604800 3600
upenn.edu.      3565    IN      RRSIG   SOA 5 2 3600 20160613154650 20160514144650 50475 upenn.edu. rRLMIEkXwjNVQjHYdgxN3TKjQ0O5P322r0IAMkkbF6BZ7SIBFfWm7nuO6MX5mhbX3xzpGmn8/vB6soCQcSPszSDGQYln0WyPUUSMltyslu/0kvCiOvwZXtaKEiOWStjJRwip8GF5VAuKF/88smOBOkfApV04yD8xauoqcUx8xq8=
upenn.edu.      3565    IN      RRSIG   NSEC 5 2 3600 20160608043139 20160509041707 50475 upenn.edu. OabJHr2Q9UuC/F3WnSfXimYbxpN8A81Wtwz4PO/UjvRTBtW1q9/Enyio+mnLSqKLo3FkV/lSClBv+eblBDfApYLvRJZyTfwLn3mamU3hO54FYQuO5Fe1CYQGbnYPA0VZmNIJEVW3rwcbQDcgfFPfAqG07Z0ZA3qVUg7kXCLxayo=
upenn.edu.      3565    IN      NSEC    _kerberos.upenn.edu. NS SOA MX TXT RRSIG NSEC DNSKEY TYPE65534
www-test1.upenn.edu.    3565    IN      RRSIG   NSEC 5 3 3600 20160528015531 20160428013502 50475 upenn.edu. hZ/iaf4+lQLUz0bMdZ2UnQgV2aKEeCP1W+VGYFVQqldXw/howrv8FOXIhrFN28akqDW9ukpMV+98EJYUhdaFdNW9pcrmdzD85lmEYQV+uJqM8TrppMNE2d4m8Z+2izCqQAZ3/7UVtQiWFwUxaLm6cauUm0OGaMX+GfLOsVWGKg0=
www-test1.upenn.edu.    3565    IN      NSEC    cardpanel.3.WXPN.upenn.edu. CNAME RRSIG NSEC

;; ResponseTime: 0.481ms
```

```
$ godig +cookie www.verisignlabs.com
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 59815
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: ; udp: 4096
; COOKIE: 00ab0535dac3a3b0a0f18f7957392e44a111bb9b8ad96847

;; QUESTION SECTION:
;www.verisignlabs.com.	IN	 A

;; ANSWER SECTION:
www.verisignlabs.com.	3592	IN	CNAME	verisignlabs.com.
verisignlabs.com.	3592	IN	A	72.13.58.64

;; ResponseTime: 0.768ms
```

```
$ godig @8.8.8.8 +subnet=128.91.13.250/24 +cookie www.huque.com. A
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 63606
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: ; udp: 512
; SUBNET: 128.91.13.0/24/0

;; QUESTION SECTION:
;www.huque.com.	IN	 A

;; ANSWER SECTION:
www.huque.com.	299	IN	CNAME	cheetara.huque.com.
cheetara.huque.com.	21599	IN	A	50.116.63.23

;; ResponseTime: 77.168ms
```

```
$ godig -4 @l.root-servers.net. +nsid . SOA
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 43374
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 13, ADDITIONAL: 25

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: ; udp: 4096
; NSID: 69616435322e6c2e726f6f742d736572766572732e6f7267  (i)(a)(d)(5)(2)(.)(l)(.)(r)(o)(o)(t)(-)(s)(e)(r)(v)(e)(r)(s)(.)(o)(r)(g)

;; QUESTION SECTION:
;.      IN       SOA

;; ANSWER SECTION:
.       86400   IN      SOA     a.root-servers.net. nstld.verisign-grs.com. 2016051400 1800 900 604800 86400

;; AUTHORITY SECTION:
.       518400  IN      NS      a.root-servers.net.
.       518400  IN      NS      b.root-servers.net.
.       518400  IN      NS      c.root-servers.net.
.       518400  IN      NS      d.root-servers.net.
.       518400  IN      NS      e.root-servers.net.
.       518400  IN      NS      f.root-servers.net.
.       518400  IN      NS      g.root-servers.net.
.       518400  IN      NS      h.root-servers.net.
.       518400  IN      NS      i.root-servers.net.
.       518400  IN      NS      j.root-servers.net.
.       518400  IN      NS      k.root-servers.net.
.       518400  IN      NS      l.root-servers.net.
.       518400  IN      NS      m.root-servers.net.

;; ADDITIONAL SECTION:
a.root-servers.net.     518400  IN      A       198.41.0.4
b.root-servers.net.     518400  IN      A       192.228.79.201
c.root-servers.net.     518400  IN      A       192.33.4.12
d.root-servers.net.     518400  IN      A       199.7.91.13
e.root-servers.net.     518400  IN      A       192.203.230.10
f.root-servers.net.     518400  IN      A       192.5.5.241
g.root-servers.net.     518400  IN      A       192.112.36.4
h.root-servers.net.     518400  IN      A       198.97.190.53
i.root-servers.net.     518400  IN      A       192.36.148.17
j.root-servers.net.     518400  IN      A       192.58.128.30
k.root-servers.net.     518400  IN      A       193.0.14.129
l.root-servers.net.     518400  IN      A       199.7.83.42
m.root-servers.net.     518400  IN      A       202.12.27.33
a.root-servers.net.     518400  IN      AAAA    2001:503:ba3e::2:30
b.root-servers.net.     518400  IN      AAAA    2001:500:84::b
c.root-servers.net.     518400  IN      AAAA    2001:500:2::c
d.root-servers.net.     518400  IN      AAAA    2001:500:2d::d
f.root-servers.net.     518400  IN      AAAA    2001:500:2f::f
h.root-servers.net.     518400  IN      AAAA    2001:500:1::53
i.root-servers.net.     518400  IN      AAAA    2001:7fe::53
j.root-servers.net.     518400  IN      AAAA    2001:503:c27::2:30
k.root-servers.net.     518400  IN      AAAA    2001:7fd::1
l.root-servers.net.     518400  IN      AAAA    2001:500:9f::42
m.root-servers.net.     518400  IN      AAAA    2001:dc3::35

;; ResponseTime: 7.246ms
```

```
$ godig +0x20 www.upenn.edu A
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 25139
;; flags: qr rd ra; QUERY: 1, ANSWER: 4, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;WwW.UPenN.Edu. IN       A

;; ANSWER SECTION:
www.upenn.edu.  268     IN      CNAME   www.upenn.edu-dscg.edgesuite.net.
www.upenn.edu-dscg.edgesuite.net.       21568   IN      CNAME   a1165.dscg.akamai.net.
a1165.dscg.akamai.net.  20      IN      A       23.220.148.65
a1165.dscg.akamai.net.  20      IN      A       23.220.148.48

;; ResponseTime: 2.420ms
```

```
$ godig txtrecord.huque.com. TXT
UDP response was truncated. Retried over TCP.
;; ->>HEADER<<- ;; opcode: QUERY, status: NOERROR, id: 1742
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;txtrecord.huque.com.   IN       TXT

;; ANSWER SECTION:
txtrecord.huque.com.    86400   IN      TXT     "a TXT record contains a sequence of 1 or more character strings" "a character string is a single length octet followed by that number of characters" "hence, if we exclude the leading length octet, each string can contain a maximum of 255 characters" "the maximum total size of all strings including the length octets is the maximum rdata size of 65,536" "following this string are 2 maximally sized strings" "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"

;; ResponseTime: 0.448ms
```

```
$ godig +edns=1 www.yahoo.com. A
;; ->>HEADER<<- ;; opcode: QUERY, status: BADVERS, id: 53967
;; flags: qr rd; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: ; udp: 4096

;; QUESTION SECTION:
;www.yahoo.com.	IN	 A

;; ResponseTime: 0.230ms
```

```
$ godig +opcode=1 www.verisign.com. A
;; ->>HEADER<<- ;; opcode: IQUERY, status: NOTIMPL, id: 57229
;; flags: qr rd ra; QUERY: 0, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0

;; ResponseTime: 0.173ms
```

