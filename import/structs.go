package structs

import (
	"io"
	"regexp"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	cpb "github.com/minkezhang/archive-pipeline/api/constants_go_proto"
	dpb "github.com/minkezhang/archive-pipeline/api/data_go_proto"
)

type R struct {
	contacts     map[string]*C
	transactions []*TextTransaction
}

func NewR() *R {
	return &R{
		contacts: map[string]*C{},
	}
}

func (r *R) AddContact(c *C) {
	if tc, ok := r.contacts[c.Nickname]; !ok {
		r.contacts[c.Nickname] = c
	} else {
		tc.AddAddresses(c.Addresses())
	}
}

func (r *R) AddTransaction(t *TextTransaction) {
	r.transactions = append(r.transactions, t)
}

func (r *R) Export() *dpb.Record {
	pb := &dpb.Record{}
	for _, c := range r.contacts {
		pb.Contacts = append(pb.GetContacts(), c.Export())
	}
	for _, tx := range r.transactions {
		pb.Transactions = append(pb.GetTransactions(), tx.Export())
	}
	return pb
}

type C struct {
	Nickname  string
	addresses map[string]bool
}

func NewC(nickname string, addresses []string) *C {
	c := &C{
		Nickname:  nickname,
		addresses: map[string]bool{},
	}
	c.AddAddresses(addresses)
	return c
}
func (c *C) AddAddresses(addresses []string) {
	for _, a := range addresses {
		c.addresses[a] = true
	}
}
func (c *C) Export() *dpb.C {
	pb := &dpb.C{
		Nickname: c.Nickname,
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

type TextTransaction struct {
	date         *timestamppb.Timestamp
	source       string
	participants map[string]bool
	message      string
}

func NewTextTransaction(date string, source string, message string) (*TextTransaction, error) {
	t, err := time.Parse(timestampLayout, date)
	if err != nil {
		return nil, err
	}

	tx := &TextTransaction{
		date:         timestamppb.New(t),
		source:       source,
		message:      message,
		participants: map[string]bool{},
	}
	tx.AddParticipants([]string{source})

	return tx, nil
}
func (t *TextTransaction) AddParticipants(ps []string) {
	for _, p := range ps {
		t.participants[p] = true
	}
}
func (t *TextTransaction) Participants() []string {
	var ps []string
	for p := range t.participants {
		ps = append(ps, p)
	}
	return ps
}
func (t *TextTransaction) Export() *dpb.Transaction {
	var ps []*dpb.Address
	for p := range t.participants {
		ps = append(ps, &dpb.Address{Address: p})
	}
	return &dpb.Transaction{
		Source:       &dpb.Address{Address: t.source},
		Participants: ps,
		Protocol:     cpb.Protocol_PROTOCOL_SMS,
		Timestamp:    t.date,
		Data:         &dpb.Transaction_TextData{TextData: t.message},
	}
}
