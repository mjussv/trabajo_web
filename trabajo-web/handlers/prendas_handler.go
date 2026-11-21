package handlers

import (
	"net/http"
	"strconv"

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

// GET / - Muestra Login/Registro si no hay sesión, o el Guardarropa si hay cookie válida
func (h *PrendasHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	filtro := r.URL.Query().Get("filtro")
	if filtro == "" {
		filtro = "todos"
	}

	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, nil, filtro, "").Render(r.Context(), w)
		return
	}

	usuario, errUser := h.usuariosService.GetUsuariosByID(r.Context(), usuarioID)
	if errUser != nil {
		borrarCookieSesion(w)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, nil, filtro, "Sesión expirada. Volvé a ingresar.").Render(r.Context(), w)
		return
	}

	prendas, _ := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	looks, _ := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.IndexPage(&usuario, prendas, looks, filtro, "").Render(r.Context(), w)
}

// POST /login - Valida credenciales con bcrypt y crea cookie firmada
func (h *PrendasHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	password := r.FormValue("password")

	if nombre == "" || password == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, nil, "todos", "Por favor completá usuario y contraseña.").Render(r.Context(), w)
		return
	}

	usuario, errUser := h.usuariosService.GetUsuariosByNombre(r.Context(), nombre)
	if errUser != nil || !h.usuariosService.VerificarPassword(usuario.Password, password) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		views.IndexPage(nil, nil, nil, "todos", "Usuario o contraseña incorrectos.").Render(r.Context(), w)
		return
	}

	establecerCookieSesion(w, usuario.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET /logout - Cierra sesión eliminando la cookie
func (h *PrendasHandler) Logout(w http.ResponseWriter, r *http.Request) {
	borrarCookieSesion(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /prendas - Crea prenda y devuelve la grilla + select y contadores OOB
func (h *PrendasHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("HX-Redirect", "/")
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
		http.Error(w, "Error al guardar la prenda: "+err.Error(), http.StatusBadRequest)
		return
	}

	prendas, err := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	looks, _ := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.PrendasList(prendas).Render(r.Context(), w)
	views.PrendasSelect(prendas, true).Render(r.Context(), w)
	views.ContadoresOOB(len(prendas), len(looks)).Render(r.Context(), w)
}

// DELETE /prendas/{id} - Elimina prenda y actualiza grilla, select y contadores OOB
func (h *PrendasHandler) DeleteFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	prenda, err := h.prendasService.GetPrendasByID(r.Context(), int32(id))
	if err != nil || prenda.UsuarioID != usuarioID {
		http.Error(w, "Prenda no encontrada", http.StatusNotFound)
		return
	}

	if err := h.prendasService.DeletePrendas(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al eliminar la prenda: "+err.Error(), http.StatusInternalServerError)
		return
	}

	prendas, _ := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	looks, _ := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.PrendasList(prendas).Render(r.Context(), w)
	views.PrendasSelect(prendas, true).Render(r.Context(), w)
	views.ContadoresOOB(len(prendas), len(looks)).Render(r.Context(), w)
}
