package handlers

import (
	"encoding/json"
	"net/http"

	"trabajo-web/logic"
)

type MeGustaLooksHandler struct {
	meGustaLooksService *logic.MeGustaLooksService
}

func NewMeGustaLooksHandler(meGustaLooksService *logic.MeGustaLooksService) *MeGustaLooksHandler {
	return &MeGustaLooksHandler{
		meGustaLooksService: meGustaLooksService,
	}
}


func (h *MeGustaLooksHandler) AddMeGusta(w http.ResponseWriter, r *http.Request) {
	var req logic.CreateMeGustaParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	if err := h.meGustaLooksService.AddMeGusta(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func (h *MeGustaLooksHandler) DeleteMeGusta(w http.ResponseWriter, r *http.Request) {
	var req logic.CreateMeGustaParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	if err := h.meGustaLooksService.DeleteMeGusta(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
