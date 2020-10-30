package godnschange

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
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
				spew.Dump(v)
				// NetInterface.SetInterfaceDNSConfig(v)
				d.NetInterface = NetInterface
				// d.NetInterface.(netsh.Interface).SetDNSServer(dns)
			}
		}
	}
}

func (d *DNSStruct) RestoreDNS(dns string) {
	// d.NetInterface.(netsh.Interface).ResetDNSServer()
}
