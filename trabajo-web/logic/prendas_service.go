package logic

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"unicode/utf8"

	db "trabajo-web/db/sqlc"
)

type PrendasService struct {
	queries *db.Queries
}

func NewPrendasService(queries *db.Queries) *PrendasService {
	return &PrendasService{queries: queries}
}

func validarURLImagen(rawURL string, obligatoria bool) error {
	if rawURL == "" {
		if obligatoria {
			return ValidationError{Msg: "la URL de imagen es obligatoria"}
		}
		return nil
	}
	u, err := url.ParseRequestURI(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return ValidationError{Msg: "la URL de imagen debe comenzar con http:// o https://"}
	}
	return nil
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
	return s.queries.BuscarPrendas(ctx, sql.NullString{String: strings.TrimSpace(q), Valid: true})
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

func (s *PrendasService) DeletePrendas(ctx context.Context, id int32) error {
	if id <= 0 {
		return ValidationError{Msg: "ID inválido"}
	}
	return s.queries.EliminarPrenda(ctx, id)
}
