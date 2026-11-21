package handlers

import (
	"net/http"
)

func SetupRoutes(
	usuariosHandler *UsuariosHandler,
	prendasHandler *PrendasHandler,
	looksHandler *LooksHandler,
	meGustaLooksHandler *MeGustaLooksHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Páginas y sesión
	mux.HandleFunc("GET /", prendasHandler.IndexPage)
	mux.HandleFunc("POST /login", prendasHandler.Login)
	mux.HandleFunc("GET /logout", prendasHandler.Logout)
	mux.HandleFunc("POST /logout", prendasHandler.Logout)
	mux.HandleFunc("POST /usuarios/crear-web", usuariosHandler.CreateFromForm)

	// Prendas (HTMX)
	mux.HandleFunc("POST /prendas", prendasHandler.CreateFromForm)
	mux.HandleFunc("DELETE /prendas/{id}", prendasHandler.DeleteFromForm)

	// Looks (HTMX)
	mux.HandleFunc("POST /looks", looksHandler.CreateFromForm)
	mux.HandleFunc("DELETE /looks/{id}", looksHandler.DeleteFromForm)

	// Explorar y perfiles
	mux.HandleFunc("GET /explorar", looksHandler.ExplorarPage)
	mux.HandleFunc("GET /perfil/{id}", looksHandler.PerfilUsuarioPage)

	// Me gusta (HTMX)
	mux.HandleFunc("POST /me-gusta/toggle-web", meGustaLooksHandler.ToggleFromForm)

	// Archivos estáticos
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	return mux
}
