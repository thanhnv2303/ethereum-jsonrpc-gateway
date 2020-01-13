package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HydroProtocol/ethereum-jsonrpc-gateway/core"
	"github.com/sirupsen/logrus"
)

func main() {
	os.Exit(Run())
}

func waitExitSignal(ctxStop context.CancelFunc) {
	var exitSignal = make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGTERM)
	signal.Notify(exitSignal, syscall.SIGINT)

	<-exitSignal

	logrus.Info("Stopping...")
	ctxStop()
}

func Run() int {

	ctx, stop := context.WithCancel(context.Background())
	go waitExitSignal(stop)

	config := &core.Config{}

	logrus.Info("load config from file")
	bts, err := ioutil.ReadFile("./config.json")

	if err != nil {
		logrus.Fatal(err)
	}

	_ = json.Unmarshal(bts, config)

	_, err = core.BuildRunningConfigFromConfig(ctx, config)

	// test reload config
	//go func() {
	//	time.Sleep(5 * time.Second)
	//
	//	oldRunningConfig := currentRunningConfig
	//	newRcfg, err := BuildRunningConfigFromConfig(ctx, config)
	//
	//	if err == nil {
	//		currentRunningConfig = newRcfg
	//		oldRunningConfig.stop()
	//		logrus.Info("running config changes successfully")
	//	} else {
	//		logrus.Info("running config changes failed, err: %+v", err)
	//	}
	//}()

	if err != nil {
		logrus.Fatal(err)
	}

	httpServer := &http.Server{Addr: ":3005", Handler: &core.Server{}}

	// http server graceful shutdown
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logrus.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
	}()

	logrus.Infof("Listening on http://0.0.0.0%s\n", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	logrus.Info("Stopped")
	return 0
}
