package handlers

import (
	"central-config-service/model"
	"central-config-service/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	service *services.ConfigGroupService
}

func NewConfigGroupHandler(s *services.ConfigGroupService) *ConfigGroupHandler {
	return &ConfigGroupHandler{service: s}
}

func (h *ConfigGroupHandler) AddGroup(w http.ResponseWriter, r *http.Request) {
	var group model.ConfigGroup
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group.ID = uuid.New().String()

	if err := h.service.AddGroup(group); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *ConfigGroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	group, err := h.service.GetGroup(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(group)
}

func (h *ConfigGroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.service.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

func (h *ConfigGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	if err := h.service.DeleteGroup(name, version); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ConfigGroupHandler) AddConfigToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	var config model.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.ID = uuid.New().String()

	if err := h.service.AddConfigToGroup(name, version, config); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

func (h *ConfigGroupHandler) AddExistingConfigToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	groupVersion := vars["version"]

	var req struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Nevalidan JSON u body-ju", http.StatusBadRequest)
		return
	}

	err = h.service.AddExistingConfigToGroup(groupName, groupVersion, req.Name, req.Version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Uspesno povezan postojeci config sa grupom!"}`))
}

func (h *ConfigGroupHandler) DeleteConfigFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	groupVersion := vars["version"]
	configName := vars["configName"]

	config := model.Config{
		Name: configName,
	}

	if err := h.service.DeleteConfigFromGroup(groupName, groupVersion, config); err != nil {
		http.Error(w, "config ne postoji u grupi", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Konfiguracija je uspesno obrisana iz grupe"}`))
}
