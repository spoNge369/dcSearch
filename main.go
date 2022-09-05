package main

import (
	"fmt"
	"reflect"
	"net"
	"net/netip"
	"github.com/docopt/docopt-go"
	"os"
	r "dcSearch/resolvconf"
)


func dcfind(domain, service string, aip []string) {


	r.Add("/etc/resolv.conf", aip)
	_, addrs, err := net.LookupSRV(service, "tcp", domain)
	
	if err != nil {
		r.Remove("/etc/resolv.conf", aip)
		panic(err)
	}


	for _, addr := range addrs {

		if (addr.Target != "") {
			fmt.Println("FQDN:", addr.Target)

			ips, _ := net.LookupIP(addr.Target)
			for _, ip := range ips {
				fmt.Println("IP DC:", ip.String())
			}
		}
	}

	r.Remove("/etc/resolv.conf", aip)


}

func Hosts(cidr string) ([]netip.Addr, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}

	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}

	if len(ips) < 2 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}



func netip_to_string(hosts []netip.Addr) (ip []string) {
	
	for _, host := range hosts {
		host_string := fmt.Sprintf("%s", host)
		ip = append(ip, host_string)
	}
	return

}

func ip_array(t interface{}) ([]string) {

	var out []string
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)


		for i := 0; i < s.Len(); i++ {
			value := fmt.Sprintf("%s", s.Index(i))
			out = append(out, value)
		}
	}
	return out
}


func main() {
	
	usage := `dcSearch
	
Usage:
  dcSearch Options Arguments
  dcSearch [-s <service> | --service <service>] (-d <domain> | --domain <domain>) (DC ...)
  dcSearch -h | --help
  dcSearch -v | --version

Arguments:
  DC               IP of Possible DCs.

Options:
  -h, --help        Show this screen.
  -v, --version     Show version.
  -d, --domain      AD domain.
  -s, --service     Service (Default: kerberos)

Examples:
  dcSearch -d Enterprise.local 192.168.0.6
  dcSearch -d Enterprise.local 192.168.0.{2..254}
  dcSearch -s ldap -d Enterprise.local 192.168.0.133`

  arg := os.Args[1:]
  arguments, _ := docopt.ParseArgs(usage, arg, "version 0.1")
	domain, _ := arguments.String("<domain>")
	service, _ := arguments.String("<service>")

	if service == "" {
		service = "kerberos"
	}


	i := ip_array(arguments["DC"])

	fmt.Println(service)
	//tmp, _ := Hosts(host)
	//ip_string := netip_to_string(tmp)
	dcfind(domain, service, i)
}
