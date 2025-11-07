package notifier

import "context"

type Notifier interface {
	Notify(ctx context.Context, id string, dest string, opts ...NotifyOption) error
}
