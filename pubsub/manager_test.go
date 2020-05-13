package pubsub

import (
	"context"
	"testing"
)

func setup(t *testing.T) (Manager, context.Context) {
	t.Parallel()
	m := New()
	return m, context.Background()
}
func TestPublishAndSusbcribe(t *testing.T) {
	m, ctx := setup(t)
	message := "test pubsub"
	s := m.Subscribe(ctx, "subs1")
	m.Publish(ctx, message)

	msgChan := <-s
	if msgChan != message {
		t.Errorf("error: expected [%s] got [%s]", message, msgChan)
	}

}

func TestResolveAllMessage(t *testing.T) {
	m, ctx := setup(t)
	m.Publish(ctx, "test")
	m.Publish(ctx, "test2")

	msgs := m.ResolveAllMessages(ctx)
	if len(msgs) < 2 {
		t.Errorf("error: expected two messages got %d messages", len(msgs))
	}
}
