package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/w-h-a/demo-go/api/user"
	"github.com/w-h-a/demo-go/internal/client/notifier"
	userrepo "github.com/w-h-a/demo-go/internal/client/user_repo"
)

type Service struct {
	repo      userrepo.UserRepo
	notifier  notifier.Notifier
	isRunning bool
	mtx       sync.RWMutex
}

func (s *Service) Run(stop chan struct{}) error {
	s.mtx.RLock()
	if s.isRunning {
		s.mtx.RUnlock()
		return errors.New("user service already running")
	}
	s.mtx.RUnlock()

	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to start user service: %w", err)
	}

	<-stop

	return s.Stop()
}

func (s *Service) Start() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.isRunning {
		return errors.New("user service already started")
	}

	s.isRunning = true

	return nil
}

func (s *Service) Stop() error {
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	return s.stop(stopCtx)
}

func (s *Service) stop(ctx context.Context) error {
	s.mtx.Lock()

	if !s.isRunning {
		s.mtx.Unlock()
		return errors.New("user service not running")
	}

	s.isRunning = false

	s.mtx.Unlock()

	gracefulStopDone := make(chan struct{})
	go func() {
		// TODO: close clients gracefully
		close(gracefulStopDone)
	}()

	var stopErr error

	select {
	case <-gracefulStopDone:
	case <-ctx.Done():
		stopErr = ctx.Err()
	}

	return stopErr
}

// CreateUser contains the business logic for creating a new user.
func (s *Service) CreateUser(ctx context.Context, dto user.CreateUserDTO) (user.User, error) {
	// 1. Business Logic: Validation
	if dto.Name == "" || dto.Email == "" {
		return user.User{}, ErrInvalidInput
	}

	dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
	dto.Name = strings.TrimSpace(dto.Name)

	// 2. Business Logic: Check for duplicates
	_, err := s.repo.GetByEmail(ctx, dto.Email)
	if err == nil {
		return user.User{}, ErrEmailInUse
	}
	if !errors.Is(err, userrepo.ErrUserNotFound) {
		return user.User{}, err
	}

	// 3. Call the repository to create the user
	u, err := s.repo.Create(ctx, dto)
	if err != nil {
		return user.User{}, err
	}

	// 4. Orchestration: Send a welcome email (fire-and-forget)
	// We run this in a goroutine so it doesn't block the HTTP response.
	// We also create a new background context in case the original request is cancelled.
	go func() {
		if err := s.notifier.Notify(context.Background(), u.Name, u.Email); err != nil {
			// In a real app, we'd log this to a proper monitoring service.
			log.Printf("Error sending welcome email: %v\n", err)
		}
	}()

	return u, nil
}

// GetUser is the business logic for retrieving a single user.
func (s *Service) GetUser(ctx context.Context, id string) (user.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllUsers is the business logic for retrieving all users.
func (s *Service) GetAllUsers(ctx context.Context) ([]user.User, error) {
	return s.repo.GetAll(ctx)
}

func New(repo userrepo.UserRepo, notifier notifier.Notifier) *Service {
	return &Service{repo, notifier, false, sync.RWMutex{}}
}
