package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	db "trabajo-web/db/sqlc"
)

type PrendasService struct {
	queries *db.Queries
}

func NewPrendasService(queries *db.Queries) *PrendasService {
	return &PrendasService{queries: queries}
}

type CreatePrendasParams struct {
	UsuarioID int32  `json:"usuario_id"`
	Nombre    string `json:"nombre"`
	Categoria string `json:"categoria"`
	UrlImagen string `json:"url_imagen"`
}

func (s *PrendasService) CreatePrendas(ctx context.Context, params CreatePrendasParams) (db.Prenda, error) {
	if params.UsuarioID <= 0 {
		return db.Prenda{}, errors.New("el ID de usuario debe ser mayor a 0")
	}
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Prenda{}, errors.New("el nombre de la prenda no puede estar vacío")
	}
	if strings.TrimSpace(params.Categoria) == "" {
		return db.Prenda{}, errors.New("la categoría no puede estar vacía")
	}

	return s.queries.CrearPrenda(ctx, db.CrearPrendaParams{
		UsuarioID: params.UsuarioID,
		Nombre:    params.Nombre,
		Categoria: params.Categoria,
		Column4:   params.UrlImagen,
	})
}

func (s *PrendasService) ListPrendasPorUsuario(ctx context.Context, usuarioID int32) ([]db.Prenda, error) {
	if usuarioID <= 0 {
		return nil, errors.New("el ID de usuario debe ser mayor a 0")
	}
	return s.queries.ListarPrendasPorUsuario(ctx, usuarioID)
}

func (s *PrendasService) ListAllPrendas(ctx context.Context) ([]db.ListarTodasPrendasConAutoraRow, error) {
	return s.queries.ListarTodasPrendasConAutora(ctx)
}

func (s *PrendasService) SearchPrendas(ctx context.Context, q string) ([]db.BuscarPrendasRow, error) {
	return s.queries.BuscarPrendas(ctx, sql.NullString{String: q, Valid: true})
}

func (s *PrendasService) GetPrendasByID(ctx context.Context, id int32) (db.Prenda, error) {
	if id <= 0 {
		return db.Prenda{}, errors.New("ID inválido")
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
	if params.ID <= 0 {
		return db.Prenda{}, errors.New("ID de prenda inválido")
	}
	if params.UsuarioID <= 0 {
		return db.Prenda{}, errors.New("el ID de usuario debe ser mayor a 0")
	}
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Prenda{}, errors.New("el nombre de la prenda no puede estar vacío")
	}
	if strings.TrimSpace(params.Categoria) == "" {
		return db.Prenda{}, errors.New("la categoría no puede estar vacía")
	}

	return s.queries.ActualizarPrenda(ctx, db.ActualizarPrendaParams{
		ID:        params.ID,
		UsuarioID: params.UsuarioID,
		Nombre:    params.Nombre,
		Categoria: params.Categoria,
		Column5:   params.UrlImagen,
	})
}

func (s *PrendasService) DeletePrendas(ctx context.Context, id int32) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	return s.queries.EliminarPrenda(ctx, id)
}
