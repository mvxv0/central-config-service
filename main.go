package main

import (
	"central-config-service/handlers"
	"central-config-service/repositories"
	"central-config-service/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	handler := handlers.NewConfigHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")

	fmt.Println("server uspesno pokrenut")

	http.ListenAndServe(":8000", router)
}
