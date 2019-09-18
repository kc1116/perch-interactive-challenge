package core

import (
	"encoding/base64"
	"fmt"
	data "github.com/Pallinder/go-randomdata"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/kc1116/perch-interactive-challenge/core/protos"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)
//i. Interactions are typically clustered together every 5-30 seconds
// the overall session lasting 5-180 seconds.

//rand.Intn(max - min) + min

const (
	minEvtIter = 5
	maxEvtIter = 30
	minSession = 5
	maxSession = 180
)

type Session struct {
	ID string
	EventTick 	<-chan time.Time
	SessionTimeout <-chan time.Time
	Done bool
}

func NewSession() *Session {
	rand.Seed(time.Now().UnixNano())
	tick := time.Second * time.Duration(rand.Intn(maxEvtIter - minEvtIter) + minEvtIter)
	timeout := time.Second * time.Duration(rand.Intn(maxSession - minSession) + minSession)

	return &Session{
		ID: fmt.Sprintf("%-%", data.FirstName(rand.Intn(1-0) + 0), uuid.NewV1()),
		EventTick: time.Tick(tick),
		SessionTimeout: time.After(timeout),
	}
}

func (s *Session) Start()  {
	for {
		select {
		case <-s.SessionTimeout:
			return
		case <-s.EventTick:
			s.SendEvent()
		}
	}
}

func (s *Session) SendEvent()  {

}

// RandomEvent creates a random event filled with random data
func RandomEvent() *protos.Event {
	evt := &protos.Event{}
	evt.ProductId = uuid.NewV4().String()
	evt.ProductName = data.SillyName()
	evt.Timestamp = ptypes.TimestampNow()

	rand.Seed(time.Now().Unix())
	n := rand.Intn(0 - 2) + 0
	evt.InteractionType = protos.INTERACTION_TYPE(n)

	if evt.InteractionType == protos.INTERACTION_TYPE_SCREEN_TOUCH {
		evt.ButtonName = fmt.Sprintf("button-%s", data.SillyName())
	}

	return evt
}

// EncodeEvent base64 encode event proto bytes
func EncodeEvent(evt *protos.Event) (string, error) {
	b, err := proto.Marshal(evt)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(b), nil
}

// DecodeEvt decode incoming event
func DecodeEvt(encodedEvtStr string) *protos.Event {
	evt := &protos.Event{}
	b, _ := base64.RawStdEncoding.DecodeString(encodedEvtStr)

	_ = proto.Unmarshal(b, evt)
	return evt
}
