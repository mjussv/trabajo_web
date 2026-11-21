package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Look{}, errors.New("el nombre del look no puede estar vacío")
	}
	if params.UsuarioID <= 0 {
		return db.Look{}, errors.New("el ID de usuario es obligatorio")
	}
	if strings.TrimSpace(params.Ocasion) == "" {
		return db.Look{}, errors.New("la ocasión es obligatoria")
	}
	if strings.TrimSpace(params.Temporada) == "" {
		return db.Look{}, errors.New("la temporada es obligatoria")
	}
	if strings.TrimSpace(params.UrlImagen) == "" {
		return db.Look{}, errors.New("la URL de imagen del look es obligatoria")
	}
	if params.PrendaEstrellaID > 0 {
		prenda, err := s.queries.ObtenerPrenda(ctx, params.PrendaEstrellaID)
		if err != nil || prenda.UsuarioID != params.UsuarioID {
			return db.Look{}, errors.New("la prenda estrella seleccionada no pertenece al usuario")
		}
	}

	return s.queries.CrearLook(ctx, db.CrearLookParams{
		UsuarioID: params.UsuarioID,
		Nombre:    params.Nombre,
		Descripcion: sql.NullString{
			String: params.Descripcion,
			Valid:  strings.TrimSpace(params.Descripcion) != "",
		},
		Ocasion:   params.Ocasion,
		Temporada: params.Temporada,
		PrendaEstrellaID: sql.NullInt32{
			Int32: params.PrendaEstrellaID,
			Valid: params.PrendaEstrellaID > 0,
		},
		UrlImagen: params.UrlImagen,
		Estilo: sql.NullString{
			String: params.Estilo,
			Valid:  strings.TrimSpace(params.Estilo) != "",
		},
	})
}

func (s *LooksService) ListLooks(ctx context.Context, viewerID int32) ([]db.ListarLooksConDetalleRow, error) {
	return s.queries.ListarLooksConDetalle(ctx, viewerID)
}

func (s *LooksService) SearchLooks(ctx context.Context, viewerID int32, query string) ([]db.BuscarLooksRow, error) {
	return s.queries.BuscarLooks(ctx, db.BuscarLooksParams{
		UsuarioID: viewerID,
		Column2:   sql.NullString{String: query, Valid: true},
	})
}

func (s *LooksService) ListLooksPorUsuario(ctx context.Context, perfilUsuarioID, viewerID int32) ([]db.ListarLooksPorUsuarioRow, error) {
	if perfilUsuarioID <= 0 {
		return nil, errors.New("ID de usuario inválido")
	}
	return s.queries.ListarLooksPorUsuario(ctx, db.ListarLooksPorUsuarioParams{
		UsuarioID:   perfilUsuarioID,
		UsuarioID_2: viewerID,
	})
}

func (s *LooksService) GetLooksByID(ctx context.Context, id int32) (db.ObtenerLookConDetalleRow, error) {
	if id <= 0 {
		return db.ObtenerLookConDetalleRow{}, errors.New("ID inválido")
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
	if params.ID <= 0 {
		return db.Look{}, errors.New("ID inválido")
	}
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Look{}, errors.New("el nombre del look no puede estar vacío")
	}
	if params.PrendaEstrellaID > 0 && params.UsuarioID > 0 {
		prenda, err := s.queries.ObtenerPrenda(ctx, params.PrendaEstrellaID)
		if err != nil || prenda.UsuarioID != params.UsuarioID {
			return db.Look{}, errors.New("la prenda estrella seleccionada no pertenece al usuario")
		}
	}

	return s.queries.ActualizarLook(ctx, db.ActualizarLookParams{
		ID:        params.ID,
		Nombre:    params.Nombre,
		Descripcion: sql.NullString{
			String: params.Descripcion,
			Valid:  strings.TrimSpace(params.Descripcion) != "",
		},
		Ocasion:   params.Ocasion,
		Temporada: params.Temporada,
		PrendaEstrellaID: sql.NullInt32{
			Int32: params.PrendaEstrellaID,
			Valid: params.PrendaEstrellaID > 0,
		},
		UrlImagen: params.UrlImagen,
		Estilo: sql.NullString{
			String: params.Estilo,
			Valid:  strings.TrimSpace(params.Estilo) != "",
		},
	})
}

func (s *LooksService) DeleteLooks(ctx context.Context, id int32) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	return s.queries.EliminarLook(ctx, id)
}
