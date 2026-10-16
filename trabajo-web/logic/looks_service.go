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
	return &LooksService{
		queries: queries,
	}
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

	return s.queries.CrearLook(ctx, db.CrearLookParams{
		UsuarioID: params.UsuarioID,
		Nombre:    params.Nombre,
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

func (s *LooksService) ListLooks(ctx context.Context) ([]db.ListarLooksConDetalleRow, error) {
	return s.queries.ListarLooksConDetalle(ctx)
}

func (s *LooksService) GetLooksByID(ctx context.Context, id int32) (db.ObtenerLookConDetalleRow, error) {
	if id <= 0 {
		return db.ObtenerLookConDetalleRow{}, errors.New("ID inválido")
	}
	return s.queries.ObtenerLookConDetalle(ctx, id)
}

func (s *LooksService) DeleteLooks(ctx context.Context, id int32) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	return s.queries.EliminarLook(ctx, id)
}
