package transaction

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	cpb "github.com/minkezhang/archive-pipeline/api/constants_go_proto"
	dpb "github.com/minkezhang/archive-pipeline/api/data_go_proto"
)

const (
	timestampLayout = "2006-01-02T15:04:05.999-07:00"
)

type T struct {
	timestamp    *timestamppb.Timestamp
	source       string
	participants map[string]bool
	message      string
}

func New(timestamp string, source string, message string) (*T, error) {
	t, err := time.Parse(timestampLayout, timestamp)
	if err != nil {
		return nil, err
	}

	tx := &T{
		timestamp:    timestamppb.New(t),
		source:       source,
		message:      message,
		participants: map[string]bool{},
	}
	tx.AddParticipants([]string{source})

	return tx, nil
}
func (t *T) AddParticipants(ps []string) {
	for _, p := range ps {
		t.participants[p] = true
	}
}
func (t *T) Participants() []string {
	var ps []string
	for p := range t.participants {
		ps = append(ps, p)
	}
	return ps
}
func (t *T) Export() *dpb.Transaction {
	var ps []*dpb.Address
	for p := range t.participants {
		ps = append(ps, &dpb.Address{Address: p})
	}
	return &dpb.Transaction{
		Source:       &dpb.Address{Address: t.source},
		Participants: ps,
		Protocol:     cpb.Protocol_PROTOCOL_SMS,
		Timestamp:    t.timestamp,
		Data:         &dpb.Transaction_TextData{TextData: t.message},
	}
}
