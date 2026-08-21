package poll

import (
	"context"
	"testing"
)

func TestNewAggregationContext(t *testing.T) {
	t.Run("second call for the same message cancels the first context", func(t *testing.T) {
		first := NewAggregationContext("channel-1", "message-1")
		second := NewAggregationContext("channel-1", "message-1")

		select {
		case <-first.Done():
		default:
			t.Fatal("expected the first context to be canceled once a second one is created for the same message")
		}
		if err := first.Err(); err != context.Canceled {
			t.Errorf("first.Err() = %v, want %v", err, context.Canceled)
		}

		select {
		case <-second.Done():
			t.Fatal("expected the second (latest) context to still be active")
		default:
		}
	})

	t.Run("different messages get independent contexts", func(t *testing.T) {
		a := NewAggregationContext("channel-2", "message-a")
		b := NewAggregationContext("channel-2", "message-b")

		select {
		case <-a.Done():
			t.Fatal("expected context for message-a to remain active")
		default:
		}
		select {
		case <-b.Done():
			t.Fatal("expected context for message-b to remain active")
		default:
		}
	})
}
