package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/kc1116/perch-interactive-challenge/core/protos"
	"github.com/olahol/melody"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"gopkg.in/vrecan/death.v3"
	"os"
	SYS "syscall"
	"time"
)

type Interaction struct {
	ID              string `gorethink:"identifier,omitempty"`
	Timestamp       string `gorethink:"timestamp,omitempty"`
	ProductName     string `gorethink:"productName,omitempty"`
	InteractionType string `gorethink:"interactionType,omitempty"`
}

type Store struct {
	session *r.Session
	Shutdown chan bool
}

func (s *Store) Init() error {
	_ = r.DBDrop("interactions").Exec(s.session)
	_ = r.DBCreate("interactions").Exec(s.session)
	_ = r.DB("interactions").TableCreate("events").Exec(s.session)
	_ = r.DB("interactions").Table("events").IndexCreate("productName").Exec(s.session)


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

func (s *Store) Close() error {
	logger.Infoln("shutting down socket server gracefully")

	close(s.Shutdown)
	return nil
}

func (s *Store) StartWSProxy()  {
	r := gin.Default()
	m := melody.New()
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	signalWatcher := death.NewDeath(SYS.SIGINT, SYS.SIGTERM, SYS.SIGKILL, os.Interrupt)
	logger.Infoln("starting websocket client proxy . . .")
	go r.Run(":8000")
	go func(m *melody.Melody) {
		interactionCh := make(chan interface{})
		stream, err := s.GetStream()
		if err != nil {
			logger.Fatalln(err)
		}

		stream.Listen(interactionCh)

		for {
			select {
			case <- s.Shutdown:
				s.session.Close()
				m.Close()
				return
			case interaction := <- interactionCh:
				if interaction == nil {
					continue
				}
				data := interaction.(map[string]interface{})
				if val, ok := data["new_val"]; ok {
					b, err := json.Marshal(val)
					if err != nil {
						logger.Errorln("error unmarshalling interaction from chan ", err)
						continue
					}

					err = m.Broadcast(b)
					if err != nil {
						logger.Errorln("error broadcasting interaction to clients ", err)
					}
				}
			}
		}
	}(m)
	err := signalWatcher.WaitForDeath(s)
	if err != nil {
		logger.Fatalln(err)
	}
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

	return &Store{session,make(chan bool)}, nil
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
