package logic

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	db "trabajo-web/db/sqlc"
)

type PrendasService struct {
	queries *db.Queries
}

func NewPrendasService(queries *db.Queries) *PrendasService {
	return &PrendasService{queries: queries}
}

var extensionesImagen = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
	".avif": true,
	".svg":  true,
}

var prefijosDataImagen = []string{
	"data:image/jpeg;base64,",
	"data:image/jpg;base64,",
	"data:image/png;base64,",
	"data:image/webp;base64,",
	"data:image/gif;base64,",
}

func esRespuestaDeImagen(rawURL string) bool {
	client := &http.Client{Timeout: 4 * time.Second}

	reqHead, err := http.NewRequest(http.MethodHead, rawURL, nil)
	if err != nil {
		return false
	}
	reqHead.Header.Set("User-Agent", "Mozilla/5.0")

	if resp, err := client.Do(reqHead); err == nil {
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			ct := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
			if strings.HasPrefix(ct, "image/") {
				return true
			}
			if ct != "" && !strings.HasPrefix(ct, "application/octet-stream") {
				return false
			}
		}
	}

	reqGet, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return false
	}
	reqGet.Header.Set("User-Agent", "Mozilla/5.0")
	reqGet.Header.Set("Range", "bytes=0-511")

	resp, err := client.Do(reqGet)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}

	ct := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
	if strings.HasPrefix(ct, "image/") {
		return true
	}

	buf, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil || len(buf) == 0 {
		return false
	}

	return strings.HasPrefix(http.DetectContentType(buf), "image/")
}

