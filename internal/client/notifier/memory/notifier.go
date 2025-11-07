package memory

import (
	"context"
	"log"

	"github.com/w-h-a/demo-go/internal/client/notifier"
)

type memoryNotifier struct{}

func (n *memoryNotifier) Notify(ctx context.Context, id string, dest string, opts ...notifier.NotifyOption) error {
	log.Printf("[MemoryNotifier] Sending welcome email for %s at %s\n", id, dest)

	return nil
}

func NewNotifier(opts ...notifier.Option) notifier.Notifier {
	return &memoryNotifier{}
}
