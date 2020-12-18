// +build windows
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
	SetDNSServer(dns string, domains []string, peers []string, internal string, api string) error
	ResetDNSServer(dns string) error
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

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`, registry.READ)
	devices, _ := k.ReadSubKeyNames(20)
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

func AddNRPT(dns string, name []string) error {
	uuidWithHyphen := uuid.New()
	r, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\{`+strings.ToUpper(uuidWithHyphen.String())+"}", registry.ALL_ACCESS)

	if err == nil {
		err = r.SetStringValue("Comment", "PacketFence ZTN")
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("DisplayName", "ZTN")
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("GenericDNSservers", dns)
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringValue("IPSECCARestriction", "")
		if err != nil {
			log.Println(err)
		}
		err = r.SetStringsValue("Name", name)
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
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`, registry.READ)
	nrptrules, err := k.ReadSubKeyNames(50)
	if err != nil {
		log.Println(err)
	}

	for _, nrptrule := range nrptrules {
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`+nrptrule, registry.QUERY_VALUE)
		if err != nil {
			log.Println(err)
		}
		s, _, err := k.GetStringValue("Comment")
		if err != nil {
			log.Println(err)
		}
		if s == "PacketFence ZTN" {
			err = registry.DeleteKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Services\Dnscache\Parameters\DnsPolicyConfig\`+nrptrule)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return err
}
func (runner *runner) SetDNSServer(dns string, domains []string, peers []string, internal string, api string) error {
	var err error
	var Name []string
	for _, v := range domains {
		Name = append(Name, v)
		Name = append(Name, "."+v)
	}
	Name = append(Name, internal)
	Name = append(Name, "."+internal)
	for _, v := range peers {
		if v != "" {
			for _, searchDomain := range runner.ReturnDomainSearch() {
				if searchDomain == "" {
					continue
				}
				Name = append(Name, v+"."+searchDomain)
			}
		}
	}
	err = AddNRPT(dns, Name)

	// Forward api fqdn to original dns server

	if net.ParseIP(api) == nil {
		var localDNS string
		var localDNSIP []string
		if len(runner.InterFaceDNSConfig.DNSServers) > 1 {
			for _, v := range runner.InterFaceDNSConfig.DNSServers {
				localDNSIP = append(localDNSIP, v.String())
			}
			localDNS = strings.Join(localDNSIP, ";")
		} else {
			localDNS = runner.InterFaceDNSConfig.DNSServers[0].String()
		}
		err = AddNRPT(localDNS, []string{api})
	}

	return err
}

func (runner *runner) ResetDNSServer(dns string) error {
	err := DelNRPT(dns)
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
