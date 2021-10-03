package eventbus

import (
	"sync"

	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
)

type EventBus struct {
	in     chan *messages.Message
	out    chan *messages.Message
	mx     sync.Mutex
	fanout []chan *messages.Message
	wg     sync.WaitGroup
}

func NewEventBus() *EventBus {

	mgr := &EventBus{
		in:     make(chan *messages.Message),
		out:    make(chan *messages.Message),
		fanout: []chan *messages.Message{},
	}

	go mgr.gofanout()
	go mgr.Boss()
	return mgr
}

func (mgr *EventBus) Boss() {
	for msg := range mgr.in {
		log.Infof("Boss: %s", msg.Key)
		mgr.out <- msg
	}
}

//mgr -> gofanout -> (many)outer -> (many)widgets
//   out       mgr.fanout
func (mgr *EventBus) gofanout() {
	for msg := range mgr.out {
		mgr.mx.Lock()
		for _, ch := range mgr.fanout {
			ch <- msg
		}
		mgr.mx.Unlock()
	}
}

func (mgr *EventBus) outer(fromMgr chan *messages.Message, toClient chan *messages.Message) {
	var (
		to      chan *messages.Message
		msgs    []*messages.Message
		currmsg *messages.Message
	)

	for {
		select {
		case msg := <-fromMgr:
			msgs = append(msgs, msg)
			to = toClient
			currmsg = msgs[0]
		case to <- currmsg:
			log.Infof("Messages %v", msgs)
			msgs = msgs[1:]
			if len(msgs) == 0 {
				// mark the toClient channels out of bound
				// write to nil channel blocks forever
				to = nil
			} else {
				currmsg = msgs[0]
			}
		}
	}
}

// (many)widget -> (many)gofanin -> mgr
//                              (in)
func (mgr *EventBus) gofanin(fromClient chan *messages.Message, toMgr chan *messages.Message) {
	var (
		to      chan *messages.Message
		msgs    []*messages.Message
		currmsg *messages.Message
	)

	for {
		select {
		case msg := <-fromClient:
			msgs = append(msgs, msg)
			to = toMgr
			currmsg = msgs[0]
		case to <- currmsg:
			msgs = msgs[1:]
			if len(msgs) == 0 {
				// mark the toClient channels out of bound
				// write to nil channel blocks forever
				to = nil
			} else {
				currmsg = msgs[0]
			}
		}
	}
}

//mgr -> gofanout -> (many)outer -> (many)widgets
//   out       mgr.fanout
// (many)widget -> (many)gofanin -> mgr
//                              (in)
func (mgr *EventBus) RegisterWindow() (cltrd chan *messages.Message, cltwr chan *messages.Message) {
	cltrd = make(chan *messages.Message)
	cltwr = make(chan *messages.Message)
	fanoutchan := make(chan *messages.Message)

	mgr.mx.Lock()
	mgr.fanout = append(mgr.fanout, fanoutchan)
	mgr.mx.Unlock()

	go mgr.outer(fanoutchan, cltrd)
	go mgr.gofanin(cltwr, mgr.in)

	return cltrd, cltwr
}
