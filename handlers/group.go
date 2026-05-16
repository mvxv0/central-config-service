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
	w.Header().Set("Content-Type", "application/json")
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

	var configDTO model.ConfigDTO
	if err := json.NewDecoder(r.Body).Decode(&configDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddConfigToGroup(name, version, configDTO); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(configDTO)
}

func (h *ConfigGroupHandler) AddExistingConfigToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	groupVersion := vars["version"]

	var req model.LabelsConfigDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Nevalidan JSON u body-ju", http.StatusBadRequest)
		return
	}

	err = h.service.AddExistingConfigToGroup(groupName, groupVersion, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Uspesno povezan postojeci config sa grupom i dodeljene su mu labele!"}`))
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Konfiguracija je uspesno obrisana iz grupe"}`))
}

func (h *ConfigGroupHandler) GetConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	groupVersion := vars["version"]

	// query paramter after '?'
	labelsQuery := r.URL.Query().Get("labels")
	if labelsQuery == "" {
		http.Error(w, "Query parametar 'labels' je obavezan", http.StatusBadRequest)
		return
	}

	// parseLabels call
	searchLabels, err := services.ParseLabels(labelsQuery)
	if err != nil {
		http.Error(w, "Greska pri parsiranju labela: "+err.Error(), http.StatusBadRequest) // 400 Bad Request
		return
	}

	// service filter
	configs, err := h.service.GetConfigsByLabels(groupName, groupVersion, searchLabels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

func (h *ConfigGroupHandler) DeleteConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]
	groupVersion := vars["version"]

	labelsQuery := r.URL.Query().Get("labels")
	if labelsQuery == "" {
		http.Error(w, "Query parametar 'labels' je obavezan", http.StatusBadRequest)
		return
	}

	searchLabels, err := services.ParseLabels(labelsQuery)
	if err != nil {
		http.Error(w, "Greska pri parsiranju labela: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteConfigsByLabels(groupName, groupVersion, searchLabels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Konfiguracije su uspesno obrisane iz grupe na osnovu zadatih labela"}`))
}
