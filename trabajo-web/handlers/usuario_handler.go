package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "trabajo-web/db/sqlc"
	"trabajo-web/logic"
	"trabajo-web/views"
)

type UsuariosHandler struct {
	usuariosService *logic.UsuariosService
}

func NewUsuariosHandler(service *logic.UsuariosService) *UsuariosHandler {
	return &UsuariosHandler{usuariosService: service}
}

func (h *UsuariosHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, "Formulario inválido").Render(r.Context(), w)
		return
	}

	nombre := r.FormValue("nombre")
	email := r.FormValue("email")
	password := r.FormValue("password")

	usuario, err := h.usuariosService.CreateUsuarios(r.Context(), logic.CreateUsuariosParams{
		Nombre:   nombre,
		Email:    email,
		Password: password,
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, "Error al crear el usuario: "+err.Error()).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/?usuario_id="+strconv.Itoa(int(usuario.ID)), http.StatusSeeOther)
}

func (h *UsuariosHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params logic.CreateUsuariosParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usuario, err := h.usuariosService.CreateUsuarios(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

func (h *UsuariosHandler) List(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.usuariosService.ListUsuarios(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if usuarios == nil {
		usuarios = []db.Usuario{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func (h *UsuariosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.usuariosService.GetUsuariosByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func (h *UsuariosHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var params logic.UpdateUsuariosParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.ID = int32(id)

	usuario, err := h.usuariosService.UpdateUsuarios(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func (h *UsuariosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.usuariosService.DeleteUsuarios(r.Context(), int32(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
