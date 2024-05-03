package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func StartHTTPServer(
	wg *sync.WaitGroup,
	ctrls ...httpController,
) (shutdown func() error, err error) {
	e := echo.New()
	e.Use(middleware.CORS())

	v1 := e.Group("/api/v1")

	for _, ctrl := range ctrls {
		ctrl.Register(v1)
	}

	s := http.Server{
		Handler:           e,
		ReadHeaderTimeout: 3 * time.Second,
	}

	port := getPort()
	addr := "0.0.0.0" + ":" + port

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("start http server: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("HTTP server stopped unexpectedly", err)
		} else {
			fmt.Println("HTTP server stopped")
		}
	}()

	fmt.Println(fmt.Sprintf("HTTP server started. Listening on %s", addr))

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		fmt.Println("Stopping HTTP server")

		return s.Shutdown(ctx)
	}, nil
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	return port
}
