package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/finance-dashboard/backend/internal/pkg/middlewares"
)

type HandlersMap map[string]http.HandlerFunc

type Implementation struct {
	server *http.Server
}

func New(handlers HandlersMap) (*Implementation, error) {
	mux := http.NewServeMux()
	server := &http.Server{
		// перенсти адрес в конфиг
		Addr:    "0.0.0.0:3000",
		Handler: mux,

		// todo перенести таймауты в конфиг
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	implementation := &Implementation{
		server: server,
	}
	implementation.RegisterHandlers(mux, handlers)

	return implementation, nil
}

func (i *Implementation) RegisterHandlers(mux *http.ServeMux, handlers HandlersMap) {
	for path, handler := range handlers {
		mux.Handle(path, middlewares.Log(handler))
	}
}

func (i *Implementation) Run() {
	// Запускаем сервер в горутине
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- i.server.ListenAndServe()
	}()

	// Канал для получения сигналов ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		slog.Error(fmt.Sprintf("i.server.ListenAndServe err: %v", err))

	case <-stop:
		// Даем время на завершение текущих запросов
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := i.server.Shutdown(ctx); err != nil {
			fmt.Printf("Server got termination signal: %v\n", err)
		}
	}
}
