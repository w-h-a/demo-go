package unit

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	testmock "github.com/stretchr/testify/mock"
	"github.com/w-h-a/demo-go/api/user"
	mocknotifier "github.com/w-h-a/demo-go/internal/client/notifier/mock"
	userrepo "github.com/w-h-a/demo-go/internal/client/user_repo"
	mockrepo "github.com/w-h-a/demo-go/internal/client/user_repo/mock"
	userservice "github.com/w-h-a/demo-go/internal/service/user"
)

func TestUserService_CreateUser(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	ctx := context.Background()
	dto := user.CreateUserDTO{Name: "Test User", Email: "test@test.com"}
	expectedUser := user.User{ID: "some-uuid", Name: dto.Name, Email: dto.Email}

	t.Run("Success", func(t *testing.T) {
		// Arrange
		var wg sync.WaitGroup
		wg.Add(1)

		mockRepo := mockrepo.NewUserRepo()
		mockNotifier := mocknotifier.NewNotifier()
		userService := userservice.New(mockRepo, mockNotifier)

		mockRepo.On("GetByEmail", ctx, dto.Email).Return(user.User{}, userrepo.ErrUserNotFound)
		mockRepo.On("Create", ctx, dto).Return(expectedUser, nil)
		mockNotifier.On("Notify", testmock.Anything, expectedUser.Name, expectedUser.Email, testmock.Anything).Return(nil).Run(func(args testmock.Arguments) {
			wg.Done()
		})

		// Act
		u, err := userService.CreateUser(ctx, dto)
		wg.Wait()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, u)
		mockRepo.AssertExpectations(t)
		mockNotifier.AssertExpectations(t)
	})
}
