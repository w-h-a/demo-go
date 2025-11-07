package mock

import (
	"context"

	testmock "github.com/stretchr/testify/mock"
	"github.com/w-h-a/demo-go/internal/client/notifier"
)

type mockNotifier struct {
	*testmock.Mock
}

func (m *mockNotifier) Notify(ctx context.Context, id string, dest string, opts ...notifier.NotifyOption) error {
	args := m.Called(testmock.Anything, id, dest, opts)
	return args.Error(0)
}

func NewNotifier(opts ...notifier.Option) *mockNotifier {
	return &mockNotifier{&testmock.Mock{}}
}
