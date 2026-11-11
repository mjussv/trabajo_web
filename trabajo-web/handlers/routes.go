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

	// ----------------------------------------------------
	// 1. RUTAS DE LA VISTA WEB (SSR con Templ + PRG)
	// ----------------------------------------------------
	mux.HandleFunc("GET /", prendasHandler.IndexPage)
	mux.HandleFunc("POST /login", prendasHandler.Login)
	mux.HandleFunc("POST /prendas", prendasHandler.CreateFromForm)
	mux.HandleFunc("POST /prendas/eliminar", prendasHandler.DeleteFromForm)
	mux.HandleFunc("POST /usuarios/crear-web", usuariosHandler.CreateFromForm)

	// Servidor de archivos estáticos para CSS e imágenes
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// ----------------------------------------------------
	// 2. RUTAS DE LA API REST (Mantener compatibilidad)
	// ----------------------------------------------------
	
	// Rutas de usuarios
	mux.HandleFunc("POST /usuarios", usuariosHandler.Create)
	mux.HandleFunc("GET /usuarios", usuariosHandler.List)
	mux.HandleFunc("GET /usuarios/{id}", usuariosHandler.GetByID)
	mux.HandleFunc("PUT /usuarios/{id}", usuariosHandler.Update)
	mux.HandleFunc("DELETE /usuarios/{id}", usuariosHandler.Delete)

	// Rutas API de prendas
	mux.HandleFunc("GET /prendas", prendasHandler.List)
	mux.HandleFunc("GET /prendas/{id}", prendasHandler.GetByID)
	mux.HandleFunc("PUT /prendas/{id}", prendasHandler.Update)
	mux.HandleFunc("DELETE /prendas/{id}", prendasHandler.Delete)

	// Rutas de looks
	mux.HandleFunc("POST /looks", looksHandler.Create)
	mux.HandleFunc("GET /looks", looksHandler.List)
	mux.HandleFunc("GET /looks/{id}", looksHandler.GetByID)
	mux.HandleFunc("PUT /looks/{id}", looksHandler.Update)
	mux.HandleFunc("DELETE /looks/{id}", looksHandler.Delete)

	// Rutas de me gusta
	mux.HandleFunc("POST /me-gusta", meGustaLooksHandler.AddMeGusta)
	mux.HandleFunc("DELETE /me-gusta", meGustaLooksHandler.DeleteMeGusta)

	return mux
}
