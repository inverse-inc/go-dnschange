package godnschange

import (
	"fmt"
	"os"
)

const (
	resolvConf     = "/etc/resolv.conf"
	resolvConfSave = "/etc/resolv.conf.save"
)

func (d *DNSStruct) Change(dns string) {

	err := os.Rename(resolvConf, resolvConfSave)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("/etc/resolv.conf")

	f.WriteString("nameserver " + dns + "\n")
	f.Sync()
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
