package handlers

import (
	"central-config-service/model"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (h *ConfigHandler) Create(w http.ResponseWriter, r *http.Request) {

	var cfg model.Config

	err := json.NewDecoder(r.Body).Decode(&cfg)

	if err != nil {

		http.Error(w, "nevalidan json format", http.StatusBadRequest)
		return
	}

	err = h.service.Add(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (h *ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	err := h.service.Delete(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

func (h *ConfigHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	configs, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configs)
}
