package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"skuld/env"
	"skuld/xlogger"

	_ "skuld/encoding/form"
	_ "skuld/encoding/json"
	_ "skuld/encoding/xml"
	_ "skuld/encoding/yaml"
)

type App struct {
	envInfo env.Info
	logger  xlogger.Logger
	svr     Server
}

func New(envInfo env.Info, logger xlogger.Logger, svr Server) *App {
	if !envInfo.Envv.Valid() {
		panic(fmt.Sprintf("env is not valid: %s", envInfo.Envv))
	}
	if logger == nil {
		panic("logger is nil")
	}
	if svr == nil {
		panic("svr is nil")
	}

	return &App{
		envInfo: envInfo,
		logger:  logger,
		svr:     svr,
	}
}

func (a *App) Run() error {
	if !a.envInfo.Envv.IsProduction() {
		go func() {
			if err := http.ListenAndServe(":6789", nil); err != nil {
				a.logger.Warn("Run http.ListenAndServe", "err", err)
			}
		}()
	}

	go func() {
		if err := a.svr.Run(); err != nil {
			a.logger.Warn("Run a.svr.Run", "err", err)
		}
	}()

	signal := listenSignal()
	a.logger.Info("receive signal", "signal", signal.String())
	a.Close()

	return nil
}

func (a *App) Close() {
	ctx, fn := context.WithTimeout(context.Background(), 5*time.Second)
	_ = a.svr.Close(ctx)
	fn()
	a.logger.Info("server closed")
}
