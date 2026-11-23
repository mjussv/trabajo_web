package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	db "trabajo-web/db/sqlc"
)

type LooksService struct {
	queries *db.Queries
}

func NewLooksService(queries *db.Queries) *LooksService {
	return &LooksService{queries: queries}
}

type CreateLooksParams struct {
	UsuarioID        int32  `json:"usuario_id"`
	Nombre           string `json:"nombre"`
	Descripcion      string `json:"descripcion"`
	Ocasion          string `json:"ocasion"`
	Temporada        string `json:"temporada"`
	PrendaEstrellaID int32  `json:"prenda_estrella_id"`
	UrlImagen        string `json:"url_imagen"`
	Estilo           string `json:"estilo"`
}

func (s *LooksService) CreateLooks(ctx context.Context, params CreateLooksParams) (db.Look, error) {
	nombreLimpio := strings.TrimSpace(params.Nombre)
	descLimpia := strings.TrimSpace(params.Descripcion)
	ocasionLimpia := strings.TrimSpace(params.Ocasion)
	temporadaLimpia := strings.TrimSpace(params.Temporada)
	urlLimpia := strings.TrimSpace(params.UrlImagen)
	estiloLimpio := strings.TrimSpace(params.Estilo)

	if nombreLimpio == "" {
		return db.Look{}, ValidationError{Msg: "el nombre del look no puede estar vacío"}
	}
	if utf8.RuneCountInString(nombreLimpio) > 255 {
		return db.Look{}, ValidationError{Msg: "el nombre del look no puede superar los 255 caracteres"}
	}
	if utf8.RuneCountInString(descLimpia) > 500 {
		return db.Look{}, ValidationError{Msg: "la descripción no puede superar los 500 caracteres"}
	}
	if params.UsuarioID <= 0 {
		return db.Look{}, ValidationError{Msg: "el ID de usuario es obligatorio"}
	}
	if ocasionLimpia == "" {
		return db.Look{}, ValidationError{Msg: "la ocasión es obligatoria"}
	}
	if utf8.RuneCountInString(ocasionLimpia) > 100 {
		return db.Look{}, ValidationError{Msg: "la ocasión no puede superar los 100 caracteres"}
	}
	if temporadaLimpia == "" {
		return db.Look{}, ValidationError{Msg: "la temporada es obligatoria"}
	}
	if utf8.RuneCountInString(temporadaLimpia) > 100 {
		return db.Look{}, ValidationError{Msg: "la temporada no puede superar los 100 caracteres"}
	}
	if utf8.RuneCountInString(estiloLimpio) > 100 {
		return db.Look{}, ValidationError{Msg: "el estilo no puede superar los 100 caracteres"}
	}
	if err := validarURLImagen(urlLimpia, true); err != nil {
		return db.Look{}, err
	}
	if params.PrendaEstrellaID > 0 {
		prenda, err := s.queries.ObtenerPrenda(ctx, params.PrendaEstrellaID)
		if errors.Is(err, sql.ErrNoRows) {
			return db.Look{}, ValidationError{Msg: "la prenda estrella seleccionada no existe"}
		}
		if err != nil {
			return db.Look{}, err
		}
		if prenda.UsuarioID != params.UsuarioID {
			return db.Look{}, ValidationError{Msg: "la prenda estrella seleccionada no pertenece al usuario"}
		}
	}

	return s.queries.CrearLook(ctx, db.CrearLookParams{
		UsuarioID: params.UsuarioID,
		Nombre:    nombreLimpio,
		Descripcion: sql.NullString{
			String: descLimpia,
			Valid:  descLimpia != "",
		},
		Ocasion:   ocasionLimpia,
		Temporada: temporadaLimpia,
		PrendaEstrellaID: sql.NullInt32{
			Int32: params.PrendaEstrellaID,
			Valid: params.PrendaEstrellaID > 0,
		},
		UrlImagen: urlLimpia,
		Estilo: sql.NullString{
			String: estiloLimpio,
			Valid:  estiloLimpio != "",
		},
	})
}

func (s *LooksService) ListLooks(ctx context.Context, viewerID int32) ([]db.ListarLooksConDetalleRow, error) {
	return s.queries.ListarLooksConDetalle(ctx, viewerID)
}

func (s *LooksService) SearchLooks(ctx context.Context, viewerID int32, query string) ([]db.BuscarLooksRow, error) {
	return s.queries.BuscarLooks(ctx, db.BuscarLooksParams{
		UsuarioID: viewerID,
		Column2:   sql.NullString{String: escaparLike(query), Valid: true},
	})
}

