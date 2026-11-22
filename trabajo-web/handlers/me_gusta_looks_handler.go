package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"trabajo-web/logic"
	"trabajo-web/views"
)

type MeGustaLooksHandler struct {
	meGustaLooksService *logic.MeGustaLooksService
	usuariosService     *logic.UsuariosService
}

func NewMeGustaLooksHandler(service *logic.MeGustaLooksService, usuariosService *logic.UsuariosService) *MeGustaLooksHandler {
	return &MeGustaLooksHandler{
		meGustaLooksService: service,
		usuariosService:     usuariosService,
	}
}

// POST /me-gusta/toggle-web - Da o quita me gusta y devuelve solo el botón actualizado
func (h *MeGustaLooksHandler) ToggleFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		responderNoAutenticado(w, r)
		return
	}

	if !validarUsuarioActivo(w, r, h.usuariosService, usuarioID) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	lookID64, err := strconv.ParseInt(r.FormValue("look_id"), 10, 32)
	if err != nil || lookID64 <= 0 {
		http.Error(w, "Look inválido", http.StatusBadRequest)
		return
	}
	lookID := int32(lookID64)

	if err := h.meGustaLooksService.ToggleMeGusta(r.Context(), usuarioID, lookID); err != nil {
		var valErr logic.ValidationError
		switch {
		case errors.Is(err, logic.ErrSelfLike):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, logic.ErrLookNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.As(err, &valErr):
			http.Error(w, valErr.Msg, http.StatusBadRequest)
		default:
			log.Printf("Error interno al alternar me gusta (usuario %d, look %d): %v", usuarioID, lookID, err)
			http.Error(w, "Error interno al actualizar me gusta", http.StatusInternalServerError)
		}
		return
	}

	totalMeGusta, err := h.meGustaLooksService.ContarMeGustaPorLook(r.Context(), lookID)
	if err != nil {
		log.Printf("Error contando me gusta del look %d: %v", lookID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	leDioMeGusta, err := h.meGustaLooksService.UsuarioLeDioMeGusta(r.Context(), usuarioID, lookID)
	if err != nil {
		log.Printf("Error verificando me gusta (usuario %d, look %d): %v", usuarioID, lookID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	renderHTML(w, r, http.StatusOK, views.LikeButton(lookID, leDioMeGusta, totalMeGusta))
}
