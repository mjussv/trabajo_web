package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"trabajo-web/logic"
)

type MeGustaLooksHandler struct {
	meGustaLooksService *logic.MeGustaLooksService
}

func NewMeGustaLooksHandler(service *logic.MeGustaLooksService) *MeGustaLooksHandler {
	return &MeGustaLooksHandler{meGustaLooksService: service}
}

// POST /me-gusta/toggle-web - Da o quita me gusta usando la sesión activa
func (h *MeGustaLooksHandler) ToggleFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	lookID, _ := strconv.Atoi(r.FormValue("look_id"))
	redirectTo := r.FormValue("redirect_to")
	if redirectTo == "" {
		redirectTo = "/explorar"
	}

	if err := h.meGustaLooksService.ToggleMeGusta(r.Context(), usuarioID, int32(lookID)); err != nil {
		http.Error(w, "Error al actualizar me gusta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

// ==========================================
// MÉTODOS DE LA API REST (JSON)
// ==========================================

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
