package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"trabajo-web/logic"
	"trabajo-web/views"
)

type PrendasHandler struct {
	prendasService  *logic.PrendasService
	usuariosService *logic.UsuariosService
	looksService    *logic.LooksService
}

func NewPrendasHandler(
	service *logic.PrendasService,
	usuariosService *logic.UsuariosService,
	looksService *logic.LooksService,
) *PrendasHandler {
	return &PrendasHandler{
		prendasService:  service,
		usuariosService: usuariosService,
		looksService:    looksService,
	}
}

func normalizarFiltro(filtro string) string {
	switch filtro {
	case "todos", "prendas", "looks":
		return filtro
	default:
		return "todos"
	}
}

// GET / - Muestra Login/Registro si no hay sesión, o el Guardarropa si hay cookie válida
func (h *PrendasHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	filtro := normalizarFiltro(r.URL.Query().Get("filtro"))

	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		renderHTML(w, r, http.StatusOK, views.IndexPage(nil, nil, nil, filtro, ""))
		return
	}

	usuario, errUser := h.usuariosService.GetUsuariosByID(r.Context(), usuarioID)
	if errUser != nil {
		if errors.Is(errUser, sql.ErrNoRows) {
			borrarCookieSesion(w, r)
			renderHTML(w, r, http.StatusOK, views.IndexPage(nil, nil, nil, filtro, "Sesión expirada. Volvé a ingresar."))
			return
		}
		log.Printf("Error consultando usuario %d en IndexPage: %v", usuarioID, errUser)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	prendas, errPrendas := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	if errPrendas != nil {
		log.Printf("Error obteniendo prendas del usuario %d: %v", usuarioID, errPrendas)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	looks, errLooks := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)
	if errLooks != nil {
		log.Printf("Error obteniendo looks del usuario %d: %v", usuarioID, errLooks)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	renderHTML(w, r, http.StatusOK, views.IndexPage(&usuario, prendas, looks, filtro, ""))
}

// POST /login - Valida credenciales con bcrypt y crea cookie firmada
func (h *PrendasHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	password := r.FormValue("password")

	if nombre == "" || password == "" {
		renderHTML(w, r, http.StatusBadRequest, views.IndexPage(nil, nil, nil, "todos", "Por favor completá usuario y contraseña."))
		return
	}

	usuario, errUser := h.usuariosService.GetUsuariosByNombre(r.Context(), nombre)
	if errUser != nil {
		if errors.Is(errUser, sql.ErrNoRows) {
			h.usuariosService.VerificarPasswordDummy(password)
			renderHTML(w, r, http.StatusUnauthorized, views.IndexPage(nil, nil, nil, "todos", "Usuario o contraseña incorrectos."))
			return
		}
		log.Printf("Error de base de datos consultando usuario en Login: %v", errUser)
		renderHTML(w, r, http.StatusInternalServerError, views.IndexPage(nil, nil, nil, "todos", "Error interno del servidor. Intentá nuevamente."))
		return
	}

	if !h.usuariosService.VerificarPassword(usuario.Password, password) {
		renderHTML(w, r, http.StatusUnauthorized, views.IndexPage(nil, nil, nil, "todos", "Usuario o contraseña incorrectos."))
		return
	}

	establecerCookieSesion(w, r, usuario.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET / POST /logout - Cierra sesión eliminando la cookie
func (h *PrendasHandler) Logout(w http.ResponseWriter, r *http.Request) {
	borrarCookieSesion(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /prendas - Crea prenda y devuelve la grilla + select y contadores OOB
func (h *PrendasHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.prendasService.CreatePrendas(r.Context(), logic.CreatePrendasParams{
		UsuarioID: usuarioID,
		Nombre:    r.FormValue("nombre"),
		Categoria: r.FormValue("categoria"),
		UrlImagen: r.FormValue("url_imagen"),
	})
	if err != nil {
		var valErr logic.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, "Error al guardar la prenda: "+valErr.Msg, http.StatusBadRequest)
			return
		}
		log.Printf("Error interno al crear prenda: %v", err)
		http.Error(w, "Error interno al guardar la prenda", http.StatusInternalServerError)
		return
	}

	prendas, errPrendas := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	looks, errLooks := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)
	if errPrendas != nil || errLooks != nil {
		log.Printf("Prenda guardada pero falló recarga OOB (prendas=%v, looks=%v)", errPrendas, errLooks)
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	renderHTML(
		w, r, http.StatusOK,
		views.PrendasList(prendas),
		views.PrendasSelect(prendas, true),
		views.ContadoresOOB(len(prendas), len(looks)),
	)
}

// DELETE /prendas/{id} - Elimina prenda y actualiza grilla, select, looks y contadores OOB
func (h *PrendasHandler) DeleteFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		responderNoAutenticado(w, r)
		return
	}

	id64, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id64 <= 0 {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}
	id := int32(id64)

	if err := h.prendasService.DeletePrendas(r.Context(), id, usuarioID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Prenda no encontrada", http.StatusNotFound)
			return
		}
		log.Printf("Error eliminando prenda %d: %v", id, err)
		http.Error(w, "Error interno al eliminar la prenda", http.StatusInternalServerError)
		return
	}

	prendas, errPrendas := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	looks, errLooks := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)
	if errPrendas != nil || errLooks != nil {
		log.Printf("Prenda eliminada pero falló recarga OOB (prendas=%v, looks=%v)", errPrendas, errLooks)
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	renderHTML(
		w, r, http.StatusOK,
		views.PrendasList(prendas),
		views.PrendasSelect(prendas, true),
		views.LooksList(looks, true),
		views.ContadoresOOB(len(prendas), len(looks)),
	)
}
