package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
