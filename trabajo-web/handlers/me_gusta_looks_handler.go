package handlers

import (
	"net/http"
	"strconv"

	"trabajo-web/logic"
	"trabajo-web/views"
)

type MeGustaLooksHandler struct {
	meGustaLooksService *logic.MeGustaLooksService
	looksService        *logic.LooksService
}

func NewMeGustaLooksHandler(
	service *logic.MeGustaLooksService,
	looksService *logic.LooksService,
) *MeGustaLooksHandler {
	return &MeGustaLooksHandler{
		meGustaLooksService: service,
		looksService:        looksService,
	}
}

// POST /me-gusta/toggle-web - Da o quita me gusta y devuelve solo el botón actualizado
func (h *MeGustaLooksHandler) ToggleFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	lookID, err := strconv.Atoi(r.FormValue("look_id"))
	if err != nil || lookID <= 0 {
		http.Error(w, "Look inválido", http.StatusBadRequest)
		return
	}

	if err := h.meGustaLooksService.ToggleMeGusta(r.Context(), usuarioID, int32(lookID)); err != nil {
		http.Error(w, "Error al actualizar me gusta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Consultamos el look actualizado para saber el total y si el usuario le dio like
	lookDetalle, err := h.looksService.GetLooksByID(r.Context(), int32(lookID))
	if err != nil {
		http.Error(w, "Look no encontrado", http.StatusNotFound)
		return
	}

	looksUsuario, err := h.looksService.ListLooksPorUsuario(r.Context(), lookDetalle.UsuarioID, usuarioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var totalMeGusta int64
	var leDioMeGusta bool
	for _, l := range looksUsuario {
		if l.ID == int32(lookID) {
			totalMeGusta = l.TotalMeGusta
			leDioMeGusta = l.LeDioMeGusta
			break
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.LikeButton(int32(lookID), leDioMeGusta, totalMeGusta).Render(r.Context(), w)
}
