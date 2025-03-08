package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	// Canal para manejar señales de sistema
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("Server starting on port %s", s.httpServer.Addr)
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Esperar por señales o errores
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error starting server: %w", err)
	case sig := <-signals:
		log.Printf("Received signal: %v", sig)

		// Crear context con timeout para shutdown graceful
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Intentar shutdown graceful
		if err := s.httpServer.Shutdown(ctx); err != nil {
			// Forzar cierre si el shutdown graceful falla
			s.httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
