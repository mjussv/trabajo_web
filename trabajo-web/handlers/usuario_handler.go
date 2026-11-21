package handlers

import (
	"net/http"
	"strings"

	"trabajo-web/logic"
	"trabajo-web/views"
)

type UsuariosHandler struct {
	usuariosService *logic.UsuariosService
}

func NewUsuariosHandler(service *logic.UsuariosService) *UsuariosHandler {
	return &UsuariosHandler{usuariosService: service}
}

// POST /usuarios/crear-web - Registro (form normal + cookie firmada + redirect)
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
		msg := "No se pudo crear la cuenta: " + err.Error()
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			msg = "No se pudo crear la cuenta: ese nombre de usuario o email ya está en uso."
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, nil, "todos", msg).Render(r.Context(), w)
		return
	}

	establecerCookieSesion(w, usuario.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
