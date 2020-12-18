package godnschange

import (
	"fmt"
	"net"

	"github.com/inverse-inc/go-dnschange/darwin"
	"github.com/jackpal/gateway"
)

func (d *DNSStruct) Change(dns string, domains []string, peers []string, internal string, api string) error {
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
	NetInterface.SetDNSServer(dns, domains, peers, internal, api)

	d.NetInterface = NetInterface
	return err
}

func (d *DNSStruct) GetDNS() *DNSInfo {
	InfoDNS := &DNSInfo{}
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
	InfoDNS.NameServers = d.NetInterface.(darwin.Interface).ReturnDNS()
	InfoDNS.SearchDomain = d.NetInterface.(darwin.Interface).ReturnDomainSearch()
	return InfoDNS
}

func (d *DNSStruct) RestoreDNS(dns string) {
	d.NetInterface.(darwin.Interface).ResetDNSServer(dns)
	// d.NetInterface.(darwin.Interface).ResetDNSServer(d.NetInterface.(darwin.Interface).InterFaceDNSConfig.NameServers[0]
	d.NetInterface.(darwin.Interface).RemoveInterfaceAlias(dns)
}
