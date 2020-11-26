package godnschange

import (
	"fmt"

	"github.com/inverse-inc/go-dnschange/windows"
	"github.com/jackpal/gateway"
)

func (d *DNSStruct) Change(dns string) {
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
				d.NetInterface.(windows.Interface).SetDNSServer(dns)
			}
		}
	}
}

func (d *DNSStruct) GetDNS() []string {
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
	return d.NetInterface.(windows.Interface).ReturnDNS()
}

func (d *DNSStruct) RestoreDNS(dns string) {
	d.NetInterface.(windows.Interface).ResetDNSServer()
}
