// sd-test-srv runs a simple service discovery system; it returns two hardcoded
// IP addresses for every A query.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.Rcode = dns.RcodeSuccess

	if len(r.Question) != 1 {
		m.Rcode = dns.RcodeServerFailure
		w.WriteMsg(m)
		return
	}

	fmt.Printf("sd-test-srv: got query with question name %q\n", r.Question[0].Name)
	if !strings.HasSuffix(r.Question[0].Name, ".boulder.") {
		m.Rcode = dns.RcodeServerFailure
		w.WriteMsg(m)
		return
	}

	if r.Question[0].Qtype == dns.TypeA {
		hdr := dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		}
		// These two hardcoded IPs correspond to the configured addresses for boulder
		// in docker-compose.yml. In our Docker setup, boulder is present on two
		// networks, rednet and bluenet, with a different IP address on each. This
		// allows us to test load balance across gRPC backends.
		m.Answer = append(m.Answer, &dns.A{
			A:   net.ParseIP("10.77.77.77"),
			Hdr: hdr,
		}, &dns.A{
			A:   net.ParseIP("10.88.88.88"),
			Hdr: hdr,
		})
		w.WriteMsg(m)
		return
	}

	if r.Question[0].Qtype == dns.TypeSRV {
		fmt.Printf("sd-test-srv: Sending SRV record response!\n")
		hdr := dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeSRV,
			Class:  dns.ClassINET,
			Ttl:    0,
		}
		// These two hardcoded IPs correspond to the configured addresses for boulder
		// in docker-compose.yml. In our Docker setup, boulder is present on two
		// networks, rednet and bluenet, with a different IP address on each. This
		// allows us to test load balance across gRPC backends.
		// These two hardcoded names:port combos correspond to the configured names
		// in docker-compose.yml, which in turn point to the local IPs on which our
		// local resolver runs.
		m.Answer = append(m.Answer, &dns.SRV{
			Target: "dns1.boulder",
			Port:   8053,
			Hdr:    hdr,
		}, &dns.SRV{
			Target: "dns2.boulder",
			Port:   8054,
			Hdr:    hdr,
		})
		w.WriteMsg(m)
		return
	}

	// Just return a NOERROR message for non-A, non-SRV questions
	w.WriteMsg(m)
	return
}

func main() {
	listen := flag.String("listen", ":53", "Address and port to listen on.")
	flag.Parse()
	if *listen == "" {
		flag.Usage()
		return
	}
	dns.HandleFunc(".", dnsHandler)
	go func() {
		srv := dns.Server{
			Addr:         *listen,
			Net:          "tcp",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	srv := dns.Server{
		Addr:         *listen,
		Net:          "udp",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
