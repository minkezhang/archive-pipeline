package contact

import (
	dpb "github.com/minkezhang/archive-pipeline/api/data_go_proto"
)

type C struct {
	nickname  string
	addresses map[string]bool
}

func New(nickname string, addresses []string) *C {
	c := &C{
		nickname:  nickname,
		addresses: map[string]bool{},
	}
	c.appendAddresses(addresses)
	return c
}

func (c *C) Nickname() string { return c.nickname }

func (c *C) Merge(o *C) {
	if c.Nickname() == "" {
		c.nickname = o.Nickname()
	}
	c.appendAddresses(o.Addresses())
}

func (c *C) Export() *dpb.Contact {
	pb := &dpb.Contact{
		Nickname: c.Nickname(),
	}
	for _, a := range c.Addresses() {
		pb.Addresses = append(pb.GetAddresses(), &dpb.Address{Address: a})
	}
	return pb
}

func (c *C) Addresses() []string {
	var addr []string
	for a := range c.addresses {
		addr = append(addr, a)
	}
	return addr
}

func (c *C) appendAddresses(addresses []string) {
	for _, a := range addresses {
		c.addresses[a] = true
	}
}
