package steps

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"

	"net/http"

	"github.com/ONSdigital/dp-cantabular-dimension-api/config"
	"github.com/ONSdigital/dp-cantabular-dimension-api/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/maxcnunes/httpfake"

	componenttest "github.com/ONSdigital/dp-component-test"
)

var (
	BuildTime string = "1625046891"
	GitCommit string = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   string = ""
)

type Component struct {
	componenttest.ErrorFeature
	errorChan        chan error
	svc              *service.Service
	cfg              *config.Config
	signals          chan os.Signal
	CantabularSrv    *httpfake.HTTPFake
	CantabularApiExt *httpfake.HTTPFake
	wg               *sync.WaitGroup
	ctx              context.Context
	HTTPServer       *http.Server
	ServiceRunning   bool
}

func NewComponent(t testing.TB, zebedeeURL string) (*Component, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	c := &Component{
		errorChan:        make(chan error),
		CantabularSrv:    httpfake.New(),
		CantabularApiExt: httpfake.New(httpfake.WithTesting(t)),
		wg:               &sync.WaitGroup{},
		ctx:              context.Background(),
		HTTPServer:       &http.Server{},
		cfg:              cfg,
	}

	c.cfg.ZebedeeURL = zebedeeURL
	c.cfg.CantabularURL = c.CantabularSrv.ResolveURL("")
	c.cfg.CantabularExtURL = c.CantabularApiExt.ResolveURL("")

	return c, nil
}

// Init initialises the server and the mocks
func (c *Component) Init() (http.Handler, error) {
	log.Info(c.ctx, "config used by component tests", log.Data{"cfg": c.cfg})

	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	// Create service and initialise it
	c.svc = service.New()
	if err := c.svc.Init(c.ctx, c.cfg, BuildTime, GitCommit, Version); err != nil {
		return nil, fmt.Errorf("failed to initialise service: %w", err)
	}

	return c.HTTPServer.Handler, nil
}

func (c *Component) setInitialiserMock() {
	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		c.HTTPServer.Addr = bindAddr
		c.HTTPServer.Handler = router
		return c.HTTPServer
	}
}

func (c *Component) startService(ctx context.Context) error {
	c.svc.Start(ctx, c.errorChan)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		// blocks until an os interrupt or a fatal error occurs
		select {
		case err := <-c.errorChan:
			if errClose := c.svc.Close(ctx); errClose != nil {
				log.Warn(ctx, "error closing server during error handing", log.Data{"close_error": errClose})
			}
			panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
		case sig := <-c.signals:
			log.Info(ctx, "os signal received", log.Data{"signal": sig})
		}
		if err := c.svc.Close(ctx); err != nil {
			panic(fmt.Errorf("unexpected error during service graceful shutdown: %w", err))
		}
	}()

	return nil
}

// Close kills the application under test.
func (c *Component) Close() {
	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	c.wg.Wait()
}

// Reset runs before each scenario. It re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	c.setInitialiserMock()

	if _, err := c.Init(); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	c.CantabularSrv.Reset()
	c.CantabularApiExt.Reset()

	return nil
}
