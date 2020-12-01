package godnschange

type DNSStruct struct {
	NetInterface interface{}
}

type DNSInfo struct {
	NameServers  []string
	SearchDomain []string
}

func NewDNSChange() *DNSStruct {
	d := &DNSStruct{}
	return d
}
