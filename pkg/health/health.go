package health

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
)

type Service struct {
	rate         *time.Ticker
	quitChan     chan struct{}
	pingStatus   *atomic.Bool
	dependencies []Dependency
	h            healthcheck.Handler
}

type Dependency interface {
	Ping(context.Context) error
}

func NewHealthHandler(dependencies ...Dependency) (*Service, error) {
	handler := &Service{
		rate:         time.NewTicker(time.Second * 10),
		dependencies: dependencies,
		quitChan:     make(chan struct{}),
		pingStatus:   &atomic.Bool{},
	}

	if err := handler.setupHandlers(); err != nil {
		return nil, err
	}
	return handler, nil
}

func (s *Service) LiveHandler() (string, http.HandlerFunc) {
	return "/live", s.h.LiveEndpoint
}

func (s *Service) HealthHandler() (string, http.HandlerFunc) {
	return "/health", s.h.LiveEndpoint
}

func (s *Service) ReadyHandler() (string, http.HandlerFunc) {
	return "/ready", s.h.ReadyEndpoint
}

func (s *Service) status() bool {
	status := s.pingStatus.Load()
	return status
}

func (s *Service) Run(ctx context.Context) {
	go func() {
		defer s.rate.Stop()
		for {
			select {
			case <-s.rate.C:
				s.pingStatus.Store(handlePing(ctx, s.dependencies...))
			case <-ctx.Done():
				return
			}
		}
	}()
}

func handlePing(ctx context.Context, deps ...Dependency) bool {
	// Get logger
	pings := true
	for _, d := range deps {
		err := d.Ping(ctx)
		if err != nil {
			log.Printf("dependency %s; err=%v", reflect.TypeOf(d).Elem().String(), fmt.Errorf("dependency check failed: %v", err))
		}
		pings = pings && err == nil
	}
	return pings
}

func (s *Service) setupHandlers() error {
	// Get logger
	runtimeGoRoutineCountHealth := func() error {
		if err := healthcheck.GoroutineCountCheck(1000)(); err != nil {
			log.Fatalf("goroutine count failed: %v", err)
			return err
		}
		return nil
	}

	runtimeGCMaxPauseHealth := func() error {
		if err := healthcheck.GCMaxPauseCheck(100 * time.Millisecond)(); err != nil {
			log.Fatalf("runtime gc max pause failed: %v", err)
			return err
		}
		return nil
	}

	dependenciesCheck := func() error {
		if !s.pingStatus.Load() {
			return errors.New("application has one or more failed dependencies")
		}
		return nil
	}

	h := healthcheck.NewHandler()
	h.AddLivenessCheck("goroutines", runtimeGoRoutineCountHealth)
	h.AddLivenessCheck("gcMaxPause", runtimeGCMaxPauseHealth)
	h.AddLivenessCheck("dependencies", dependenciesCheck)
	h.AddReadinessCheck("dependencies", dependenciesCheck)
	s.h = h
	return nil
}

func (s *Service) ServeHealthcheckMux() (*http.ServeMux, error) {
	livePath, liveHandler := s.LiveHandler()
	readyPath, readyHandler := s.ReadyHandler()
	// Set up healthcheck routes
	mux := http.NewServeMux()
	mux.Handle(livePath, liveHandler)
	mux.Handle(readyPath, readyHandler)

	return mux, nil
}

func OSSignalWatcher(quit func()) {
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		c := <-osSigChan
		log.Printf("caught os signal %v", c.String())
		quit()
		time.Sleep(5 * time.Second)
		signal.Stop(osSigChan)
		close(osSigChan)
	}()
}

func ListenForContextCancel(ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Printf("context cancelled, shutting down...")
	}()
}
