package godnschange

type DNSStruct struct {
	NetInterface interface{}
	Success      bool
}

type DNSInfo struct {
	NameServers  []string
	SearchDomain []string
}

func NewDNSChange() *DNSStruct {
	d := &DNSStruct{}
	return d
}
