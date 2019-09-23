package core

import (
	"fmt"
	data "github.com/Pallinder/go-randomdata"
	"github.com/golang/protobuf/ptypes"
	"github.com/kc1116/perch-interactive-challenge/core/protos"
	"github.com/satori/go.uuid"
	"math/rand"
	"sync"
	"time"
)

var shoes = []string{"Ankle", "Athletic", "Boat Shoes", "Boot", "Clogs and Mules", "Crib Shoes", "Firstwalker", "Flat", "Flats", "Heel", "Heels", "Knee High", "Loafers", "Mid-Calf", "Over the Knee", "Oxfords", "Prewalker", "Prewalker Boots", "Slipper Flats", "Slipper Heels", "Sneakers and Athletic Shoes", "SubCategory"}

const (
	minEvtIter = 5
	maxEvtIter = 30
	minSession = 5
	maxSession = 180
)

type publish func(evt *protos.Event) error

type Session struct {
	DeviceID         string
	ID               string
	EventTick        <-chan time.Time
	SessionTimeout   <-chan time.Time
	Done             bool
	PubFn            publish
	Duration         string
	InteractionSleep string
}

func NewSession(deviceID string, pubFn publish) *Session {
	rand.Seed(time.Now().UnixNano())
	tick := time.Second * time.Duration(rand.Intn(maxEvtIter-minEvtIter)+minEvtIter)
	timeout := time.Second * time.Duration(rand.Intn(maxSession-minSession)+minSession)

	return &Session{
		DeviceID:         deviceID,
		ID:               fmt.Sprintf("%-%", data.FirstName(rand.Intn(1-0)+0), uuid.NewV1()),
		EventTick:        time.Tick(tick),
		SessionTimeout:   time.After(timeout),
		Duration:         timeout.String(),
		InteractionSleep: tick.String(),
		PubFn:            pubFn,
	}
}

func (s *Session) Start(wg *sync.WaitGroup) {
	logger.
		WithField("device-id", s.DeviceID).
		WithField("duration", s.Duration).
		WithField("interaction-frequency", fmt.Sprintf("%s", s.InteractionSleep)).
		Infoln("starting session")
	for {
		select {
		case <-s.SessionTimeout:
			return
		case <-s.EventTick:
			err := s.SendEvent()
			if err != nil {
				logger.WithError(err).
					WithField("device-id", s.DeviceID).
					Warnln("error sending event during session")
			}
		}
	}
}

func (s *Session) SendEvent() error {
	evt := RandomEvent()
	logger.
		WithField("product-id", evt.ProductId).
		WithField("product-name", evt.ProductName).
		WithField("button-name", evt.ButtonName).
		WithField("interaction-type", evt.InteractionType.String()).
		Infof("publishing event (device-id: %s)\n", s.DeviceID)

	return s.PubFn(RandomEvent())
}

func RandomInteraction() protos.INTERACTION_TYPE {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(2)
	return protos.INTERACTION_TYPE(n)
}

func RandomShoe() string {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(shoes) - 1)
	return shoes[n]
}

// RandomEvent creates a random event filled with random data
func RandomEvent() *protos.Event {
	evt := &protos.Event{}
	evt.ProductName = RandomShoe()
	evt.InteractionType = RandomInteraction()
	evt.ProductId = uuid.NewV1().String()
	evt.Timestamp = ptypes.TimestampNow()
	if evt.InteractionType == protos.INTERACTION_TYPE_SCREEN_TOUCH {
		evt.ButtonName = fmt.Sprintf("button-%s", data.SillyName())
	}

	return evt
}
