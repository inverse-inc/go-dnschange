package windows

import (
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/sys/windows/registry"
)

// Interface is an injectable interface for running commands.
type Interface interface {
	GetInterfaces() ([]NetworkInterface, error)
	SetInterfaceDNSConfig(NetworkInterface)
	SetDNSServer(dns string, domains []string, peers []string) error
	ResetDNSServer() error
	ReturnDNS() []string
	ReturnDomainSearch() []string
}

// New returns a new Interface
func New() Interface {

	runner := &runner{}

	return runner
}

// runner implements Interface
type runner struct {
	InterFaceDNSConfig NetworkInterface
}

// NetworkInterface structure
type NetworkInterface struct {
	Name           string
	Description    string
	DhcpEnabled    bool
	Domain         string
	IPAddress      net.IP
	Mask           net.IPMask
	DefaultGateway []net.IP
	GatewayMetric  int
	DNSServers     []net.IP
}

func (runner *runner) GetInterfaces() ([]NetworkInterface, error) {
	NetworkInterfaces := []NetworkInterface{}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\`, registry.QUERY_VALUE)
	devices, err := k.ReadSubKeyNames(20)

	for _, device := range devices {
		NetInterface := &NetworkInterface{}
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\`+device, registry.QUERY_VALUE)
		if err != nil {
			log.Println(err)
			continue
		}

		NetInterface.Description = device

		NetInterface.Name = device
		if err != nil {
			log.Println(err)
			continue
		}
		defer k.Close()
		s, _, err := k.GetIntegerValue("EnableDHCP")
		if err != nil {
			continue
		}
		if s == uint64(1) {
			NetInterface.DhcpEnabled = true
			s, _, err := k.GetStringsValue("DhcpDefaultGateway")
			if err != nil {
				log.Println(err)
			} else {
				for _, v := range s {
					NetInterface.DefaultGateway = append(NetInterface.DefaultGateway, net.ParseIP(v))
				}
			}
			t, _, err := k.GetStringValue("DhcpDomain")
			if err != nil {
				log.Println(err)
			} else {
				NetInterface.Domain = t
			}
			t, _, err = k.GetStringValue("DhcpIPAddress")
			if err != nil {
				log.Println(err)
			} else {
				NetInterface.IPAddress = net.ParseIP(t)
			}
			t, _, err = k.GetStringValue("NameServer")
			if err != nil || len(t) == 0 {
				t, _, err = k.GetStringValue("DhcpNameServer")
				if err != nil {
					log.Println(err)
				} else {
					dns := strings.Split(t, ",")
					for _, v := range dns {
						NetInterface.DNSServers = append(NetInterface.DNSServers, net.ParseIP(v))
					}
				}
			} else {
				dns := strings.Split(t, ",")
				for _, v := range dns {
					NetInterface.DNSServers = append(NetInterface.DNSServers, net.ParseIP(v))
				}
			}
			t, _, err = k.GetStringValue("DhcpSubnetMask")
			if err != nil {
				log.Println(err)
			} else {
				NetInterface.Mask = net.IPMask(net.ParseIP(t).To4())
			}

		} else {
			NetInterface.DhcpEnabled = true
			s, _, err := k.GetStringsValue("DefaultGateway")
			if err != nil {
				log.Println(err)
			} else {
				for _, v := range s {
					NetInterface.DefaultGateway = append(NetInterface.DefaultGateway, net.ParseIP(v))
				}
			}
			t, _, err := k.GetStringValue("Domain")
			if err != nil {
				log.Println(err)
			} else {
				NetInterface.Domain = t
			}
			t, _, err = k.GetStringValue("IPAddress")
			if err != nil {
				log.Println(err)
			} else {
				NetInterface.IPAddress = net.ParseIP(t)
			}
			t, _, err = k.GetStringValue("NameServer")
			if err != nil {
				log.Println(err)

			} else {
				dns := strings.Split(t, ",")
				for _, v := range dns {
					NetInterface.DNSServers = append(NetInterface.DNSServers, net.ParseIP(v))
				}
			}
			t, _, err = k.GetStringValue("SubnetMask")
			if err != nil {
				log.Println(err)
			} else {
				NetInterface.Mask = net.IPMask(net.ParseIP(t).To4())
			}

		}
		if err != nil {
			log.Println(err)
		}
		NetworkInterfaces = append(NetworkInterfaces, *NetInterface)
	}
	return NetworkInterfaces, err
}

func (runner *runner) SetInterfaceDNSConfig(Int NetworkInterface) {
	runner.InterFaceDNSConfig = Int
}

func AddNRPT(dns string, domain string) error {
	uuidWithHyphen := uuid.New()
	r, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `Computer\HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`+uuidWithHyphen.String(), registry.ALL_ACCESS)
	if err == nil {
		err = r.SetStringValue("Comment", "PacketFence ZTN "+domain)
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("DisplayName", "ZTN "+domain)
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("GenericDNSservers", dns)
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("IPSECCARestriction", "ZTN")
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("Name", domain)
		if err != nil {
			log.Println(err)
		}
		err = r.SetDWordValue("ConfigOptions", uint32(8))
		if err != nil {
			log.Println(err)
		}
		err = r.SetDWordValue("Version", uint32(2))
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func DelNRPT(dns string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`, registry.QUERY_VALUE)
	nrptrules, err := k.ReadSubKeyNames(20)
	if err != nil {
		log.Println(err)
	}
	for _, nrptrule := range nrptrules {
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`+nrptrule, registry.QUERY_VALUE)
		if err != nil {
			log.Println(err)
		}
		s, _, err := k.GetStringValue("Name")
		if err != nil {
			log.Println(err)
		}
		if s == dns {
			err = registry.DeleteKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`+nrptrule)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return err
}
func (runner *runner) SetDNSServer(dns string, domains []string, peers []string) error {
	var err error
	for _, v := range domains {
		err = AddNRPT(dns, v)
		err = AddNRPT(dns, "."+v)
	}
	for _, v := range peers {
		err = AddNRPT(dns, v)
		for _, searchDomain := range runner.ReturnDomainSearch() {
			err = AddNRPT(dns, v+"."+searchDomain)
		}
	}

	return err
}

func (runner *runner) ResetDNSServer() error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\`+runner.InterFaceDNSConfig.Name, registry.SET_VALUE)
	if err != nil {
		log.Println(err)
	}
	defer k.Close()
	var dnsServers []string
	for _, v := range runner.InterFaceDNSConfig.DNSServers {
		dnsServers = append(dnsServers, v.String())
	}
	dnsservers := strings.Join(dnsServers, ",")
	err = k.SetStringValue("NameServer", dnsservers)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (runner *runner) ReturnDNS() []string {
	var dnsServers []string
	for _, v := range runner.InterFaceDNSConfig.DNSServers {
		dnsServers = append(dnsServers, v.String())
	}
	return dnsServers
}

func (runner *runner) ReturnDomainSearch() []string {
	var searchDomain []string
	searchDomain = append(searchDomain, runner.InterFaceDNSConfig.Domain)
	return searchDomain
}
