package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheRTK/http-multiplexer/internal/api"
	"github.com/TheRTK/http-multiplexer/internal/app"
	"github.com/TheRTK/http-multiplexer/internal/config"
	"github.com/TheRTK/http-multiplexer/internal/services/request"
)

func main() {
	// Так как разрешены только встроенные пакеты, то с чтением из env заморачиваться не буду, взял бы "github.com/caarlos0/env/v6".
	cfg := &config.Config{
		PortHTTP:                "8080",      // Порт на котором работает HTTP сервер
		RequestLimitCount:       2,           // Сколько максимум запросов обслуживает сервер
		MaxOutputRequestsPerURL: 4,           // Сколько максимум запросов посылается одновременно
		OutputRequestTimeout:    time.Second, // Сколько таймаут у исходящих запросов
	}

	rs := request.NewRequestService(cfg.MaxOutputRequestsPerURL, cfg.OutputRequestTimeout)

	httpServer := api.New(app.ConfigServer(cfg, rs))

	go func() {
		if err := httpServer.Run(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("HTTP server stopped!")
			} else {
				log.Println(err.Error())
			}
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	sig := <-sigs

	log.Printf("Shutting down application... Reason: %s...", sig.String())

	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Println(err.Error())
	}

	log.Println("Bye!")
}
