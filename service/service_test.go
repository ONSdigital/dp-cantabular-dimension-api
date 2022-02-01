package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ONSdigital/dp-cantabular-dimension-api/config"
	"github.com/ONSdigital/dp-cantabular-dimension-api/service/mock"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "12"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	errServer     = errors.New("HTTP Server error")
)

var errHealthcheck = errors.New("could not get healthcheck")

func TestInit(t *testing.T) {

	Convey("Having a set of mocked dependencies", t, func() {

		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error{return nil},
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() {},
		}

		GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
			return hcMock, nil
		}

		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				return nil
			},
			ShutdownFunc: func(ctx context.Context) error {
				return nil
			},
		}
		GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
			return serverMock
		}

		svc := &Service{}

		Convey("Given that initialising healthcheck returns an error", func() {
			GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
				return nil, errHealthcheck
			}
			// setup (run before each `Convey` at this scope / indentation):
			svc := New()
			err := svc.Init(ctx,testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init fails with an error", func() {
				So(errors.Is(err, errHealthcheck), ShouldBeTrue)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			// setup (run before each `Convey` at this scope / indentation):
			err := svc.Init(ctx,testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init succeeds", func() {
				So(err, ShouldBeNil)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

	})
}

func TestClose(t *testing.T) {

	Convey("Having a correctly initialised service", t, func() {

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error{return nil},
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}
		GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
			return hcMock, nil
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("server stopped before healthcheck")
				}
				return nil
			},
		}

		GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
			return serverMock
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			svcErrors := make(chan error, 1)
			svc := New()
			err := svc.Init(ctx, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			svc.Start(context.Background(), svcErrors)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("failed to stop http server")
				},
			}

			GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
				return failingserverMock
			}

			svcErrors := make(chan error, 1)
			svc := New()
			err := svc.Init(ctx, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			svc.Start(context.Background(), svcErrors)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			timeoutServerMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(2 * time.Millisecond)
					return nil
				},
			}

			svc := Service{
				config:      cfg,
				server:      timeoutServerMock,
				healthCheck: hcMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(errors.Is(err, context.DeadlineExceeded), ShouldBeTrue)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
