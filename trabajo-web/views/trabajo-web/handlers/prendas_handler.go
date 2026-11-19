package handlers

import (
	"encoding/json"
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

// Helper interno para leer el usuario autenticado desde la cookie de sesión
func obtenerUsuarioDeSesion(r *http.Request) (int32, bool) {
	cookie, err := r.Cookie("session_user_id")
	if err != nil || cookie.Value == "" {
		return 0, false
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil || id <= 0 {
		return 0, false
	}
	return int32(id), true
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
		// Si la cookie apunta a un usuario borrado, invalidamos la cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_user_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, nil, filtro, "Sesión expirada. Volvé a ingresar.").Render(r.Context(), w)
		return
	}

	prendas, _ := h.prendasService.ListPrendasPorUsuario(r.Context(), usuarioID)
	looks, _ := h.looksService.ListLooksPorUsuario(r.Context(), usuarioID, usuarioID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.IndexPage(&usuario, prendas, looks, filtro, "").Render(r.Context(), w)
}

// POST /login - Valida credenciales y crea cookie de sesión
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
	if errUser != nil || usuario.Password != password {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		views.IndexPage(nil, nil, nil, "todos", "Usuario o contraseña incorrectos.").Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_user_id",
		Value:    strconv.Itoa(int(usuario.ID)),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET /logout - Cierra sesión eliminando la cookie
func (h *PrendasHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /prendas - Crea prenda para el usuario autenticado en sesión
func (h *PrendasHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	categoria := r.FormValue("categoria")
	urlImagen := r.FormValue("url_imagen")

	_, err := h.prendasService.CreatePrendas(r.Context(), logic.CreatePrendasParams{
		UsuarioID: usuarioID,
		Nombre:    nombre,
		Categoria: categoria,
		UrlImagen: urlImagen,
	})

	if err != nil {
		http.Error(w, "Error al guardar la prenda: "+err.Error(), http.StatusBadRequest)
		return
	}

   //se vuelve a consultar lista actualizada
   prendas, err:= h.prendasService.ListPrendasPorUsuario(r.Context(),usuarioID)
   if err!=nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
   }

   //se devuelve solo el componente de lista
   w.Header().Set("Content-Type","text/html;charset=utf-8")
   views.EntityList(prendas).Render(r.Context(), w)    

}

// POST /prendas/eliminar - Elimina prenda verificando sesión
func (h *PrendasHandler) DeleteFromForm(w http.ResponseWriter, r *http.Request) {
	usuarioID, logueado := obtenerUsuarioDeSesion(r)
	if !logueado {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))

    prenda,err:= h.prendasService.GetPrendasByID(r.Context(), int32(id))
    if err != nil || prenda.UsuarioID!= usuarioID {
        http.Error(w, "Prenda no encontrada", http.StatusNotFound)
        return    
    }

	if err := h.prendasService.DeletePrendas(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al eliminar la prenda: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

