package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/kc1116/perch-interactive-challenge/core/protos"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"time"
)

type Interaction struct {
	ID              string `gorethink:"id,omitempty"`
	Timestamp       string `gorethink:"timestamp,omitempty"`
	ProductName     string `gorethink:"productName,omitempty"`
	InteractionType string `gorethink:"interactionType,omitempty"`
}

type Store struct {
	session *r.Session
}

func (s *Store) Init() error {
	err := r.DBCreate("interactions").Exec(s.session)
	err = r.DB("interactions").TableCreate("events").Exec(s.session)
	err = r.DB("interactions").Table("events").IndexCreate("productName").Exec(s.session)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	return nil
}

func (s *Store) PutEvt(evt *protos.Event) error {
	interaction := NewInteraction(evt)
	_, err := r.DB("interactions").Table("events").Insert(interaction).RunWrite(s.session)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	return nil
}

func (s *Store) GetStream() (*r.Cursor, error)  {
	return r.Table("events").Changes().Run(s.session)
}

func NewStore(host, database string) (*Store, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:    host, // "127.0.0.1:28015" default
		Database:   database,
		InitialCap: 10,
		MaxOpen:    10,
	})

	if err != nil {
		logger.Errorf("error creating rethinkdb connection: %s", err)
		return nil, err
	}

	return &Store{session}, nil
}

func NewInteraction(evt *protos.Event) *Interaction {
	t, _ := ptypes.Timestamp(evt.GetTimestamp())

	return &Interaction{
		ID:              evt.GetProductId(),
		Timestamp:       t.Format(time.RFC3339),
		ProductName:     evt.GetProductName(),
		InteractionType: evt.GetInteractionType().String(),
	}
}