func validarURLImagen(rawURL string, obligatoria bool) error {
	if rawURL == "" {
		if obligatoria {
			return ValidationError{Msg: "la URL de imagen es obligatoria"}
		}
		return nil
	}

	if strings.HasPrefix(rawURL, "data:image/") {
		if len(rawURL) > 2_000_000 {
			return ValidationError{Msg: "la imagen pegada es demasiado pesada (máx. 1.5 MB)"}
		}
		for _, pref := range prefijosDataImagen {
			if strings.HasPrefix(rawURL, pref) && len(rawURL) > len(pref) {
				return nil
			}
		}
		return ValidationError{Msg: "formato de imagen no válido (usá JPG, PNG, WEBP o GIF)"}
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return ValidationError{Msg: "la URL de imagen debe comenzar con http:// o https://"}
	}

	host := strings.ToLower(u.Hostname())
	if host != "localhost" && !strings.Contains(host, ".") {
		return ValidationError{Msg: "la URL de imagen no tiene un dominio válido"}
	}

	ext := strings.ToLower(path.Ext(u.Path))
	if ext != "" {
		if !extensionesImagen[ext] {
			return ValidationError{Msg: "la URL apunta a un archivo que no es imagen (.jpg, .jpeg, .png, .webp, .gif o .avif)"}
		}
		return nil
	}

	if !esRespuestaDeImagen(rawURL) {
		return ValidationError{Msg: "el enlace no corresponde a una imagen válida (copiá la dirección directa de la imagen, no la página web)"}
	}

	return nil
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func escaparLike(q string) string {
	return likeEscaper.Replace(strings.TrimSpace(q))
}

type CreatePrendasParams struct {
	UsuarioID int32  `json:"usuario_id"`
	Nombre    string `json:"nombre"`
	Categoria string `json:"categoria"`
	UrlImagen string `json:"url_imagen"`
}

func (s *PrendasService) CreatePrendas(ctx context.Context, params CreatePrendasParams) (db.Prenda, error) {
	nombreLimpio := strings.TrimSpace(params.Nombre)
	categoriaLimpia := strings.TrimSpace(params.Categoria)
	urlLimpia := strings.TrimSpace(params.UrlImagen)

	if params.UsuarioID <= 0 {
		return db.Prenda{}, ValidationError{Msg: "el ID de usuario debe ser mayor a 0"}
	}
	if nombreLimpio == "" {
		return db.Prenda{}, ValidationError{Msg: "el nombre de la prenda no puede estar vacío"}
	}
	if utf8.RuneCountInString(nombreLimpio) > 100 {
		return db.Prenda{}, ValidationError{Msg: "el nombre de la prenda no puede superar los 100 caracteres"}
	}
	if categoriaLimpia == "" {
		return db.Prenda{}, ValidationError{Msg: "la categoría no puede estar vacía"}
	}
	if utf8.RuneCountInString(categoriaLimpia) > 50 {
		return db.Prenda{}, ValidationError{Msg: "la categoría no puede superar los 50 caracteres"}
	}
	if err := validarURLImagen(urlLimpia, false); err != nil {
		return db.Prenda{}, err
	}

	return s.queries.CrearPrenda(ctx, db.CrearPrendaParams{
		UsuarioID: params.UsuarioID,
		Nombre:    nombreLimpio,
		Categoria: categoriaLimpia,
		Column4:   urlLimpia,
	})
}

func (s *PrendasService) ListPrendasPorUsuario(ctx context.Context, usuarioID int32) ([]db.Prenda, error) {
	if usuarioID <= 0 {
		return nil, ValidationError{Msg: "el ID de usuario debe ser mayor a 0"}
	}
	return s.queries.ListarPrendasPorUsuario(ctx, usuarioID)
}

func (s *PrendasService) ListAllPrendas(ctx context.Context) ([]db.ListarTodasPrendasConAutoraRow, error) {
	return s.queries.ListarTodasPrendasConAutora(ctx)
}

func (s *PrendasService) SearchPrendas(ctx context.Context, q string) ([]db.BuscarPrendasRow, error) {
	return s.queries.BuscarPrendas(ctx, sql.NullString{String: escaparLike(q), Valid: true})
}

func (s *PrendasService) GetPrendasByID(ctx context.Context, id int32) (db.Prenda, error) {
	if id <= 0 {
		return db.Prenda{}, ValidationError{Msg: "ID inválido"}
	}
	return s.queries.ObtenerPrenda(ctx, id)
}

type UpdatePrendasParams struct {
	ID        int32  `json:"id"`
	UsuarioID int32  `json:"usuario_id"`
	Nombre    string `json:"nombre"`
	Categoria string `json:"categoria"`
	UrlImagen string `json:"url_imagen"`
}

func (s *PrendasService) UpdatePrendas(ctx context.Context, params UpdatePrendasParams) (db.Prenda, error) {
	nombreLimpio := strings.TrimSpace(params.Nombre)
	categoriaLimpia := strings.TrimSpace(params.Categoria)
	urlLimpia := strings.TrimSpace(params.UrlImagen)

	if params.ID <= 0 {
		return db.Prenda{}, ValidationError{Msg: "ID de prenda inválido"}
	}
	if params.UsuarioID <= 0 {
		return db.Prenda{}, ValidationError{Msg: "el ID de usuario debe ser mayor a 0"}
	}
	if nombreLimpio == "" {
		return db.Prenda{}, ValidationError{Msg: "el nombre de la prenda no puede estar vacío"}
	}
	if utf8.RuneCountInString(nombreLimpio) > 100 {
		return db.Prenda{}, ValidationError{Msg: "el nombre de la prenda no puede superar los 100 caracteres"}
	}
	if categoriaLimpia == "" {
		return db.Prenda{}, ValidationError{Msg: "la categoría no puede estar vacía"}
	}
	if utf8.RuneCountInString(categoriaLimpia) > 50 {
		return db.Prenda{}, ValidationError{Msg: "la categoría no puede superar los 50 caracteres"}
	}
	if err := validarURLImagen(urlLimpia, false); err != nil {
		return db.Prenda{}, err
	}

	return s.queries.ActualizarPrenda(ctx, db.ActualizarPrendaParams{
		ID:        params.ID,
		UsuarioID: params.UsuarioID,
		Nombre:    nombreLimpio,
		Categoria: categoriaLimpia,
		Column5:   urlLimpia,
	})
}

func (s *PrendasService) DeletePrendas(ctx context.Context, id, usuarioID int32) error {
	if id <= 0 || usuarioID <= 0 {
		return ValidationError{Msg: "ID inválido"}
	}
	rows, err := s.queries.EliminarPrenda(ctx, db.EliminarPrendaParams{
		ID:        id,
		UsuarioID: usuarioID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
