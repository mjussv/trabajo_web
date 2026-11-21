package handlers

import (
	"errors"
	"log"
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
		renderHTML(w, r, http.StatusBadRequest, views.IndexPage(nil, nil, nil, "todos", "Formulario inválido."))
		return
	}

	usuario, err := h.usuariosService.CreateUsuarios(r.Context(), logic.CreateUsuariosParams{
		Nombre:   r.FormValue("nombre"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		var valErr logic.ValidationError
		status := http.StatusBadRequest
		msg := "No se pudo crear la cuenta: datos inválidos."

		switch {
		case strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key"):
			msg = "No se pudo crear la cuenta: ese nombre de usuario o email ya está en uso."
		case errors.As(err, &valErr):
			msg = "No se pudo crear la cuenta: " + valErr.Msg
		default:
			log.Printf("Error interno al registrar usuario: %v", err)
			status = http.StatusInternalServerError
			msg = "No se pudo crear la cuenta debido a un error interno."
		}

		renderHTML(w, r, status, views.IndexPage(nil, nil, nil, "todos", msg))
		return
	}

	establecerCookieSesion(w, r, usuario.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
