package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "trabajo-web/db/sqlc"
	"trabajo-web/logic"
	"trabajo-web/views"
)

type LooksHandler struct {
	looksService    *logic.LooksService
	usuariosService *logic.UsuariosService
	prendasService  *logic.PrendasService
}

func NewLooksHandler(
	service *logic.LooksService,
	usuariosService *logic.UsuariosService,
	prendasService *logic.PrendasService,
) *LooksHandler {
	return &LooksHandler{
		looksService:    service,
		usuariosService: usuariosService,
		prendasService:  prendasService,
	}
}

// ==========================================
// MÉTODOS PARA VISTAS WEB (SSR + PRG + Cookie)
// ==========================================

func (h *LooksHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	prendaEstrellaID, _ := strconv.Atoi(r.FormValue("prenda_estrella_id"))

	_, err := h.looksService.CreateLooks(r.Context(), logic.CreateLooksParams{
		UsuarioID:        usuarioID,
		Nombre:           r.FormValue("nombre"),
		Descripcion:      r.FormValue("descripcion"),
		Ocasion:          r.FormValue("ocasion"),
		Temporada:        r.FormValue("temporada"),
		PrendaEstrellaID: int32(prendaEstrellaID),
		UrlImagen:        r.FormValue("url_imagen"),
		Estilo:           r.FormValue("estilo"),
	})
	if err != nil {
		http.Error(w, "Error al crear look: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	looks, err := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.LooksList(looks).Render(r.Context(), w)
}

func (h *LooksHandler) DeleteFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	look, err := h.looksService.GetLooksByID(r.Context(), int32(id))
	if err != nil || look.UsuarioID != usuarioID {
		http.Error(w, "Look no encontrado", http.StatusNotFound)
		return
	}


	if err := h.looksService.DeleteLooks(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al eliminar look: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LooksHandler) ExplorarPage(w http.ResponseWriter, r *http.Request) {
	viewerID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	viewer, err := h.usuariosService.GetUsuariosByID(r.Context(), viewerID)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	query := r.URL.Query().Get("q")
	tipo := r.URL.Query().Get("tipo")
	if tipo == "" {
		tipo = "todos"
	}

	looks, _ := h.looksService.SearchLooks(r.Context(), viewerID, query)
	prendas, _ := h.prendasService.SearchPrendas(r.Context(), query)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.ExplorarPage(viewer, query, tipo, looks, prendas).Render(r.Context(), w)
}

func (h *LooksHandler) PerfilUsuarioPage(w http.ResponseWriter, r *http.Request) {
	viewerID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	perfilID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID de perfil inválido", http.StatusBadRequest)
		return
	}

	if int32(perfilID) == viewerID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	viewer, errV := h.usuariosService.GetUsuariosByID(r.Context(), viewerID)
	if errV != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	perfil, errP := h.usuariosService.GetUsuariosByID(r.Context(), int32(perfilID))
	if errP != nil {
		http.NotFound(w, r)
		return
	}

	looks, _ := h.looksService.ListLooksPorUsuario(r.Context(), int32(perfilID), viewerID)
	prendas, _ := h.prendasService.ListPrendasPorUsuario(r.Context(), int32(perfilID))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.PerfilUsuarioPage(viewer, perfil, looks, prendas).Render(r.Context(), w)
}
