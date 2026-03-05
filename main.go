package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type MyHandler struct {
}

func (MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.Write([]byte("Hello World!!!"))
	case "/test":
		w.Write([]byte("Test for test"))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Requested URI %s not found...", r.URL.Path)))
	}
}

func main() {
	var port string = "localhost:8080"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.RFC1123}
	log := zerolog.New(output).With().Timestamp().Logger()
	srv := http.Server{Addr: port, Handler: MyHandler{}}

	osChann := make(chan os.Signal, 1)
	signal.Notify(osChann, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Info().Err(err)
		}
	}()

	<-osChann

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Info().Msgf("Server shutdown failed: %v\n", err)
	}
	log.Info().Msg("Server gracefully stopped.")
}
