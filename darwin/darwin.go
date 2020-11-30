package darwin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	utilexec "k8s.io/utils/exec"
)

const (
	cmdNetworksetup string = "networksetup"
	cmdIfconfig     string = "ifconfig"
	cmdScutil       string = "scutil"
)

// Interface is an injectable interface for running scutil/networkconfig/ifconfig commands.
type Interface interface {
	// GetDNSServers retreive the dns servers
	GetDNSServers(iface string) error
	// Set DNS server
	SetDNSServer(dns string) error
	// Reset DNS server
	ResetDNSServer() error
	AddInterfaceAlias(string) error
	RemoveInterfaceAlias(string) error
	ReturnDNS() []string
}

// runner implements Interface
type runner struct {
	mu                 sync.Mutex
	exec               utilexec.Interface
	InterFaceDNSConfig DNSConfig
}

// DNSConfig structure
type DNSConfig struct {
	Domain       string
	SearchDomain []string
	NameServers  []string
	IfIndex      string
	IfName       string
	Flags        string
	Reach        string
	Options      string
}

// New returns a new Interface which will exec scutil.
func New(exec utilexec.Interface) Interface {

	if exec == nil {
		exec = utilexec.New()
	}

	runner := &runner{
		exec: exec,
	}

	return runner
}

// GetDNSServers uses the show addresses command and returns a formatted structure
func (runner *runner) GetDNSServers(ifname string) error {
	runner.InterFaceDNSConfig = DNSConfig{}

	err := runner.InterfaceAliasName(ifname)

	args := []string{
		"-getdnsservers", runner.InterFaceDNSConfig.IfName,
	}

	output, _ := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput()

	DNSString := string(output[:])

	if strings.Contains("There aren't any DNS Servers set on", DNSString) {
		args := []string{
			"--dns",
		}

		output, _ := runner.exec.Command(cmdScutil, args...).CombinedOutput()

		DNSString := string(output[:])

		outputLines := strings.Split(DNSString, "\n")

		interfacePattern := regexp.MustCompile("^\\d+\\s+\\((.*)\\)")

		runner.InterFaceDNSConfig = DNSConfig{}

		found := false

		for _, outputLine := range outputLines {
			if !found {
				if strings.Contains(outputLine, "DNS configuration (for scoped queries)") {
					found = true
				} else {
					continue
				}
			}

			parts := strings.SplitN(outputLine, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if strings.HasPrefix(key, "if_index") {
				match := interfacePattern.FindStringSubmatch(value)
				if match[1] == ifname {
					found = true
					runner.InterFaceDNSConfig.IfIndex = ifname
				}
			} else if strings.HasPrefix(key, "nameserver") {
				runner.InterFaceDNSConfig.NameServers = append(runner.InterFaceDNSConfig.NameServers, value)
			} else if strings.HasPrefix(key, "search domain") {
				runner.InterFaceDNSConfig.SearchDomain = append(runner.InterFaceDNSConfig.SearchDomain, value)
			} else if strings.HasPrefix(key, "flags") {
				runner.InterFaceDNSConfig.Flags = value
			} else if strings.HasPrefix(key, "reach") {
				runner.InterFaceDNSConfig.Reach = value
			} else if strings.HasPrefix(key, "domain") {
				runner.InterFaceDNSConfig.Domain = value
			} else if strings.HasPrefix(key, "reach") {
				runner.InterFaceDNSConfig.Reach = value
			} else if strings.HasPrefix(key, "options") {
				runner.InterFaceDNSConfig.Options = value
			}
		}

	} else {

		outputLines := strings.Split(DNSString, "\n")

		for _, outputLine := range outputLines {
			runner.InterFaceDNSConfig.NameServers = append(runner.InterFaceDNSConfig.NameServers, outputLine)
		}
	}
	return err
}

func (runner *runner) InterfaceAliasName(ifname string) error {

	args := []string{
		"-listnetworkserviceorder",
	}

	output, _ := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput()

	DNSString := string(output[:])

	outputLines := strings.Split(DNSString, "\n")

	interfacePattern := regexp.MustCompile("\\(Hardware Port:\\s+(.*),\\s+Device:\\s+(.*)\\)")

	err := errors.New("Unable to find the interface alias")

	for _, outputLine := range outputLines {
		if strings.Contains(outputLine, "Hardware Port") {
			match := interfacePattern.FindStringSubmatch(outputLine)
			if match[2] == ifname {
				runner.InterFaceDNSConfig.IfName = match[1]
				err = nil
			}
		} else {
			continue
		}
	}
	return err
}

// Set DNS server
func (runner *runner) SetDNSServer(dns string) error {
	args := []string{
		"-setdnsservers", runner.InterFaceDNSConfig.IfName, dns,
	}
	cmd := strings.Join(args, " ")
	if stdout, err := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set dns servers on [%v], error: %v. cmd: %v. stdout: %v", runner.InterFaceDNSConfig.IfName, err.Error(), cmd, string(stdout))
	}
	return nil
}

// Reset DNS
func (runner *runner) ResetDNSServer() error {
	args := []string{
		"-setdnsservers", runner.InterFaceDNSConfig.IfName, strings.Join(runner.InterFaceDNSConfig.NameServers[:], " "),
	}
	cmd := strings.Join(args, " ")

	if stdout, err := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reset dns servers on [%v], error: %v. cmd: %v. stdout: %v", runner.InterFaceDNSConfig.IfName, err.Error(), cmd, string(stdout))
	}

	return nil
}

// Add interface alias
func (runner *runner) AddInterfaceAlias(ip string) error {
	args := []string{
		"lo0", "alias", ip,
	}
	if _, err := runner.exec.Command(cmdIfconfig, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add alias on interface lo0")
	}
	return nil
}

// Remove interface alias
func (runner *runner) RemoveInterfaceAlias(ip string) error {
	args := []string{
		"lo0", "-alias", ip,
	}
	if _, err := runner.exec.Command(cmdIfconfig, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove alias on interface lo0")
	}
	return nil
}

func (runner *runner) ReturnDNS() []string {
	var dnsServers []string
	for _, v := range runner.InterFaceDNSConfig.NameServers {
		dnsServers = append(dnsServers, v)
	}
	return dnsServers
}
