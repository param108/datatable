package manager

import (
	"sync"

	"github.com/param108/datatable/messages"
)

type Manager struct {
	in     chan *messages.Message
	out    chan *messages.Message
	mx     sync.Mutex
	fanout []chan *messages.Message
	wg     sync.WaitGroup
}

func NewManager() *Manager {

	mgr := &Manager{
		in:     make(chan *messages.Message),
		out:    make(chan *messages.Message),
		fanout: []chan *messages.Message{},
	}

	go mgr.gofanout()
	go mgr.Boss()
	return mgr
}

func (mgr *Manager) Boss() {
	for msg := range mgr.in {
		mgr.out <- msg
	}
}

//mgr -> gofanout -> (many)outer -> (many)widgets
//   out       mgr.fanout
func (mgr *Manager) gofanout() {
	for msg := range mgr.out {
		mgr.mx.Lock()
		for _, ch := range mgr.fanout {
			ch <- msg
		}
		mgr.mx.Unlock()
	}
}

func (mgr *Manager) outer(fromMgr chan *messages.Message, toClient chan *messages.Message) {
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
			msgs := msgs[1:]
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
func (mgr *Manager) gofanin(fromClient chan *messages.Message, toMgr chan *messages.Message) {
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
			msgs := msgs[1:]
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
func (mgr *Manager) RegisterWindow() (cltrd chan *messages.Message, cltwr chan *messages.Message) {
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
