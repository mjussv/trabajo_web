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

	// servir el archivo index.html en la raíz (http://localhost:8080/)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	// servir otros archivos estáticos desde ./static
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// rutas de usuarios
	mux.HandleFunc("POST /usuarios", usuariosHandler.Create)
	mux.HandleFunc("GET /usuarios", usuariosHandler.List)
	mux.HandleFunc("GET /usuarios/{id}", usuariosHandler.GetByID)
	mux.HandleFunc("PUT /usuarios/{id}", usuariosHandler.Update)
	mux.HandleFunc("DELETE /usuarios/{id}", usuariosHandler.Delete)

	// rutas de prendas
	mux.HandleFunc("POST /prendas", prendasHandler.Create)
	mux.HandleFunc("GET /prendas", prendasHandler.List)
	mux.HandleFunc("GET /prendas/{id}", prendasHandler.GetByID)
	mux.HandleFunc("DELETE /prendas/{id}", prendasHandler.Delete)

	// rutas de looks
	mux.HandleFunc("POST /looks", looksHandler.Create)
	mux.HandleFunc("GET /looks", looksHandler.List)
	mux.HandleFunc("GET /looks/{id}", looksHandler.GetByID)
	mux.HandleFunc("DELETE /looks/{id}", looksHandler.Delete)

	// rutas de me gusta
	mux.HandleFunc("POST /me-gusta", meGustaLooksHandler.AddMeGusta)
	mux.HandleFunc("DELETE /me-gusta", meGustaLooksHandler.DeleteMeGusta)

	return mux
}
