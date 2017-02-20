# godig: Todo List

Fix:  
* ";; * SECTION" etc lines dont' print to stdout.
* Program exit code - set appropriately.

A list of additional features I might implement later.

* Implement test suite
* Use more log.Fatal* and improve its messages
* Should results channel be buffered? and with what capacity?
* For batch processing, provide serial (rather than concurrent) option
* Print out time for the entire batch process.
* Fallback to other servers in resolv.conf if available
* Print more stats: request/response size, answering server IP address
* For cookie, retry query with a newly obtained server cookie
* Implement the following dig options:
  * +[no]all
  * +[no]question
  * +[no]answer
  * +[no]authority
  * +[no]additional
  * +trace
  * +[no]cmd
  * +[no]fail
  * +[no]keepopen
  * +[no]stats
* Implement DNS over TLS
* Implement EDNS chain query option
* Implement DNSSEC validation with configured trust anchor.
