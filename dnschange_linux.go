package godnschange

import (
	"fmt"
	"os"
)

const (
	resolvConf     = "/etc/resolv.conf"
	resolvConfSave = "/etc/resolv.conf.save"
)

func (d *DNSStruct) Change(dns string) error {

	err := os.Rename(resolvConf, resolvConfSave)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(resolvConf)
	defer f.Close()
	f.WriteString("nameserver " + dns + "\n")
	f.Sync()
	return err
}

func (d *DNSStruct) GetDNS() *DNSInfo {
	InfoDNS := &DNSInfo{}

	var DNS []string

	DNS = append(DNS, resolvConfSave)
	InfoDNS.NameServers = DNS
	InfoDNS.SearchDomain = append(InfoDNS.SearchDomain, "packetfence")
	return InfoDNS
}

func (d *DNSStruct) RestoreDNS(dns string) {
	err := os.Remove(resolvConf)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Rename(resolvConfSave, resolvConf)
	if err != nil {
		fmt.Println(err)
	}
}
