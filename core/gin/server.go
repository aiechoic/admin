package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Engine *gin.Engine
	// api router for registering api routes
	ApiRouter gin.IRouter
	HttpPort  int
}

func (s *Server) Run(ctx context.Context) {
	address := fmt.Sprintf(":%d", s.HttpPort)
	srv := &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("received system signal, shutting down gracefully")
	case <-ctx.Done():
		log.Println("context cancelled, shutting down gracefully")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
}
