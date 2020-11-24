package godnschange

import (
	"fmt"
	"net"

	"github.com/inverse-inc/go-dnschange/darwin"
	"github.com/jackpal/gateway"
)

func (d *DNSStruct) Change(dns string) {
	gatewayIP, _ := gateway.DiscoverGateway()
	var gatewayInterface string
	Interfaces, _ := net.Interfaces()
	for _, v := range Interfaces {
		eth, _ := net.InterfaceByName(v.Name)
		adresses, _ := eth.Addrs()
		for _, adresse := range adresses {
			_, NetIP, _ := net.ParseCIDR(adresse.String())
			if NetIP.Contains(gatewayIP) {
				gatewayInterface = v.Name
			}
		}
	}
	NetInterface := darwin.New(nil)
	err := NetInterface.GetDNSServers(gatewayInterface)
	if err != nil {
		fmt.Println(err)
	}
	NetInterface.AddInterfaceAlias(dns)
	NetInterface.SetDNSServer(dns)

	d.NetInterface = NetInterface
}

func (d *DNSStruct) GetDNS() []string {
	gatewayIP, _ := gateway.DiscoverGateway()
	var gatewayInterface string
	Interfaces, _ := net.Interfaces()
	for _, v := range Interfaces {
		eth, _ := net.InterfaceByName(v.Name)
		adresses, _ := eth.Addrs()
		for _, adresse := range adresses {
			_, NetIP, _ := net.ParseCIDR(adresse.String())
			if NetIP.Contains(gatewayIP) {
				gatewayInterface = v.Name
			}
		}
	}
	NetInterface := darwin.New(nil)
	err := NetInterface.GetDNSServers(gatewayInterface)
	if err != nil {
		fmt.Println(err)
	}
	d.NetInterface = NetInterface
	return NetInterface.InterFaceDNSConfig.NameServers
}

func (d *DNSStruct) RestoreDNS(dns string) {
	d.NetInterface.(darwin.Interface).ResetDNSServer()
	d.NetInterface.(darwin.Interface).RemoveInterfaceAlias(dns)
}
