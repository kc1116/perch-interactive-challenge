package core

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/vrecan/death.v3"
	"os"
	SYS "syscall"
)

type EventAggregator struct {
	Registry    *DeviceRegistry
	Threads     int
	StopWorkers chan bool
	MsgQueue    chan *pubsub.Message
	Stop        bool
	sub         *pubsub.Subscription
	Store       *Store
}

type Worker struct {
	MsgQueue chan *pubsub.Message
	Store    *Store
}

func (w *Worker) Work() {
	for msg := range w.MsgQueue {
		w.ProcessMsg(context.Background(), msg)
	}
}

func (w *Worker) ProcessMsg(ctx context.Context, msg *pubsub.Message) {
	interactionEvt := DecodeEvt(string(msg.Data))
	logger.
		WithField("product-name", interactionEvt.GetProductName()).
		WithField("interaction-type", interactionEvt.GetInteractionType().String()).
		WithField("timestamp", interactionEvt.GetTimestamp().String()).
		WithField("productid", interactionEvt.GetProductId()).
		WithField("button-name", interactionEvt.GetButtonName()).
		Infoln("incoming interaction event")

	err := w.Store.PutEvt(interactionEvt)
	if err != nil {
		logger.Errorln("worker failed to store event")
	}

	msg.Ack()
}

func (e *EventAggregator) Start(host, database string) error {
	store, err := NewStore(host, database)
	if err != nil {
		return fmt.Errorf("unable to start event aggregator err when connecting to event store %s", err)
	}

	e.Store = store

	subConf := pubsub.SubscriptionConfig{
		Topic: e.Registry.Topic,
	}

	sub, err := e.Registry.PubSubClient.CreateSubscription(context.Background(), fmt.Sprintf("sub-%s", uuid.NewV1().String()), subConf)
	if err != nil {
		return fmt.Errorf("error creating subscription %s", err)
	}

	e.sub = sub
	go e.StartWorkers(sub, e.Store)

	signalWatcher := death.NewDeath(SYS.SIGINT, SYS.SIGTERM, SYS.SIGKILL, os.Interrupt)

	logger.Infoln("listening for incoming pubsub events . . .")
	err = signalWatcher.WaitForDeath(e)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventAggregator) StartWorkers(sub *pubsub.Subscription, store *Store) {
	logger.Infof("starting event aggregate workers (subscription: %s, num of workers: %x) ", sub.String(), e.Threads)
	for i := 0; i < e.Threads; i++ {
		w := &Worker{MsgQueue: e.MsgQueue, Store: store}
		go w.Work()
	}

	for e.Stop == false {
		err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			e.MsgQueue <- msg
		})
		if err != nil {
			logger.Warnln("error receiving publish", err)
		}
	}

	close(e.StopWorkers)
}

func (e *EventAggregator) Close() error {
	logger.Infof("received stop signal, killing worker threads and exiting gracefully\n")
	close(e.MsgQueue)
	e.Stop = true

	logger.Infof("deleting subscription %s\n", e.sub.String())
	err := e.sub.Delete(context.Background())
	if err != nil {
		logger.Errorf("error deleting subscription: %s %s", e.sub.String(), err)
	}
	return nil
}

func NewEventListener(registry *DeviceRegistry, threads int) *EventAggregator {
	return &EventAggregator{
		Registry:    registry,
		StopWorkers: make(chan bool),
		MsgQueue:    make(chan *pubsub.Message),
		Threads:     threads,
	}
}
