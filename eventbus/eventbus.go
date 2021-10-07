package eventbus

import (
	"context"
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
	ctx    context.Context
}

func NewEventBus(ctx context.Context) *EventBus {

	mgr := &EventBus{
		in:     make(chan *messages.Message),
		out:    make(chan *messages.Message),
		fanout: []chan *messages.Message{},
		ctx:    ctx,
	}

	mgr.wg.Add(1)
	go func() {
		defer mgr.wg.Done()
		mgr.gofanout()
	}()

	mgr.wg.Add(1)
	go func() {
		defer mgr.wg.Done()
		mgr.Boss()
	}()

	return mgr
}

func (mgr *EventBus) Boss() {
	for {
		select {
		case msg := <-mgr.in:
			log.Infof("Boss: %s", msg.Key)
			select {
			case mgr.out <- msg:
			case <-mgr.ctx.Done():
				log.Infof("Exitting Boss: context")
				return
			}
		case <-mgr.ctx.Done():
			log.Infof("Exitting Boss: context")
			return
		}
	}
}

//mgr -> gofanout -> (many)outer -> (many)widgets
//   out       mgr.fanout
func (mgr *EventBus) gofanout() {
	for {
		select {
		case msg := <-mgr.out:
			mgr.mx.Lock()
			for _, ch := range mgr.fanout {
				select {
				case ch <- msg:
				case <-mgr.ctx.Done():
					log.Infof("Exitting gofanout")
					mgr.mx.Unlock()
					return
				}
			}
			mgr.mx.Unlock()
		case <-mgr.ctx.Done():
			log.Infof("Exitting gofanout")
			return
		}
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
		case <-mgr.ctx.Done():
			log.Infof("exitting outer")
			return
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
		case <-mgr.ctx.Done():
			log.Infof("exitting gofanin")
			return

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

	mgr.wg.Add(1)
	go func() {
		defer mgr.wg.Done()
		mgr.outer(fanoutchan, cltrd)
	}()

	mgr.wg.Add(1)
	go func() {
		defer mgr.wg.Done()
		mgr.gofanin(cltwr, mgr.in)
	}()

	return cltrd, cltwr
}

// cancel the context and then call this
func (mgr *EventBus) Wait() {
	mgr.wg.Wait()
}