func (s *LooksService) ListLooksPorUsuario(ctx context.Context, perfilUsuarioID, viewerID int32) ([]db.ListarLooksPorUsuarioRow, error) {
	if perfilUsuarioID <= 0 {
		return nil, ValidationError{Msg: "ID de usuario inválido"}
	}
	return s.queries.ListarLooksPorUsuario(ctx, db.ListarLooksPorUsuarioParams{
		UsuarioID:   perfilUsuarioID,
		UsuarioID_2: viewerID,
	})
}

func (s *LooksService) GetLooksByID(ctx context.Context, id int32) (db.ObtenerLookConDetalleRow, error) {
	if id <= 0 {
		return db.ObtenerLookConDetalleRow{}, ValidationError{Msg: "ID inválido"}
	}
	return s.queries.ObtenerLookConDetalle(ctx, id)
}

type UpdateLooksParams struct {
	ID               int32  `json:"id"`
	UsuarioID        int32  `json:"usuario_id"`
	Nombre           string `json:"nombre"`
	Descripcion      string `json:"descripcion"`
	Ocasion          string `json:"ocasion"`
	Temporada        string `json:"temporada"`
	PrendaEstrellaID int32  `json:"prenda_estrella_id"`
	UrlImagen        string `json:"url_imagen"`
	Estilo           string `json:"estilo"`
}

func (s *LooksService) UpdateLooks(ctx context.Context, params UpdateLooksParams) (db.Look, error) {
	nombreLimpio := strings.TrimSpace(params.Nombre)
	descLimpia := strings.TrimSpace(params.Descripcion)
	ocasionLimpia := strings.TrimSpace(params.Ocasion)
	temporadaLimpia := strings.TrimSpace(params.Temporada)
	urlLimpia := strings.TrimSpace(params.UrlImagen)
	estiloLimpio := strings.TrimSpace(params.Estilo)

	if params.ID <= 0 {
		return db.Look{}, ValidationError{Msg: "ID inválido"}
	}
	if params.UsuarioID <= 0 {
		return db.Look{}, ValidationError{Msg: "el ID de usuario es obligatorio"}
	}
	if nombreLimpio == "" {
		return db.Look{}, ValidationError{Msg: "el nombre del look no puede estar vacío"}
	}
	if utf8.RuneCountInString(nombreLimpio) > 255 {
		return db.Look{}, ValidationError{Msg: "el nombre del look no puede superar los 255 caracteres"}
	}
	if utf8.RuneCountInString(descLimpia) > 500 {
		return db.Look{}, ValidationError{Msg: "la descripción no puede superar los 500 caracteres"}
	}
	if ocasionLimpia == "" || utf8.RuneCountInString(ocasionLimpia) > 100 {
		return db.Look{}, ValidationError{Msg: "la ocasión es obligatoria y no puede superar los 100 caracteres"}
	}
	if temporadaLimpia == "" || utf8.RuneCountInString(temporadaLimpia) > 100 {
		return db.Look{}, ValidationError{Msg: "la temporada es obligatoria y no puede superar los 100 caracteres"}
	}
	if utf8.RuneCountInString(estiloLimpio) > 100 {
		return db.Look{}, ValidationError{Msg: "el estilo no puede superar los 100 caracteres"}
	}
	if err := validarURLImagen(urlLimpia, true); err != nil {
		return db.Look{}, err
	}
	if params.PrendaEstrellaID > 0 {
		prenda, err := s.queries.ObtenerPrenda(ctx, params.PrendaEstrellaID)
		if errors.Is(err, sql.ErrNoRows) {
			return db.Look{}, ValidationError{Msg: "la prenda estrella seleccionada no existe"}
		}
		if err != nil {
			return db.Look{}, err
		}
		if prenda.UsuarioID != params.UsuarioID {
			return db.Look{}, ValidationError{Msg: "la prenda estrella seleccionada no pertenece al usuario"}
		}
	}

	return s.queries.ActualizarLook(ctx, db.ActualizarLookParams{
		ID:        params.ID,
		UsuarioID: params.UsuarioID,
		Nombre:    nombreLimpio,
		Descripcion: sql.NullString{
			String: descLimpia,
			Valid:  descLimpia != "",
		},
		Ocasion:   ocasionLimpia,
		Temporada: temporadaLimpia,
		PrendaEstrellaID: sql.NullInt32{
			Int32: params.PrendaEstrellaID,
			Valid: params.PrendaEstrellaID > 0,
		},
		UrlImagen: urlLimpia,
		Estilo: sql.NullString{
			String: estiloLimpio,
			Valid:  estiloLimpio != "",
		},
	})
}

func (s *LooksService) DeleteLooks(ctx context.Context, id, usuarioID int32) error {
	if id <= 0 || usuarioID <= 0 {
		return ValidationError{Msg: "ID inválido"}
	}
	rows, err := s.queries.EliminarLook(ctx, db.EliminarLookParams{
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
