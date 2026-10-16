package logic

import (
	"context"
	"errors"
	"strings"

	db "trabajo-web/db/sqlc"
)

type PrendasService struct {
	queries *db.Queries
}

func NewPrendasService(queries *db.Queries) *PrendasService {
	return &PrendasService{
		queries: queries,
	}
}

type CreatePrendasParams struct {
	UsuarioID int32  `json:"usuario_id"`
	Nombre    string `json:"nombre"`
	Categoria string `json:"categoria"`
}

func (s *PrendasService) CreatePrendas(ctx context.Context, params CreatePrendasParams) (db.Prenda, error) {
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Prenda{}, errors.New("el nombre de la prenda no puede estar vacío")
	}

	return s.queries.CrearPrenda(ctx, db.CrearPrendaParams{
		UsuarioID: params.UsuarioID,
		Nombre:    params.Nombre,
		Categoria: params.Categoria,
	})
}

func (s *PrendasService) ListPrendasPorUsuario(ctx context.Context, usuarioID int32) ([]db.Prenda, error) {
	return s.queries.ListarPrendasPorUsuario(ctx, usuarioID)
}

func (s *PrendasService) GetPrendasByID(ctx context.Context, id int32) (db.Prenda, error) {
	if id <= 0 {
		return db.Prenda{}, errors.New("ID inválido")
	}
	return s.queries.ObtenerPrenda(ctx, id)
}

func (s *PrendasService) DeletePrendas(ctx context.Context, id int32) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	return s.queries.EliminarPrenda(ctx, id)
}
