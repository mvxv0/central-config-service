package handlers

import (
	"central-config-service/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	service *services.ConfigService
}

func NewConfigHandler(s *services.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: s}
}

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	cfg, err := h.service.Get(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(cfg)
}
