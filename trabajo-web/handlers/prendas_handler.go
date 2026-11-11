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
}

func NewPrendasHandler(service *logic.PrendasService, usuariosService *logic.UsuariosService) *PrendasHandler {
	return &PrendasHandler{
		prendasService:  service,
		usuariosService: usuariosService,
	}
}

// GET / - Muestra la Pantalla 1 (Sin login) o Pantalla 2 (Logueado con usuario_id)
func (h *PrendasHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	usuarioIDStr := r.URL.Query().Get("usuario_id")

	if usuarioIDStr == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, "").Render(r.Context(), w)
		return
	}

	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, "ID de usuario inválido").Render(r.Context(), w)
		return
	}

	usuario, errUser := h.usuariosService.GetUsuariosByID(r.Context(), int32(usuarioID))
	if errUser != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, "El usuario no existe.").Render(r.Context(), w)
		return
	}

	prendas, _ := h.prendasService.ListPrendasPorUsuario(r.Context(), int32(usuarioID))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.IndexPage(&usuario, prendas, "").Render(r.Context(), w)
}

// POST /login - Procesa el inicio de sesión por Nombre de Usuario y Contraseña
func (h *PrendasHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	password := r.FormValue("password")

	if nombre == "" || password == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(nil, nil, "Por favor completa usuario y contraseña").Render(r.Context(), w)
		return
	}

	// Buscar usuario por nombre en la base de datos
	usuario, errUser := h.usuariosService.GetUsuariosByNombre(r.Context(), nombre)
	if errUser != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, "Usuario o contraseña incorrectos.").Render(r.Context(), w)
		return
	}

	// Validar contraseña (asumiendo campo Password/Contrasena en la struct de Usuario)
	// Si no tenés encriptación aún, se compara directamente. Si tenés hash, usás bcrypt.
	if usuario.Password != password {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		views.IndexPage(nil, nil, "Usuario o contraseña incorrectos.").Render(r.Context(), w)
		return
	}

	// Redirección PRG pasando el ID interno del usuario
	http.Redirect(w, r, "/?usuario_id="+strconv.Itoa(int(usuario.ID)), http.StatusSeeOther)
}

// POST /prendas - Guarda prenda y redirige manteniendo la sesión activa
func (h *PrendasHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	usuarioID, _ := strconv.Atoi(r.FormValue("usuario_id"))
	nombre := r.FormValue("nombre")
	categoria := r.FormValue("categoria")

	_, err := h.prendasService.CreatePrendas(r.Context(), logic.CreatePrendasParams{
		UsuarioID: int32(usuarioID),
		Nombre:    nombre,
		Categoria: categoria,
	})
	if err != nil {
		http.Error(w, "Error al guardar la prenda: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirección PRG volviendo a la pantalla 2 del usuario logueado
	http.Redirect(w, r, "/?usuario_id="+strconv.Itoa(usuarioID), http.StatusSeeOther)
}

// POST /prendas/eliminar - Borra prenda y redirige manteniendo el usuario
func (h *PrendasHandler) DeleteFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	usuarioIDStr := r.FormValue("usuario_id")

	if err := h.prendasService.DeletePrendas(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al eliminar la prenda: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?usuario_id="+usuarioIDStr, http.StatusSeeOther)
}

// ==========================================
// MÉTODOS API REST (JSON)
// ==========================================

func (h *PrendasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params logic.CreatePrendasParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prenda, err := h.prendasService.CreatePrendas(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prenda)
}

func (h *PrendasHandler) List(w http.ResponseWriter, r *http.Request) {
	usuarioIDStr := r.URL.Query().Get("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		http.Error(w, "usuario_id requerido y debe ser entero", http.StatusBadRequest)
		return
	}

	prendas, err := h.prendasService.ListPrendasPorUsuario(r.Context(), int32(usuarioID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prendas)
}

func (h *PrendasHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	prenda, err := h.prendasService.GetPrendasByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prenda)
}

func (h *PrendasHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var params logic.UpdatePrendasParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.ID = int32(id)

	prenda, err := h.prendasService.UpdatePrendas(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prenda)
}

func (h *PrendasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.prendasService.DeletePrendas(r.Context(), int32(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
