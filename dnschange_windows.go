package godnschange

import (
	"fmt"

	"github.com/inverse-inc/go-dnschange/windows"
	"github.com/jackpal/gateway"
)

func (d *DNSStruct) Change(dns string, domains []string, peers []string) error {
	gatewayIP, _ := gateway.DiscoverGateway()
	NetInterface := windows.New()
	NetInterfaces, err := NetInterface.GetInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range NetInterfaces {
		for _, w := range v.DefaultGateway {
			if gatewayIP.String() == w.String() {
				NetInterface.SetInterfaceDNSConfig(v)
				d.NetInterface = NetInterface
				d.NetInterface.(windows.Interface).SetDNSServer(dns, domains, peers)
			}
		}
	}
	return err
}

func (d *DNSStruct) GetDNS() *DNSInfo {
	InfoDNS := &DNSInfo{}
	gatewayIP, _ := gateway.DiscoverGateway()
	NetInterface := windows.New()
	NetInterfaces, err := NetInterface.GetInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range NetInterfaces {
		for _, w := range v.DefaultGateway {
			if gatewayIP.String() == w.String() {
				NetInterface.SetInterfaceDNSConfig(v)
				d.NetInterface = NetInterface
			}
		}
	}
	InfoDNS.NameServers = d.NetInterface.(windows.Interface).ReturnDNS()
	InfoDNS.SearchDomain = d.NetInterface.(windows.Interface).ReturnDomainSearch()
	return InfoDNS
}

func (d *DNSStruct) RestoreDNS(dns string) {
	d.NetInterface.(windows.Interface).ResetDNSServer()
}
