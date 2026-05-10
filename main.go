package main

import (
	"central-config-service/handlers"
	"central-config-service/repositories"
	"central-config-service/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	configRepo := repositories.NewConfigInMemRepository()
	configService := services.NewConfigService(configRepo)
	configHandler := handlers.NewConfigHandler(configService)

	groupRepo := repositories.NewConfigGroupInMemRepository()
	groupService := services.NewConfigGroupService(groupRepo)
	groupHandler := handlers.NewConfigGroupHandler(groupService)

	router := mux.NewRouter()

	router.HandleFunc("/configs", configHandler.Add).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")

	router.HandleFunc("/groups", groupHandler.AddGroup).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.GetGroup).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.DeleteGroup).Methods("DELETE")
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.AddConfigToGroup).Methods("POST")

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		log.Println("Server je uspešno pokrenut na portu 8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Greška pri pokretanju servera: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Gasim server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Greška prilikom gašenja servera: %v", err)
	}

	log.Println("Server je uspešno ugašen")
}
