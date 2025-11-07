package mock

import (
	"context"

	testmock "github.com/stretchr/testify/mock"
	"github.com/w-h-a/demo-go/api/user"
	userrepo "github.com/w-h-a/demo-go/internal/client/user_repo"
)

type mockUserRepo struct {
	*testmock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, dto user.CreateUserDTO) (user.User, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepo) GetAll(ctx context.Context, opts ...userrepo.GetAllOption) ([]user.User, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]user.User), args.Error(1)
}

func NewUserRepo(opts ...userrepo.Option) *mockUserRepo {
	return &mockUserRepo{&testmock.Mock{}}
}
