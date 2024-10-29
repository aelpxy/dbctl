package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aelpxy/dbctl/handlers"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:     "http <host:port>",
	Short:   "Start the API and serve it",
	Example: "dbctl http localhost:5000",
	Aliases: []string{"serve"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startHttpServer(args[0])
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func startHttpServer(addr string) {
	http.HandleFunc("/healthcheck", handlers.HealthCheckHandler)
	http.HandleFunc("/databases", handlers.ListDatabaseHandler)
	http.HandleFunc("/databases/{id}", handlers.RetrieveDatabaseHandler)

	log.Printf("Starting server on http://%s\n", addr)

	server := &http.Server{
		Addr: addr,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
