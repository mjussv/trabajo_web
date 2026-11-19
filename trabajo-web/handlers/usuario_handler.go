package handlers

import (
	"net/http"
	"strconv"

	"trabajo-web/logic"
	"trabajo-web/views"
)

type UsuariosHandler struct {
	usuariosService *logic.UsuariosService
}

func NewUsuariosHandler(service *logic.UsuariosService) *UsuariosHandler {
	return &UsuariosHandler{usuariosService: service}
}

// POST /usuarios/crear-web - Registro (form normal + cookie + redirect)
func (h *UsuariosHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, nil, "todos", "Formulario inválido.").Render(r.Context(), w)
		return
	}

	usuario, err := h.usuariosService.CreateUsuarios(r.Context(), logic.CreateUsuariosParams{
		Nombre:   r.FormValue("nombre"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, nil, "todos", "No se pudo crear la cuenta: ese nombre de usuario ya está en uso.").Render(r.Context(), w)
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
