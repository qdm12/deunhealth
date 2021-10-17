package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/qdm12/deunhealth/internal/config"
	"github.com/qdm12/deunhealth/internal/docker"
	"github.com/qdm12/deunhealth/internal/health"
	"github.com/qdm12/deunhealth/internal/loop"
	"github.com/qdm12/deunhealth/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/goshutdown"
	"github.com/qdm12/gosplash"
)

var (
	// Values set by the build system.
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	args := os.Args

	logger := logging.New(logging.Settings{})

	configReader := config.NewReader()

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, configReader)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("Caught OS signal, shutting down\n")
		stop()
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck query
			os.Exit(0)
		}
		logger.Error("Fatal error: " + err.Error())
		os.Exit(1)
	}

	err := <-errorCh
	if err != nil {
		logger.Error("shutdown error: " + err.Error())
	}
	os.Exit(1)
}

func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger logging.ParentLogger, configReader config.ConfReader) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()

		config, err := configReader.ReadHealth()
		if err != nil {
			return err
		}

		return client.Query(ctx, config.Address)
	}

	announcementExpiration, err := time.Parse("2006-01-02", "2021-07-14")
	if err != nil {
		return err
	}
	splashLines := gosplash.MakeLines(gosplash.Settings{
		User:          "qdm12",
		Repository:    "deunhealth",
		Authors:       []string{"github.com/qdm12"},
		Emails:        []string{"quentin.mcgaw@gmail.com"},
		Version:       buildInfo.Version,
		Commit:        buildInfo.Commit,
		BuildDate:     buildInfo.BuildDate,
		Announcement:  "",
		AnnounceExp:   announcementExpiration,
		PaypalUser:    "qmcgaw",
		GithubSponsor: "qdm12",
	})
	fmt.Println(strings.Join(splashLines, "\n"))

	config, warnings, err := configReader.ReadConfig()
	for _, warning := range warnings {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	logger = logger.NewChild(logging.Settings{Level: config.Log.Level})

	docker, err := docker.New(config.Docker.Host)
	if err != nil {
		return err
	}

	looper := loop.New(docker, logger)
	looperHandler, looperCtx, looperDone := goshutdown.NewGoRoutineHandler("loop")
	go func() {
		defer close(looperDone)
		if err := looper.Run(looperCtx); err != nil {
			logger.Error(err.Error())
			cancel()
		}
	}()

	healthcheck := func() error { return nil }
	heathcheckLogger := logger.NewChild(logging.Settings{Prefix: "healthcheck: "})
	healthServer := health.NewServer(config.Health.Address, heathcheckLogger, healthcheck)
	healthServerHandler, healthServerCtx, healthServerDone := goshutdown.NewGoRoutineHandler("health")
	go func() {
		defer close(healthServerDone)
		if err := healthServer.Run(healthServerCtx); err != nil {
			logger.Error(err.Error())
		}
	}()

	group := goshutdown.NewGroupHandler("group")
	group.Add(looperHandler, healthServerHandler)

	<-ctx.Done()
	return group.Shutdown(context.Background())
}
