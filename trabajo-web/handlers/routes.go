package handlers

import (
	"net/http"
	"os"
)

type soloArchivosFS struct {
	fs http.FileSystem
}

func (s soloArchivosFS) Open(name string) (http.File, error) {
	f, err := s.fs.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		f.Close()
		return nil, os.ErrNotExist
	}
	return f, nil
}

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
	mux.HandleFunc("GET /perfil/{nombre}", looksHandler.PerfilUsuarioPage)

	// Me gusta (HTMX)
	mux.HandleFunc("POST /me-gusta/toggle-web", meGustaLooksHandler.ToggleFromForm)

	// Archivos estáticos (sin listado de directorios)
	fileServer := http.FileServer(soloArchivosFS{http.Dir("./static")})
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	return mux
}
