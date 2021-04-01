package record

import (
	"github.com/minkezhang/archive-pipeline/import/contact"

	dpb "github.com/minkezhang/archive-pipeline/api/data_go_proto"
)

type T interface {
	Export() *dpb.Transaction
}

type R struct {
	contacts     map[string]*contact.C
	transactions []T
}

func New(contacts []*contact.C, transactions []T) *R {
	r := &R{
		contacts:     map[string]*contact.C{},
		transactions: transactions,
	}
	r.appendContacts(contacts)
	return r
}

func (r *R) Transactions() []T { return r.transactions }

func (r *R) Contacts() []*contact.C {
	var contacts []*contact.C
	for _, c := range r.contacts {
		contacts = append(contacts, c)
	}
	return contacts
}

func (r *R) Merge(o *R) {
	r.appendContacts(o.Contacts())
	r.transactions = append(r.transactions, o.Transactions()...)
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

func (r *R) appendContacts(cs []*contact.C) {
	for _, oc := range cs {
		if rc, ok := r.contacts[oc.Nickname()]; ok {
			rc.Merge(oc)
		} else {
			r.contacts[oc.Nickname()] = oc
		}
	}
}
