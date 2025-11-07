package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/w-h-a/demo-go/internal/client/notifier"
	memorynotifier "github.com/w-h-a/demo-go/internal/client/notifier/memory"
	userrepo "github.com/w-h-a/demo-go/internal/client/user_repo"
	"github.com/w-h-a/demo-go/internal/client/user_repo/postgres"
	userhttphandler "github.com/w-h-a/demo-go/internal/handler/http/user"
	authhttpmiddleware "github.com/w-h-a/demo-go/internal/middleware/http/auth"
	"github.com/w-h-a/demo-go/internal/server"
	httpserver "github.com/w-h-a/demo-go/internal/server/http"
	"github.com/w-h-a/demo-go/internal/service/user"
)

func Run(ctx *cli.Context) error {
	// 1. config
	if err := godotenv.Load(); err != nil {
		return err
	}

	httpAddr := os.Getenv("HTTP_ADDRESS")
	if len(httpAddr) == 0 {
		httpAddr = ":4000"
	}

	dataLocation := os.Getenv("DATA_LOCATION")
	if len(dataLocation) == 0 {
		return fmt.Errorf("DATA_LOCATION env var is not set")
	}

	// 2. resource

	// 3. logs

	// 4. traces

	// 5. stop channels
	stopChannels := map[string]chan struct{}{}

	// 6. create clients
	ur, err := InitUserRepo(dataLocation)
	if err != nil {
		return err
	}

	n, err := InitNotifier()
	if err != nil {
		return err
	}

	// 7. create services
	userService := user.New(ur, n)
	stopChannels["user"] = make(chan struct{})

	// 8. create servers
	httpSrv, err := InitHttpServer(httpAddr, userService)
	if err != nil {
		return err
	}
	stopChannels["httpserver"] = make(chan struct{})

	// 9. wait group and chans (for graceful shutdown)
	var wg sync.WaitGroup
	errCh := make(chan error, len(stopChannels))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 10. run
	wg.Add(1)
	go func() {
		defer wg.Done()
		// log
		errCh <- userService.Run(stopChannels["user"])
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// log
		errCh <- httpSrv.Run(stopChannels["httpserver"])
	}()

	// 11. block until shutdown
	select {
	case err := <-errCh:
		if err != nil {
			// log that we failed
			return err
		}
	case <-sigChan:
		for _, stop := range stopChannels {
			close(stop)
		}
	}

	wg.Wait()

	// 12. check for shutdown errors
	close(errCh)
	for err := range errCh {
		if err != nil {
			// log
		}
	}

	return nil
}

func InitUserRepo(loc string) (userrepo.UserRepo, error) {
	return postgres.NewUserRepo(
		userrepo.WithLocation(loc),
	), nil
}

func InitNotifier() (notifier.Notifier, error) {
	return memorynotifier.NewNotifier(), nil
}

func InitHttpServer(httpAddr string, userService *user.Service) (server.Server, error) {
	srv := httpserver.NewServer(
		server.WithAddress(httpAddr),
		httpserver.WithMiddleware(
			authhttpmiddleware.New(),
		),
	)

	router := mux.NewRouter()

	usersHandler := userhttphandler.New(userService)

	router.HandleFunc("/api/users", usersHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/users/{id}", usersHandler.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/api/users", usersHandler.GetAllUsers).Methods(http.MethodGet)

	// TODO: additional routes

	if err := srv.Handle(router); err != nil {
		return nil, fmt.Errorf("failed to attach root handler: %w", err)
	}

	return srv, nil
}
