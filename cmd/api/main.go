package main

import (
	log "applicationDesignTest/internal/logger"
	repository "applicationDesignTest/internal/repository/memory"
	service "applicationDesignTest/internal/service"
	transaction "applicationDesignTest/internal/transaction/memory"
	"applicationDesignTest/internal/transport/http/handler"

	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info("starting app applicationDesignTest")

	orderRepo := repository.NewInMemoryOrderRepository()
	availabilityRepo := repository.NewInMemoryAvailabilityRepository()
	txManager := transaction.NewInMemoryTransactionManager(orderRepo, availabilityRepo)
	orderService := service.NewOrderService(orderRepo, availabilityRepo, txManager)

	log.Info("starting server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr: ":8080", // move to env
	}

	http.Handle("/orders", handler.NewCreateOrderHandler(orderService))

	go func() {
		err := httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("http server closed")
		} else if err != nil {
			log.Fatalf("failed to start http server: %s\n", err)
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 sec for example, move to const or env
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorf("failed to stop server")
		return
	}

	log.Info("server stopped")
}
