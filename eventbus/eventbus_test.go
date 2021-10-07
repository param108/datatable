package eventbus

import (
	"testing"
	"time"

	"context"
	"github.com/param108/datatable/messages"
	"github.com/stretchr/testify/assert"
)

func TestEventBus(t *testing.T) {
	ctx, cancelfn := context.WithCancel(context.Background())
	eb := NewEventBus(ctx)

	defer func() {
		cancelfn()
		eb.Wait()
	}()

	clt1rd, clt1wr := eb.RegisterWindow()

	clt2rd, clt2wr := eb.RegisterWindow()

	t.Run("write from clt2 is available to all others", func(t *testing.T) {
		msg := &messages.Message{
			Key: "test_message",
		}

		timeout := time.NewTimer(time.Second * 10)
		clt2wr <- msg
		select {
		case recvd := <-clt1rd:
			timeout.Stop()
			assert.Equal(t, "test_message", string(recvd.Key))
		case <-timeout.C:
			assert.False(t, true, "Timeout when waiting for message")
		}

		select {
		case recvd := <-clt2rd:
			timeout.Stop()
			assert.Equal(t, "test_message", string(recvd.Key))
		case <-timeout.C:
			assert.False(t, true, "Timeout when waiting for message")
		}

	})

	t.Run("write from clt2 is available to all others", func(t *testing.T) {
		msg := &messages.Message{
			Key: "test_message",
		}

		timeout := time.NewTimer(time.Second * 10)
		clt1wr <- msg
		select {
		case recvd := <-clt1rd:
			timeout.Stop()
			assert.Equal(t, "test_message", string(recvd.Key))
		case <-timeout.C:
			assert.False(t, true, "Timeout when waiting for message")
		}

		select {
		case recvd := <-clt2rd:
			timeout.Stop()
			assert.Equal(t, "test_message", string(recvd.Key))
		case <-timeout.C:
			assert.False(t, true, "Timeout when waiting for message")
		}

	})

}
