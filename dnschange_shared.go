package godnschange

type DNSStruct struct {
	NetInterface interface{}
}

func NewDNSChange() *DNSStruct {
	d := &DNSStruct{}
	return d
}

func (d *DNSStruct) DoNothing() {

}
