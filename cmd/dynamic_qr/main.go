package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"dynamic-qr/internal/app"
	ctrl "dynamic-qr/internal/controller/http/v1"
)

func main() {
	_, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	var wg sync.WaitGroup

	notCtrl := ctrl.NewGeneratorCtrl()

	shutdownHTTP, err := app.StartHTTPServer(&wg, notCtrl)
	if err != nil {
		fmt.Println("Failed to start HTTP server", err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	sig := <-sigCh
	cancel(fmt.Errorf("%s signal received", sig))

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := shutdownHTTP(); err != nil {
			fmt.Println("Failed to shutdown HTTP server", err)
		}
	}()

	normalExitCh := make(chan struct{})
	go func() {
		wg.Wait()
		normalExitCh <- struct{}{}
	}()

	select {
	case sig := <-sigCh:
		fmt.Println(fmt.Sprintf("%s signal received during shutdown. Force stopping\n", sig))
		return
	case <-normalExitCh:

	}
}
