package importer

import (
	"io"
	"regexp"

	"github.com/minkezhang/archive-pipeline/import/contact"
	"github.com/minkezhang/archive-pipeline/import/googlevoice/transaction"
	"github.com/minkezhang/archive-pipeline/import/record"
)

var (
	textRE = regexp.MustCompile(
		`<div class="message"><abbr class="dt" title="(?P<date>[\d-T:.]*)">[^<>]*</abbr>:\n<cite class="sender vcard"><a class="tel" href="tel:(?P<address>[+\d]+)">(?:<abbr class="fn" title="">)?(?:<span class="fn">)?(?P<name>[\w, ]+)(?:</span>)?(?:</abbr>)?</a></cite>:\n<q>(?P<message>[^<>]*)</q>\n</div>`,
	)
)

func findSubstringMap(re *regexp.Regexp, match []string) map[string]string {
	results := map[string]string{}
	for i, m := range match {
		results[re.SubexpNames()[i]] = m
	}
	return results
}

type I struct {
	files map[string]io.Reader
}

func New(fs map[string]io.Reader) *I {
	return &I{
		files: fs,
	}
}

func (i *I) Import() (*record.R, error) {
	r := record.New(nil, nil)

	for fn, f := range i.files {
		ok, err := regexp.Match(".*html", []byte(fn))
		if err != nil {
			return nil, err
		}
		if ok {
			d, err := io.ReadAll(f)
			rec, err := i.processText(d)
			if err != nil {
				return nil, err
			}
			r.Merge(rec)
		}
	}

	return r, nil
}

func (i *I) processText(d []byte) (*record.R, error) {
	addressHash := map[string]string{}
	var txs []record.T

	for _, m := range textRE.FindAllStringSubmatch(string(d), -1) {
		match := findSubstringMap(textRE, m)
		addressHash[match["address"]] = match["name"]
		tx, err := transaction.New(match["date"], match["address"], match["message"])
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	var addresses []string
	var contacts []*contact.C
	for a, name := range addressHash {
		addresses = append(addresses, a)
		contacts = append(contacts, contact.New(name, []string{a}))
	}

	for _, tx := range txs {
		tx.(*transaction.T).AddParticipants(addresses)
	}

	return record.New(contacts, txs), nil
}
