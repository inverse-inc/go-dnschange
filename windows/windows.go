package windows

import (
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/google/gopacket/pcap"
	"golang.org/x/sys/windows/registry"
)

// Interface is an injectable interface for running scutil/networkconfig/ifconfig commands.
type Interface interface {
	GetInterfaces() ([]NetworkInterface, error)
	SetInterfaceDNSConfig(NetworkInterface)
	SetDNSServer(dns string) error
	ResetDNSServer() error
}

func New() Interface {

	runner := &runner{}

	return runner
}

// runner implements Interface
type runner struct {
	InterFaceDNSConfig NetworkInterface
}

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
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	interfacePattern := regexp.MustCompile("\\{(.*)\\}")

	NetworkInterfaces := []NetworkInterface{}

	for _, device := range devices {
		NetInterface := &NetworkInterface{}
		NetInterface.Name = device.Name
		NetInterface.Description = device.Description
		match := interfacePattern.FindStringSubmatch(strings.ToLower(device.Name))
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\`+match[0], registry.QUERY_VALUE)
		if err != nil {
			log.Println(err)
		}
		defer k.Close()
		s, _, err := k.GetIntegerValue("EnableDHCP")
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
			if err != nil {
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

func (runner *runner) SetDNSServer(dns string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\`+runner.InterFaceDNSConfig.Name, registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
	}
	defer k.Close()
	err = SetStringValue(NameServer, dns)
	if err != nil {
		log.Println(err)
	}
}

func (runner *runner) ResetDNSServer() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\`+runner.InterFaceDNSConfig.Name, registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
	}
	defer k.Close()
	for _, v := range runner.InterFaceDNSConfig.DNSServers {
		dnsServers := append(dnsServers, v.String())
	}
	dnsservers := Join(dnsServers, ",")
	err = SetStringValue(NameServer, dnsservers)
	if err != nil {
		log.Println(err)
	}
}
