package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"trabajo-web/logic"
)

type LooksHandler struct {
	looksService *logic.LooksService
}

func NewLooksHandler(looksService *logic.LooksService) *LooksHandler {
	return &LooksHandler{
		looksService: looksService,
	}
}


func (h *LooksHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req logic.CreateLooksParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	looks, err := h.looksService.CreateLooks(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(looks)
}


func (h *LooksHandler) List(w http.ResponseWriter, r *http.Request) {
	looks, err := h.looksService.ListLooks(r.Context())
	if err != nil {
		http.Error(w, "Error interno al obtener looks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(looks)
}


func (h *LooksHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	looks, err := h.looksService.GetLooksByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Looks no encontrados", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(looks)
}


func (h *LooksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.looksService.DeleteLooks(r.Context(), int32(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
