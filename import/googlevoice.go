package googlevoice

import (
	"io"
	"regexp"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	cpb "github.com/minkezhang/archive-pipeline/api/constants_go_proto"
	dpb "github.com/minkezhang/archive-pipeline/api/data_go_proto"
)

const (
	timestampLayout = "2006-01-02T15:04:05.999-07:00"
)

var (
	textRE = regexp.MustCompile(
		`<div class="message"><abbr class="dt" title="(?P<date>[\d-T:.]*)">[^<>]*</abbr>:\n<cite class="sender vcard"><a class="tel" href="tel:(?P<address>[+\d]+)">(?:<abbr class="fn" title="">)?(?:<span class="fn">)?(?P<name>[\w, ]+)(?:</span>)?(?:</abbr>)?</a></cite>:\n<q>(?P<message>[^<>]*)</q>\n</div>`,
	)
)

type Record struct {
	contacts     map[string]*Contact
	transactions []*TextTransaction
}

func NewRecord() *Record {
	return &Record{
		contacts: map[string]*Contact{},
	}
}
func (r *Record) AddContact(c *Contact) {
	if tc, ok := r.contacts[c.Nickname]; !ok {
		r.contacts[c.Nickname] = c
	} else {
		tc.AddAddresses(c.Addresses())
	}
}
func (r *Record) AddTransaction(t *TextTransaction) {
	r.transactions = append(r.transactions, t)
}
func (r *Record) Export() *dpb.Record {
	pb := &dpb.Record{}
	for _, c := range r.contacts {
		pb.Contacts = append(pb.GetContacts(), c.Export())
	}
	for _, tx := range r.transactions {
		pb.Transactions = append(pb.GetTransactions(), tx.Export())
	}
	return pb
}

type Contact struct {
	Nickname  string
	addresses map[string]bool
}

func NewContact(nickname string, addresses []string) *Contact {
	c := &Contact{
		Nickname:  nickname,
		addresses: map[string]bool{},
	}
	c.AddAddresses(addresses)
	return c
}
func (c *Contact) AddAddresses(addresses []string) {
	for _, a := range addresses {
		c.addresses[a] = true
	}
}
func (c *Contact) Export() *dpb.Contact {
	pb := &dpb.Contact{
		Nickname: c.Nickname,
	}
	for _, a := range c.Addresses() {
		pb.Addresses = append(pb.GetAddresses(), &dpb.Address{Address: a})
	}
	return pb
}
func (c *Contact) Addresses() []string {
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

func findSubstringMap(re *regexp.Regexp, match []string) map[string]string {
	results := map[string]string{}
	for i, m := range match {
		results[re.SubexpNames()[i]] = m
	}
	return results
}

type I interface {
	Import() (*dpb.Record, error)
}

type IImpl struct {
	files map[string]io.Reader
}

func New(fs map[string]io.Reader) *IImpl {
	return &IImpl{
		files: fs,
	}
}

func (i *IImpl) processText(r *Record, d []byte) error {
	addressHash := map[string]string{}
	var txs []*TextTransaction
	for _, m := range textRE.FindAllStringSubmatch(string(d), -1) {
		match := findSubstringMap(textRE, m)
		addressHash[match["address"]] = match["name"]
		tx, err := NewTextTransaction(match["date"], match["address"], match["message"])
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	var addresses []string
	for a, name := range addressHash {
		addresses = append(addresses, a)
		c := NewContact(name, []string{a})
		r.AddContact(c)
	}

	for _, tx := range txs {
		tx.AddParticipants(addresses)
		r.AddTransaction(tx)
	}
	return nil
}

func (i *IImpl) Import() (*dpb.Record, error) {
	r := NewRecord()
	for fn, f := range i.files {
		ok, err := regexp.Match(".*html", []byte(fn))
		if err != nil {
			return nil, err
		}
		if ok {
			d, err := io.ReadAll(f)
			i.processText(r, d)
			if err != nil {
				return nil, err
			}
		}
	}
	return r.Export(), nil
}
