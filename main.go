package main

import (
	"central-config-service/handlers"
	"central-config-service/middlewares"
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
	groupService := services.NewConfigGroupService(groupRepo, configRepo)
	groupHandler := handlers.NewConfigGroupHandler(groupService)

	router := mux.NewRouter()

	router.Use(middlewares.RateLimitMiddleware)

	router.HandleFunc("/configs", configHandler.Create).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs", configHandler.GetAll).Methods("GET")

	router.HandleFunc("/groups", groupHandler.AddGroup).Methods("POST")
	router.HandleFunc("/groups", groupHandler.GetAllGroups).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.GetGroup).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.DeleteGroup).Methods("DELETE")
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.AddConfigToGroup).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}/configs/link", groupHandler.AddExistingConfigToGroup).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}/configs/{configName}", groupHandler.DeleteConfigFromGroup).Methods("DELETE")

	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.GetConfigsByLabels).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.DeleteConfigsByLabels).Methods("DELETE")

	//graceful sd
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}

	go func() {
		log.Println("server se pokrece na portu 8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Greska pri pokretanju: %v", err)
		}
	}()

	<-quit
	log.Println("server se gasi")

	// max 10 sec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("greska pri gasenju servera:", err)
	}

	log.Println("Server je uspesno zaustavljen.")
}
