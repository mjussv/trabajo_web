package logic

import (
	"context"
	"errors"
	"strings"

	db "trabajo-web/db/sqlc"
)

type UsuariosService struct {
	queries *db.Queries
}

func NewUsuariosService(queries *db.Queries) *UsuariosService {
	return &UsuariosService{
		queries: queries,
	}
}

type CreateUsuariosParams struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func (s *UsuariosService) CreateUsuarios(ctx context.Context, params CreateUsuariosParams) (db.Usuario, error) {
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Usuario{}, errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(params.Email) == "" {
		return db.Usuario{}, errors.New("el email no puede estar vacío")
	}

	return s.queries.CrearUsuario(ctx, db.CrearUsuarioParams{
		Nombre: params.Nombre,
		Email:  params.Email,
	})
}

func (s *UsuariosService) ListUsuarios(ctx context.Context) ([]db.Usuario, error) {
	return s.queries.ListarUsuarios(ctx)
}

func (s *UsuariosService) GetUsuariosByID(ctx context.Context, id int32) (db.Usuario, error) {
	if id <= 0 {
		return db.Usuario{}, errors.New("el ID de usuario debe ser mayor a 0")
	}
	return s.queries.ObtenerUsuario(ctx, id)
}

type UpdateUsuariosParams struct {
	ID     int32  `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func (s *UsuariosService) UpdateUsuarios(ctx context.Context, params UpdateUsuariosParams) (db.Usuario, error) {
	if params.ID <= 0 {
		return db.Usuario{}, errors.New("el ID de usuario debe ser mayor a 0")
	}
	if strings.TrimSpace(params.Nombre) == "" {
		return db.Usuario{}, errors.New("el nombre no puede estar vacío")
	}

	return s.queries.CrearUsuario(ctx, db.CrearUsuarioParams{
		Nombre: params.Nombre,
		Email:  params.Email,
	})
}
