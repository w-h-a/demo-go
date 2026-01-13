package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/w-h-a/demo-go/internal/app"
	"github.com/w-h-a/demo-go/internal/service/user"
)

type cli struct {
	Name    string `env:"NAME" default:"gomento"`
	Version string `env:"VERSION" default:"v0.1.0"`

	DataLocation   string `env:"PERSISTER_LOCATION" default:"postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable"`
	HttpServerAddr string `env:"HTTP_SERVER_ADDR" default:":4000"`

	RunAll RunAllCmd `cmd:"" default:"1"`
}

type RunAllCmd struct{}

func (c *RunAllCmd) Run(cli *cli) error {
	// resource

	// logs

	// traces

	// stop channels
	stopChannels := map[string]chan struct{}{}

	// create clients
	ur, err := app.InitUserRepo(cli.DataLocation)
	if err != nil {
		return err
	}

	n, err := app.InitNotifier()
	if err != nil {
		return err
	}

	// create services
	userService := user.New(ur, n)
	stopChannels["user"] = make(chan struct{})

	// create servers
	httpSrv, err := app.InitHttpServer(cli.HttpServerAddr, userService)
	if err != nil {
		return err
	}
	stopChannels["httpserver"] = make(chan struct{})

	// wait group and chans (for graceful shutdown)
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

	// block until shutdown
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

	// check for shutdown errors
	close(errCh)
	for err := range errCh {
		if err != nil {
			// log
		}
	}

	return nil
}

func main() {
	var cli cli
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
