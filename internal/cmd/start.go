package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1gkx/finstar/internal/config"
	"github.com/1gkx/finstar/internal/repositories"
	"github.com/1gkx/finstar/internal/transport"
	"github.com/go-kit/kit/log"
	"github.com/urfave/cli/v2"
)

var Start = func(c *cli.Context) error {

	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	logger := log.NewLogfmtLogger(os.Stdout)
	conf := config.Init()

	if err := repositories.Migrate(conf.DbDsn(), &repositories.MigrateCongig{
		Folder:  conf.Migration.Folder,
		Version: conf.Migration.Version,
		Log:     logger,
	}); err != nil {
		logger.Log("event", "error", "desc", err)
		return err
	}

	repo, err := repositories.New(ctx, conf.DbDsn(), logger)
	if err != nil {
		logger.Log("event", "error", "desc", err)
		return err
	}

	httpAddr := fmt.Sprintf(":%s", conf.HttpPort)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%v", <-c)
	}()

	go func() {
		logger.Log("event", "Server starting...", "address", httpAddr)
		server := &http.Server{
			Addr:    httpAddr,
			Handler: transport.New(repo, logger),
		}
		errs <- server.ListenAndServe()
	}()

	logger.Log("exit: %v\n", <-errs)

	return nil
}
