package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"trabajo-web/logic"
)

type PrendasHandler struct {
	prendasService *logic.PrendasService
}

func NewPrendasHandler(service *logic.PrendasService) *PrendasHandler {
	return &PrendasHandler{prendasService: service}
}

func (h *PrendasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params logic.CreatePrendasParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prenda, err := h.prendasService.CreatePrendas(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prenda)
}

func (h *PrendasHandler) List(w http.ResponseWriter, r *http.Request) {
	usuarioIDStr := r.URL.Query().Get("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		http.Error(w, "usuario_id requerido y debe ser entero", http.StatusBadRequest)
		return
	}

	prendas, err := h.prendasService.ListPrendasPorUsuario(r.Context(), int32(usuarioID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prendas)
}

func (h *PrendasHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	prenda, err := h.prendasService.GetPrendasByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prenda)
}

func (h *PrendasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.prendasService.DeletePrendas(r.Context(), int32(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
