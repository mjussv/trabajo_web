package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func normalizarTipoExploracion(tipo string) string {
	switch tipo {
	case "todos", "looks", "prendas":
		return tipo
	default:
		return "todos"
	}
}

func (h *LooksHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		responderNoAutenticado(w, r)
		return
	}

	if !validarUsuarioActivo(w, r, h.usuariosService, usuarioID) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	var prendaEstrellaID int32
	if rawPrendaID := strings.TrimSpace(r.FormValue("prenda_estrella_id")); rawPrendaID != "" {
		id64, err := strconv.ParseInt(rawPrendaID, 10, 32)
		if err != nil || id64 <= 0 {
			http.Error(w, "Error al guardar el look: prenda estrella inválida", http.StatusBadRequest)
			return
		}
		prendaEstrellaID = int32(id64)
	}

	_, err := h.looksService.CreateLooks(r.Context(), logic.CreateLooksParams{
		UsuarioID:        usuarioID,
		Nombre:           r.FormValue("nombre"),
		Descripcion:      r.FormValue("descripcion"),
		Ocasion:          r.FormValue("ocasion"),
		Temporada:        r.FormValue("temporada"),
		PrendaEstrellaID: prendaEstrellaID,
		UrlImagen:        r.FormValue("url_imagen"),
		Estilo:           r.FormValue("estilo"),
	})
	if err != nil {
		var valErr logic.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, "Error al guardar el look: "+valErr.Msg, http.StatusBadRequest)
			return
		}
		log.Printf("Error interno al crear look: %v", err)
		http.Error(w, "Error interno al guardar el look", http.StatusInternalServerError)
		return
	}

	looks, errLooks := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)
	prendas, errPrendas := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	if errLooks != nil || errPrendas != nil {
		log.Printf("Look guardado pero falló recarga OOB (looks=%v, prendas=%v)", errLooks, errPrendas)
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	renderHTML(
		w, r, http.StatusOK,
		views.LooksList(looks, false),
		views.ContadoresOOB(len(prendas), len(looks)),
	)
}

func (h *LooksHandler) DeleteFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		responderNoAutenticado(w, r)
		return
	}

	if !validarUsuarioActivo(w, r, h.usuariosService, usuarioID) {
		return
	}

	id64, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id64 <= 0 {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}
	id := int32(id64)

	if err := h.looksService.DeleteLooks(r.Context(), id, usuarioID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Look no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Error eliminando look %d: %v", id, err)
		http.Error(w, "Error interno al eliminar el look", http.StatusInternalServerError)
		return
	}

	looks, errLooks := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)
	prendas, errPrendas := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	if errLooks != nil || errPrendas != nil {
		log.Printf("Look eliminado pero falló recarga OOB (looks=%v, prendas=%v)", errLooks, errPrendas)
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	renderHTML(
		w, r, http.StatusOK,
		views.LooksList(looks, false),
		views.ContadoresOOB(len(prendas), len(looks)),
	)
}

func (h *LooksHandler) ExplorarPage(w http.ResponseWriter, r *http.Request) {
	viewerID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	viewer, err := h.usuariosService.GetUsuariosByID(r.Context(), viewerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			borrarCookieSesion(w, r)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		log.Printf("Error obteniendo viewer %d en ExplorarPage: %v", viewerID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("q")
	tipo := normalizarTipoExploracion(r.URL.Query().Get("tipo"))

	var looks []db.BuscarLooksRow
	if tipo == "todos" || tipo == "looks" {
		var errLooks error
		looks, errLooks = h.looksService.SearchLooks(r.Context(), viewerID, query)
		if errLooks != nil {
			log.Printf("Error buscando looks en ExplorarPage: %v", errLooks)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
	}

	var prendas []db.BuscarPrendasRow
	if tipo == "todos" || tipo == "prendas" {
		var errPrendas error
		prendas, errPrendas = h.prendasService.SearchPrendas(r.Context(), query)
		if errPrendas != nil {
			log.Printf("Error buscando prendas en ExplorarPage: %v", errPrendas)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
	}

	renderHTML(w, r, http.StatusOK, views.ExplorarPage(viewer, query, tipo, looks, prendas))
}

func (h *LooksHandler) PerfilUsuarioPage(w http.ResponseWriter, r *http.Request) {
	viewerID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	viewer, err := h.usuariosService.GetUsuariosByID(r.Context(), viewerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			borrarCookieSesion(w, r)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		log.Printf("Error consultando viewer %d en PerfilUsuarioPage: %v", viewerID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	nombrePerfil := strings.TrimSpace(r.PathValue("nombre"))
	if nombrePerfil == "" {
		http.NotFound(w, r)
		return
	}

	perfil, err := h.usuariosService.GetUsuariosByNombre(r.Context(), nombrePerfil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Printf("Error consultando perfil %q: %v", nombrePerfil, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if perfil.ID == viewer.ID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	looks, err := h.looksService.ListLooksPorUsuario(r.Context(), perfil.ID, viewer.ID)
	if err != nil {
		log.Printf("Error listando looks del perfil %d: %v", perfil.ID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	prendas, err := h.prendasService.ListPrendasPorUsuario(r.Context(), perfil.ID)
	if err != nil {
		log.Printf("Error listando prendas del perfil %d: %v", perfil.ID, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	renderHTML(w, r, http.StatusOK, views.PerfilUsuarioPage(viewer, perfil, looks, prendas))
}
